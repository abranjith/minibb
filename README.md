# MiniBB

MiniBB is a small bulletin board modeled after phpBB and 4chan. It doesn't use any user authentication but it has a way to prevent impersonation using tripcodes.

## Features

- **Tripcode Authentication**: Users can use tripcodes to authenticate (##password format)
- **Multiple Boards**: Support for multiple discussion boards
- **Topics and Posts**: Hierarchical discussion structure
- **Markdown Support**: Basic Markdown formatting for posts
- **Read Status Tracking**: Client-side tracking of read posts using localStorage
- **Cursor-based Pagination**: Efficient pagination for large discussions
- **Rate Limiting**: Built-in rate limiting to prevent spam

## Tech Stack

### Backend
- Go with chi router
- SQLite database (modernc.org/sqlite)
- Goldmark for Markdown processing
- Rate limiting with golang.org/x/time

### Frontend
- React 19 with TypeScript
- TanStack Router for navigation
- TanStack Query for data fetching
- Tailwind CSS 4 for styling
- Vite for development

## Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- npm

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd minibb
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Install frontend dependencies:
   ```bash
   cd web
   npm install
   cd ..
   ```

### Development

Start both frontend and backend development servers:

```bash
make dev
```

This will:
- Start the Go backend on `http://localhost:8080`
- Start the Vite frontend dev server on `http://localhost:5173`
- Open both in separate terminal windows

Alternatively, you can start them manually:

**Backend:**
```bash
ENV=development go run cmd/minibb/main.go
```

**Frontend:**
```bash
cd web
npm run dev
```

### Building for Production

Build the complete application:

```bash
make build
```

This will:
1. Build the frontend and generate static files
2. Embed the frontend files in the Go binary
3. Create a single executable at `bin/minibb`

### Other Commands

- `make check` - Run linting and type checking
- `make format` - Format all code
- `make clean` - Clean build artifacts

## Usage

### Basic Navigation

1. Visit the home page to see available boards
2. Click on a board to see topics
3. Click on a topic to see posts and replies
4. Use "New Topic" to create a new discussion

### Tripcodes

To use a tripcode for authentication:

1. When posting, use the format: `YourName##YourPassword`
2. The system will generate a unique tripcode: `YourName !ABC123DEF`
3. Anyone using the same password will get the same tripcode

### Admin Features

Admin permissions are controlled via environment variables (not implemented in this basic version).

## API Endpoints

All API endpoints are under `/api/`:

- `GET /api/boards` - List all boards
- `GET /api/board/{slug}` - Get board details
- `POST /api/board` - Create new board (admin)
- `GET /api/topics/board/{id}` - List topics in a board
- `GET /api/topic/{id}` - Get topic details
- `POST /api/topic` - Create new topic
- `GET /api/posts/topic/{id}` - List posts in a topic
- `POST /api/post` - Create new post

## Database Schema

The application uses SQLite with these tables:

### Boards
- `id` - Primary key
- `slug` - URL slug
- `description` - Board description

### Topics  
- `id` - Primary key
- `board_id` - Foreign key to boards
- `pub_date` - Creation timestamp
- `title` - Topic title
- `status` - open/locked
- `author` - Author with tripcode
- `last_post_id` - ID of most recent post
- `post_count` - Total number of posts

### Posts
- `id` - Primary key  
- `topic_id` - Foreign key to topics
- `pub_date` - Creation timestamp
- `author` - Author with tripcode
- `content` - Markdown content

## Configuration

Environment variables:

- `ENV=development` - Enables development mode with CORS
- `PORT=8080` - Server port (default: 8080)

## License

[Add your license here]
