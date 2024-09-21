package handlers

import (
	"bookstore_api/internal/services"
	"bookstore_api/models"
	"bookstore_api/tools"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type BookHandler struct {
	*Handler
	bookService *services.BookService
}

func NewBookHandler(handler *Handler, bookService *services.BookService) *BookHandler {
	return &BookHandler{
		Handler:     handler,
		bookService: bookService,
	}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		tools.RespondWithError(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	createdBook, err := h.bookService.CreateBook(r.Context(), &book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tools.RespondWithJSON(w, createdBook, http.StatusCreated)
}

//func (h *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
//	idStr := chi.URLParam(r, "id")
//	cacheKey := "book" + idStr
//
//	var cachedBook models.Book
//	cachedBookJSON, err := h.Cache.Get(r.Context(), cacheKey).Result()
//	if err == nil {
//		tools.RespondWithCachedJSON(w, cachedBookJSON, &cachedBook, http.StatusOK)
//		return
//	}
//
//	id, err := strconv.Atoi(idStr)
//	if err != nil {
//		http.Error(w, "Invalid Id", http.StatusBadRequest)
//		return
//	}
//
//	book, err := h.bookService.GetBookById(r.Context(), id)
//	if err != nil {
//		http.Error(w, "Book not found", http.StatusNotFound)
//		return
//	}
//
//	ctx := context.WithValue(r.Context(), "cachedKey", cacheKey)
//	tools.RespondWithJSONAndCache(w, r.WithContext(ctx), h.Cache, book, http.StatusOK)
//}

func (h *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	cacheKey := "book" + idStr

	var cachedBook models.Book
	cachedBookJSON, err := h.Cache.Get(r.Context(), cacheKey).Result()
	if err == nil {
		tools.RespondWithCachedJSON(w, cachedBookJSON, &cachedBook, http.StatusOK)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetBookById(r.Context(), id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	ctx := context.WithValue(r.Context(), "cachedKey", cacheKey)
	tools.RespondWithJSONAndCache(w, r.WithContext(ctx), h.Cache, book, http.StatusOK)
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := query.Get("page")

	ctx := context.WithValue(r.Context(), "page", page)

	books, err := h.bookService.GetAllBooks(ctx)
	if err != nil {
		tools.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	tools.RespondWithJSON(w, books, http.StatusOK)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	cacheKey := "book" + idStr

	id, err := strconv.Atoi(idStr)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid Id"), http.StatusBadRequest)
		return
	}

	var book models.Book
	if err = json.NewDecoder(r.Body).Decode(&book); err != nil {
		tools.RespondWithError(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	book.ID = int64(id)
	updatedBook, err := h.bookService.UpdateBook(r.Context(), &book)
	if err != nil {
		tools.RespondWithError(w, errors.New("failed updating book"), http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "cachedKey", cacheKey)
	tools.RespondWithJSONAndCache(w, r.WithContext(ctx), h.Cache, updatedBook, http.StatusOK)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)

	idStr := chi.URLParam(r, "id")
	cacheKey := "book" + idStr

	id, err := strconv.Atoi(idStr)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid Id"), http.StatusBadRequest)
		return
	}

	err = h.bookService.DeleteBook(r.Context(), id)
	if err != nil {
		tools.RespondWithError(w, errors.New("failed deleting book"), http.StatusBadRequest)
		return
	}

	err = h.Cache.Del(r.Context(), cacheKey).Err()
	if err != nil {
		log.Printf("Failed to cache book: %s", err)
	}
}
