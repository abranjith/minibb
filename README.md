# MiniBB

MiniBB is a small bulletin board modeled after phpBB and 4chan. It doesn't use any user authentication but it has a way to prevent impersonation using tripcodes.

## Important Note

This repository is largely inspired by [this stream](https://www.youtube.com/watch?v=Y4_YYrIKLac) from Armin Ronacher where he builds a full-stack app app using claude code. You can look at his repo [here](https://github.com/mitsuhiko/minibb). So why this repo ?

   -  I primarily use github copilot (and windows os) and wanted to see how copilot (agentic) compares with something as advanced (as a coding agent) as claude code.
   -  Even though I have used the same repo name, instructions etc, apart from some of the scripts and utility code in go, everything else was generated either by copilot or by me.
      It may be pretty obvious looking at the code/ files.
   -  I was even able to enhance and add some more features which are not in the original repo. Primarily replies to the post, theme selection for the UI, styling etc
   -  Initial instructions provided to copilot (to mainly scaffold the backend/ frontend/ db) are in **copilot-instructions_init.md** file
   -  The results are quite remarkable. I would say about 80% plus of the code in this repo was generated using copilot agent mode (primarily using claude sonnet 4 model). Although 
      it may have taken more iterations than what is shown in Armin's video, it's still quite impressive & fun. I will share more in-depth thoughts in my blog!


## Features

- **Tripcode Authentication**: Users can use tripcodes to authenticate (##password format)
- **Multiple Boards**: Support for multiple discussion boards
- **Topics and Posts and Replies**: Hierarchical discussion structure
- **Markdown Support**: Basic Markdown formatting for posts
- **Read Status Tracking**: Client-side tracking of read posts using localStorage
- **Rate Limiting**: Built-in rate limiting to prevent spam

## Tech Stack

### Backend
- Go with chi router
- SQLite database (modernc.org/sqlite)
- Goldmark for Markdown processing
- Rate limiting with golang.org/x/time

### Frontend
- React + TypeScript
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
