package handlers

import (
	"encoding/json"
	"library-management/internal/database"
	"library-management/internal/logger"
	"library-management/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BookHandler struct {
	db *database.DB
}

func NewBookHandler(db *database.DB) *BookHandler {
	return &BookHandler{db: db}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	log := logger.WithFields(map[string]interface{}{
		"handler": "CreateBook",
		"method":  r.Method,
		"path":    r.URL.Path,
	})

	log.Info("Creating new book")

	var req models.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Error("Failed to decode request body")
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Basic validation
	if req.Title == "" || req.Author == "" || req.ISBN == "" {
		log.WithFields(map[string]interface{}{
			"title":  req.Title,
			"author": req.Author,
			"isbn":   req.ISBN,
		}).Warn("Validation failed: missing required fields")
		h.sendErrorResponse(w, http.StatusBadRequest, "Title, author, and ISBN are required")
		return
	}

	book, err := h.db.CreateBook(&req)
	if err != nil {
		log.WithError(err).WithFields(map[string]interface{}{
			"title":  req.Title,
			"author": req.Author,
			"isbn":   req.ISBN,
		}).Error("Failed to create book in database")
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to create book: "+err.Error())
		return
	}

	log.WithFields(map[string]interface{}{
		"book_id": book.ID,
		"title":   book.Title,
		"author":  book.Author,
	}).Info("Book created successfully")

	h.sendSuccessResponse(w, http.StatusCreated, "Book created successfully", book)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book, err := h.db.GetBookByID(id)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Book retrieved successfully", book)
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.db.GetAllBooks()
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve books: "+err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Books retrieved successfully", books)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var req models.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	book, err := h.db.UpdateBook(id, &req)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to update book: "+err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Book updated successfully", book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	err = h.db.DeleteBook(id)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete book: "+err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Book deleted successfully", nil)
}

func (h *BookHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func (h *BookHandler) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{Message: message, Data: data})
}
