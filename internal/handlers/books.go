package handlers

import (
	"bookstore_api/internal/services"
	"bookstore_api/models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"time"
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdBook, err := h.bookService.CreateBook(r.Context(), &book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(createdBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	cacheKey := "book" + idStr

	var cachedBook models.Book
	cachedBookJSON, err := h.Cache.Get(r.Context(), cacheKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cachedBookJSON), &cachedBook)
		if err != nil {
			http.Error(w, "Failed to unmarshall cached book", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(cachedBook)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	// Convert the book to JSON before caching
	bookJSON, err := json.Marshal(book)
	if err != nil {
		http.Error(w, "Failed to process book", http.StatusInternalServerError)
		return
	}

	err = h.Cache.Set(r.Context(), cacheKey, bookJSON, time.Second*3600).Err()
	if err != nil {
		log.Printf("Failed to cache book: %s", err)
	}

	_, err = w.Write(bookJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {

	// I don't think caching all the books is a good idea.
	// Because I would then have to delete this cache everytime a book is updated or deleted.

	books, err := h.bookService.GetAllBooks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	cacheKey := "book" + idStr

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	var book models.Book
	if err = json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	book.ID = int64(id)
	updatedBook, err := h.bookService.UpdateBook(r.Context(), &book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bookJSON, err := json.Marshal(updatedBook)
	if err != nil {
		http.Error(w, "Failed to process book", http.StatusInternalServerError)
		return
	}

	// Checking wether key exist or not seems like a waste of computing power... Let it just be an error

	err = h.Cache.SetEx(r.Context(), cacheKey, bookJSON, time.Second*3600).Err()
	if err != nil {
		log.Printf("Failed to cache book: %s", err)
	}

	_, err = w.Write(bookJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	cacheKey := "book" + idStr

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	err = h.bookService.DeleteBook(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Cache.Del(r.Context(), cacheKey).Err()
	if err != nil {
		log.Printf("Failed to cache book: %s", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
