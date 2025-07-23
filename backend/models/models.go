package models

import (
	"database/sql"
	"time"
)

type Board struct {
	ID          int    `json:"id"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	TopicCount  int    `json:"topic_count"`
}

type TopicStatus string

const (
	TopicStatusOpen   TopicStatus = "open"
	TopicStatusLocked TopicStatus = "locked"
)

type Topic struct {
	ID         int         `json:"id"`
	BoardID    int         `json:"board_id"`
	PubDate    time.Time   `json:"pub_date"`
	Title      string      `json:"title"`
	Status     TopicStatus `json:"status"`
	Author     string      `json:"author"`
	LastPostID int         `json:"last_post_id"`
	PostCount  int         `json:"post_count"`
}

type Post struct {
	ID      int       `json:"id"`
	TopicID int       `json:"topic_id"`
	PubDate time.Time `json:"pub_date"`
	Author  string    `json:"author"`
	Content string    `json:"content"`
}

type PostReply struct {
	ID      int       `json:"id"`
	PostID  int       `json:"post_id"`
	PubDate time.Time `json:"pub_date"`
	Author  string    `json:"author"`
	Content string    `json:"content"`
}

// Database interface
type DB interface {
	// Boards
	GetBoards() ([]Board, error)
	GetBoardBySlug(slug string) (*Board, error)
	CreateBoard(slug, description string) (*Board, error)

	// Topics
	GetTopicsByBoardID(boardID int, cursor *int, limit int) ([]Topic, error)
	GetTopicByID(id int) (*Topic, error)
	CreateTopic(boardID int, title, author string) (*Topic, error)

	// Posts
	GetPostsByTopicID(topicID int, cursor *int, limit int) ([]Post, error)
	CreatePost(topicID int, author, content string) (*Post, error)
	GetPostByID(id int) (*Post, error)

	// Post Replies
	GetRepliesByPostID(postID int) ([]PostReply, error)
	CreatePostReply(postID int, author, content string) (*PostReply, error)
	DeletePostReply(replyID int) error

	// Utilities
	Close() error
}

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(db *sql.DB) *SQLiteDB {
	return &SQLiteDB{db: db}
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}
