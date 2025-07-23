package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"minibb/backend/models"
	"minibb/backend/utils"

	_ "modernc.org/sqlite"
)

func Init() (models.DB, error) {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "minibb_new.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return models.NewSQLiteDB(db), nil
}

func runMigrations(db *sql.DB) error {
	// Create boards table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS boards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slug TEXT UNIQUE NOT NULL,
			description TEXT NOT NULL
		)
	`); err != nil {
		return err
	}

	// Create topics table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS topics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			board_id INTEGER NOT NULL,
			pub_date DATETIME NOT NULL,
			title TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'open',
			author TEXT NOT NULL,
			last_post_id INTEGER NOT NULL DEFAULT 0,
			post_count INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (board_id) REFERENCES boards (id)
		)
	`); err != nil {
		return err
	}

	// Create posts table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			topic_id INTEGER NOT NULL,
			pub_date DATETIME NOT NULL,
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			FOREIGN KEY (topic_id) REFERENCES topics (id)
		)
	`); err != nil {
		return err
	}

	// Create post_replies table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS post_replies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			pub_date DATETIME NOT NULL,
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
		)
	`); err != nil {
		return err
	}

	// Create indexes for performance
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_topics_board_id ON topics (board_id)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_topic_id ON posts (topic_id)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_post_replies_post_id ON post_replies (post_id)`); err != nil {
		return err
	}

	// Insert default boards if they don't exist
	if err := insertDefaultBoards(db); err != nil {
		return err
	}

	// Insert sample data for development/demo
	if err := insertSampleData(db); err != nil {
		return err
	}

	return nil
}

func insertDefaultBoards(db *sql.DB) error {
	defaultBoards := []struct {
		slug        string
		description string
	}{
		{"general", "General discussion board"},
		{"random", "Random topics and conversations"},
		{"tech", "Technology discussions"},
		{"gaming", "Video games and gaming discussion"},
	}

	for _, board := range defaultBoards {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO boards (slug, description) 
			VALUES (?, ?)`,
			board.slug, board.description)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertSampleData(db *sql.DB) error {
	// Check if sample data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM topics").Scan(&count)
	if err != nil {
		return err
	}

	// Only insert sample data if no topics exist
	if count > 0 {
		return nil
	}

	// Sample topics with posts
	sampleData := []struct {
		boardSlug string
		title     string
		author    string
		content   string
		replies   []struct {
			author  string
			content string
		}
	}{
		{
			boardSlug: "general",
			title:     "Welcome to MiniBB!",
			author:    "Admin !QmVyh7c8oQ",
			content:   "Welcome to **MiniBB**! This is a simple bulletin board system.\n\nFeel free to:\n- Start new topics\n- Reply to existing discussions\n- Use tripcodes for authentication (##password)\n\nEnjoy your stay!",
			replies: []struct {
				author  string
				content string
			}{
				{"NewUser", "Thanks for setting this up! Looking forward to the discussions."},
				{"TestUser##demo123", "Testing out the tripcode system. Pretty cool!"},
			},
		},
		{
			boardSlug: "tech",
			title:     "What programming languages are you learning?",
			author:    "CodeNinja##secure789",
			content:   "I'm curious about what everyone is working with lately.\n\nPersonally, I've been diving deep into:\n- **Go** for backend services\n- **TypeScript** for frontend work\n- **Rust** for systems programming\n\nWhat about you?",
			replies: []struct {
				author  string
				content string
			}{
				{"WebDev", "Been focusing on **React** and **Node.js** lately. The ecosystem is so rich!"},
				{"SystemsProgrammer", "**Rust** all the way! Memory safety without garbage collection is amazing."},
				{"DataScientist##ml2024", "Working with **Python** for machine learning. *Pandas* and *NumPy* are indispensable."},
			},
		},
		{
			boardSlug: "gaming",
			title:     "Best games of 2024?",
			author:    "Gamer##player1",
			content:   "What are your favorite games from this year?\n\nMy top picks:\n1. *Game Title A*\n2. *Game Title B*\n3. *Game Title C*\n\nLooking for recommendations!",
			replies: []struct {
				author  string
				content string
			}{
				{"RPGFan", "The new RPG that came out is incredible! 50+ hours and still discovering new content."},
				{"IndieGamer##pixel", "I'm all about indie games. Found some real gems on Steam this year."},
			},
		},
		{
			boardSlug: "random",
			title:     "Daily Chat Thread",
			author:    "Anonymous",
			content:   "How's everyone doing today?\n\nShare what's on your mind, what you're working on, or just say hello!",
			replies: []struct {
				author  string
				content string
			}{
				{"MorningPerson", "Good morning everyone! Beautiful day outside."},
				{"NightOwl##late", "Still up from yesterday... time has no meaning when you're coding."},
				{"CoffeeAddict", "On my third cup of coffee already. *Send help.*"},
			},
		},
	}

	for _, data := range sampleData {
		// Get board ID
		var boardID int
		err := db.QueryRow("SELECT id FROM boards WHERE slug = ?", data.boardSlug).Scan(&boardID)
		if err != nil {
			continue // Skip if board doesn't exist
		}

		// Create topic
		processedAuthor := utils.GenerateTripcode(data.author)
		result, err := db.Exec(`
			INSERT INTO topics (board_id, pub_date, title, status, author, last_post_id, post_count)
			VALUES (?, datetime('now', '-' || abs(random() % 168) || ' hours'), ?, 'open', ?, 0, 0)
		`, boardID, data.title, processedAuthor)
		if err != nil {
			continue
		}

		topicID, err := result.LastInsertId()
		if err != nil {
			continue
		}

		// Create initial post
		postResult, err := db.Exec(`
			INSERT INTO posts (topic_id, pub_date, author, content)
			VALUES (?, datetime('now', '-' || abs(random() % 168) || ' hours'), ?, ?)
		`, topicID, processedAuthor, data.content)
		if err != nil {
			continue
		}

		lastPostID, _ := postResult.LastInsertId()
		postCount := 1

		// Create replies
		for i, reply := range data.replies {
			processedReplyAuthor := utils.GenerateTripcode(reply.author)
			replyResult, err := db.Exec(`
				INSERT INTO posts (topic_id, pub_date, author, content)
				VALUES (?, datetime('now', '-' || abs(random() % (168 - ?)) || ' hours'), ?, ?)
			`, topicID, i+1, processedReplyAuthor, reply.content)
			if err != nil {
				continue
			}

			if replyPostID, err := replyResult.LastInsertId(); err == nil {
				lastPostID = replyPostID
			}
			postCount++
		}

		// Update topic with final post count and last post ID
		db.Exec(`
			UPDATE topics 
			SET last_post_id = ?, post_count = ?
			WHERE id = ?
		`, lastPostID, postCount, topicID)
	}

	return nil
}
