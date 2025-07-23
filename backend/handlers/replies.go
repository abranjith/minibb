package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"minibb/backend/utils"

	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
)

type CreatePostReplyRequest struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

// GET /api/posts/{id}/replies
func (h *Handler) GetRepliesByPostID(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid post ID"})
		return
	}

	replies, err := h.db.GetRepliesByPostID(postID)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	// Convert markdown to HTML for each reply
	type PostReplyResponse struct {
		ID          int    `json:"id"`
		PostID      int    `json:"post_id"`
		PubDate     string `json:"pub_date"`
		Author      string `json:"author"`
		Content     string `json:"content"`
		ContentHTML string `json:"content_html"`
	}

	var responseReplies []PostReplyResponse
	for _, reply := range replies {
		// Convert markdown to HTML
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(reply.Content), &buf); err != nil {
			utils.InternalServerError(w, err)
			return
		}

		responseReplies = append(responseReplies, PostReplyResponse{
			ID:          reply.ID,
			PostID:      reply.PostID,
			PubDate:     reply.PubDate.Format("2006-01-02T15:04:05Z07:00"),
			Author:      reply.Author,
			Content:     reply.Content,
			ContentHTML: buf.String(),
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"replies": responseReplies,
	})
}

// POST /api/posts/{id}/replies
func (h *Handler) CreatePostReply(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid post ID"})
		return
	}

	var req CreatePostReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid JSON"})
		return
	}

	// Basic validation
	req.Author = strings.TrimSpace(req.Author)
	req.Content = strings.TrimSpace(req.Content)

	if req.Author == "" {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "author is required"})
		return
	}

	if req.Content == "" {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "content is required"})
		return
	}

	if len(req.Author) > 100 {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "author name too long"})
		return
	}

	if len(req.Content) > 10000 {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "content too long"})
		return
	}

	// Process tripcode if present
	req.Author = utils.GenerateTripcode(req.Author)

	reply, err := h.db.CreatePostReply(postID, req.Author, req.Content)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(reply.Content), &buf); err != nil {
		utils.InternalServerError(w, err)
		return
	}

	response := struct {
		ID          int    `json:"id"`
		PostID      int    `json:"post_id"`
		PubDate     string `json:"pub_date"`
		Author      string `json:"author"`
		Content     string `json:"content"`
		ContentHTML string `json:"content_html"`
	}{
		ID:          reply.ID,
		PostID:      reply.PostID,
		PubDate:     reply.PubDate.Format("2006-01-02T15:04:05Z07:00"),
		Author:      reply.Author,
		Content:     reply.Content,
		ContentHTML: buf.String(),
	}

	utils.RespondWithJSON(w, http.StatusCreated, response)
}

// DELETE /api/replies/{id}
func (h *Handler) DeletePostReply(w http.ResponseWriter, r *http.Request) {
	replyIDStr := chi.URLParam(r, "id")
	replyID, err := strconv.Atoi(replyIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid reply ID"})
		return
	}

	err = h.db.DeletePostReply(replyID)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
