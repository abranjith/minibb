package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"minibb/backend/utils"

	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
)

// GET /api/posts/topic/{id}
func (h *Handler) GetPostsByTopicID(w http.ResponseWriter, r *http.Request) {
	topicIDStr := chi.URLParam(r, "id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid topic ID"})
		return
	}

	// Parse cursor and limit
	var cursor *int
	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		if c, err := strconv.Atoi(cursorStr); err == nil {
			cursor = &c
		}
	}

	limit := 20 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	posts, err := h.db.GetPostsByTopicID(topicID, cursor, limit)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	// Convert markdown to HTML for each post
	type PostResponse struct {
		ID          int    `json:"id"`
		TopicID     int    `json:"topic_id"`
		PubDate     string `json:"pub_date"`
		Author      string `json:"author"`
		Content     string `json:"content"`
		ContentHTML string `json:"content_html"`
	}

	var postsResponse []PostResponse
	for _, post := range posts {
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(post.Content), &buf); err != nil {
			// If markdown conversion fails, use plain text
			postsResponse = append(postsResponse, PostResponse{
				ID:          post.ID,
				TopicID:     post.TopicID,
				PubDate:     post.PubDate.Format("2006-01-02T15:04:05Z07:00"),
				Author:      post.Author,
				Content:     post.Content,
				ContentHTML: post.Content,
			})
		} else {
			postsResponse = append(postsResponse, PostResponse{
				ID:          post.ID,
				TopicID:     post.TopicID,
				PubDate:     post.PubDate.Format("2006-01-02T15:04:05Z07:00"),
				Author:      post.Author,
				Content:     post.Content,
				ContentHTML: buf.String(),
			})
		}
	}

	var nextCursor *int
	if len(posts) > 0 {
		nextCursor = &posts[len(posts)-1].ID
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"posts":       postsResponse,
		"next_cursor": nextCursor,
		"has_more":    len(posts) == limit,
	})
}

// POST /api/post
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TopicID int    `json:"topic_id"`
		Author  string `json:"author"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid request body"})
		return
	}

	if req.TopicID == 0 || req.Author == "" || req.Content == "" {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "topic_id, author, and content are required"})
		return
	}

	// Process tripcode if present
	processedAuthor := utils.GenerateTripcode(req.Author)

	post, err := h.db.CreatePost(req.TopicID, processedAuthor, req.Content)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	// Convert markdown to HTML
	var buf bytes.Buffer
	contentHTML := req.Content
	if err := goldmark.Convert([]byte(req.Content), &buf); err == nil {
		contentHTML = buf.String()
	}

	response := struct {
		ID          int    `json:"id"`
		TopicID     int    `json:"topic_id"`
		PubDate     string `json:"pub_date"`
		Author      string `json:"author"`
		Content     string `json:"content"`
		ContentHTML string `json:"content_html"`
	}{
		ID:          post.ID,
		TopicID:     post.TopicID,
		PubDate:     post.PubDate.Format("2006-01-02T15:04:05Z07:00"),
		Author:      post.Author,
		Content:     post.Content,
		ContentHTML: contentHTML,
	}

	utils.RespondWithJSON(w, http.StatusCreated, response)
}
