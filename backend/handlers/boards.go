package handlers

import (
	"encoding/json"
	"net/http"

	"minibb/backend/models"
	"minibb/backend/utils"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db models.DB
}

func New(db models.DB) *Handler {
	return &Handler{db: db}
}

// GET /api/boards
func (h *Handler) GetBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := h.db.GetBoards()
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"boards": boards,
	})
}

// GET /api/board/{slug}
func (h *Handler) GetBoardBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "board slug is required"})
		return
	}

	board, err := h.db.GetBoardBySlug(slug)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, utils.APIError{Detail: "board not found"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, board)
}

// POST /api/board (admin only - for now, we'll skip auth)
func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid request body"})
		return
	}

	if req.Slug == "" || req.Description == "" {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "slug and description are required"})
		return
	}

	board, err := h.db.CreateBoard(req.Slug, req.Description)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, board)
}
