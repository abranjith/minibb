package models

import (
	"time"
)

// Board methods
func (s *SQLiteDB) GetBoards() ([]Board, error) {
	rows, err := s.db.Query(`
		SELECT b.id, b.slug, b.description, COUNT(t.id) as topic_count
		FROM boards b
		LEFT JOIN topics t ON b.id = t.board_id
		GROUP BY b.id, b.slug, b.description
		ORDER BY b.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var board Board
		if err := rows.Scan(&board.ID, &board.Slug, &board.Description, &board.TopicCount); err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	return boards, rows.Err()
}

func (s *SQLiteDB) GetBoardBySlug(slug string) (*Board, error) {
	var board Board
	err := s.db.QueryRow("SELECT id, slug, description FROM boards WHERE slug = ?", slug).
		Scan(&board.ID, &board.Slug, &board.Description)
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (s *SQLiteDB) CreateBoard(slug, description string) (*Board, error) {
	result, err := s.db.Exec("INSERT INTO boards (slug, description) VALUES (?, ?)", slug, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Board{
		ID:          int(id),
		Slug:        slug,
		Description: description,
		TopicCount:  0,
	}, nil
}

// Topic methods
func (s *SQLiteDB) GetTopicsByBoardID(boardID int, cursor *int, limit int) ([]Topic, error) {
	var query string
	var args []interface{}

	if cursor != nil {
		query = `
			SELECT id, board_id, pub_date, title, status, author, last_post_id, post_count 
			FROM topics 
			WHERE board_id = ? AND id > ? 
			ORDER BY id 
			LIMIT ?`
		args = []interface{}{boardID, *cursor, limit}
	} else {
		query = `
			SELECT id, board_id, pub_date, title, status, author, last_post_id, post_count 
			FROM topics 
			WHERE board_id = ? 
			ORDER BY id 
			LIMIT ?`
		args = []interface{}{boardID, limit}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(&topic.ID, &topic.BoardID, &topic.PubDate, &topic.Title,
			&topic.Status, &topic.Author, &topic.LastPostID, &topic.PostCount); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, rows.Err()
}

func (s *SQLiteDB) GetTopicByID(id int) (*Topic, error) {
	var topic Topic
	err := s.db.QueryRow(`
		SELECT id, board_id, pub_date, title, status, author, last_post_id, post_count 
		FROM topics WHERE id = ?`, id).
		Scan(&topic.ID, &topic.BoardID, &topic.PubDate, &topic.Title,
			&topic.Status, &topic.Author, &topic.LastPostID, &topic.PostCount)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (s *SQLiteDB) CreateTopic(boardID int, title, author string) (*Topic, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert topic
	result, err := tx.Exec(`
		INSERT INTO topics (board_id, pub_date, title, status, author, last_post_id, post_count) 
		VALUES (?, ?, ?, ?, ?, 0, 0)`,
		boardID, time.Now(), title, TopicStatusOpen, author)
	if err != nil {
		return nil, err
	}

	topicID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetTopicByID(int(topicID))
}

// Post methods
func (s *SQLiteDB) GetPostsByTopicID(topicID int, cursor *int, limit int) ([]Post, error) {
	var query string
	var args []interface{}

	if cursor != nil {
		query = `
			SELECT id, topic_id, pub_date, author, content 
			FROM posts 
			WHERE topic_id = ? AND id > ? 
			ORDER BY id 
			LIMIT ?`
		args = []interface{}{topicID, *cursor, limit}
	} else {
		query = `
			SELECT id, topic_id, pub_date, author, content 
			FROM posts 
			WHERE topic_id = ? 
			ORDER BY id 
			LIMIT ?`
		args = []interface{}{topicID, limit}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.TopicID, &post.PubDate, &post.Author, &post.Content); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func (s *SQLiteDB) GetPostByID(id int) (*Post, error) {
	var post Post
	err := s.db.QueryRow(`
		SELECT id, topic_id, pub_date, author, content 
		FROM posts WHERE id = ?`, id).
		Scan(&post.ID, &post.TopicID, &post.PubDate, &post.Author, &post.Content)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *SQLiteDB) CreatePost(topicID int, author, content string) (*Post, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert post
	result, err := tx.Exec(`
		INSERT INTO posts (topic_id, pub_date, author, content) 
		VALUES (?, ?, ?, ?)`,
		topicID, time.Now(), author, content)
	if err != nil {
		return nil, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Update topic's last_post_id and post_count
	_, err = tx.Exec(`
		UPDATE topics 
		SET last_post_id = ?, post_count = post_count + 1 
		WHERE id = ?`,
		postID, topicID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetPostByID(int(postID))
}

// PostReply methods
func (s *SQLiteDB) GetRepliesByPostID(postID int) ([]PostReply, error) {
	rows, err := s.db.Query(`
		SELECT id, post_id, pub_date, author, content
		FROM post_replies
		WHERE post_id = ?
		ORDER BY pub_date ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []PostReply
	for rows.Next() {
		var reply PostReply
		if err := rows.Scan(&reply.ID, &reply.PostID, &reply.PubDate, &reply.Author, &reply.Content); err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}
	return replies, rows.Err()
}

func (s *SQLiteDB) CreatePostReply(postID int, author, content string) (*PostReply, error) {
	result, err := s.db.Exec(`
		INSERT INTO post_replies (post_id, pub_date, author, content)
		VALUES (?, ?, ?, ?)
	`, postID, time.Now(), author, content)
	if err != nil {
		return nil, err
	}

	replyID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var reply PostReply
	err = s.db.QueryRow(`
		SELECT id, post_id, pub_date, author, content
		FROM post_replies
		WHERE id = ?
	`, replyID).Scan(&reply.ID, &reply.PostID, &reply.PubDate, &reply.Author, &reply.Content)

	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (s *SQLiteDB) DeletePostReply(replyID int) error {
	_, err := s.db.Exec("DELETE FROM post_replies WHERE id = ?", replyID)
	return err
}
