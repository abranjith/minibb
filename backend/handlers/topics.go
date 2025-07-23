package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"minibb/backend/utils"

	"github.com/go-chi/chi/v5"
)

// GET /api/topics/board/{id}
func (h *Handler) GetTopicsByBoardID(w http.ResponseWriter, r *http.Request) {
	boardIDStr := chi.URLParam(r, "id")
	boardID, err := strconv.Atoi(boardIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid board ID"})
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

	topics, err := h.db.GetTopicsByBoardID(boardID, cursor, limit)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	var nextCursor *int
	if len(topics) > 0 {
		nextCursor = &topics[len(topics)-1].ID
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"topics":      topics,
		"next_cursor": nextCursor,
		"has_more":    len(topics) == limit,
	})
}

// GET /api/topic/{id}
func (h *Handler) GetTopicByID(w http.ResponseWriter, r *http.Request) {
	topicIDStr := chi.URLParam(r, "id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid topic ID"})
		return
	}

	topic, err := h.db.GetTopicByID(topicID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, utils.APIError{Detail: "topic not found"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, topic)
}

// POST /api/topic
func (h *Handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BoardID int    `json:"board_id"`
		Title   string `json:"title"`
		Author  string `json:"author"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid request body"})
		return
	}

	if req.BoardID == 0 || req.Title == "" || req.Author == "" {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "board_id, title, and author are required"})
		return
	}

	// Process tripcode if present
	processedAuthor := utils.GenerateTripcode(req.Author)

	topic, err := h.db.CreateTopic(req.BoardID, req.Title, processedAuthor)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, topic)
}
