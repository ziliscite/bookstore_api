package controller

import (
	"bookstore_api/internal/core/domain/books"
	"bookstore_api/internal/core/service"
	"bookstore_api/tools"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	InvalidId      = errors.New("invalid id")
	InvalidRequest = errors.New("invalid request body")
)

type httpBookDTORequest struct {
	Title string `json:"title"`

	CoverImage string `json:"cover_image"`
	Synopsis   string `json:"synopsis"`

	Price float64 `json:"price"`
	Stock int64   `json:"stock"`
}

func (d *httpBookDTORequest) newBook() (*books.Book, error) {
	book, err := books.NewBook(d.Title, d.CoverImage, d.Synopsis, d.Price, d.Stock)
	if err != nil {
		return nil, err
	}

	return book, nil
}

type httpBookDTOResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`

	Slug       string `json:"slug"`
	CoverImage string `json:"cover_image"`
	Synopsis   string `json:"synopsis"`

	Price float64 `json:"price"`
	Stock int64   `json:"stock"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func newResponseBook(book *books.Book) *httpBookDTOResponse {
	return &httpBookDTOResponse{
		ID:         book.ID.Get(),
		Title:      book.Title.Get(),
		Slug:       book.Slug.Get(),
		CoverImage: book.CoverImage.Get(),
		Synopsis:   book.Synopsis.Get(),
		Price:      book.Price.Get(),
		Stock:      book.Stock.Get(),
		CreatedAt:  book.CreatedAt,
		UpdatedAt:  book.UpdatedAt,
	}
}

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	bookDTO := &httpBookDTORequest{}
	if err := json.NewDecoder(r.Body).Decode(bookDTO); err != nil {
		tools.RespondWithError(w, InvalidRequest, http.StatusBadRequest)
		return
	}

	book, err := bookDTO.newBook()
	if err != nil {
		tools.RespondWithError(w, InvalidRequest, http.StatusBadRequest)
		return
	}

	createdBook, err := h.bookService.CreateBook(r.Context(), book)
	if err != nil {
		tools.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	tools.RespondWithJSON(w, createdBook, http.StatusCreated)
}

func (h *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		tools.RespondWithError(w, InvalidId, http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetBookById(r.Context(), int64(id))
	if err != nil {
		tools.RespondWithError(w, err, http.StatusNotFound)
		return
	}

	bookResp := newResponseBook(book)

	tools.RespondWithJSON(w, bookResp, http.StatusOK)
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	// Define a regex pattern to allow only positive integers
	re := regexp.MustCompile(`^[1-9]\d*$`)
	if !re.MatchString(pageStr) {
		tools.RespondWithError(w, fmt.Errorf("invalid page number"), http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		tools.RespondWithError(w, fmt.Errorf("invalid page number"), http.StatusBadRequest)
		return
	}

	allBooks, err := h.bookService.GetAllBooks(r.Context(), int64(page))
	if err != nil {
		tools.RespondWithError(w, err, http.StatusNotFound)
		return
	}

	var allBooksResponse []*httpBookDTOResponse
	for _, book := range allBooks {
		responseBook := newResponseBook(book)
		allBooksResponse = append(allBooksResponse, responseBook)
	}

	tools.RespondWithJSON(w, allBooksResponse, http.StatusOK)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		tools.RespondWithError(w, InvalidId, http.StatusBadRequest)
		return
	}

	bookDTO := &httpBookDTORequest{}
	if err = json.NewDecoder(r.Body).Decode(bookDTO); err != nil {
		tools.RespondWithError(w, InvalidRequest, http.StatusBadRequest)
		return
	}

	// Book currently has no id
	book, err := bookDTO.newBook()
	if err != nil {
		tools.RespondWithError(w, InvalidRequest, http.StatusBadRequest)
		return
	}

	updatedBook, err := h.bookService.UpdateBook(r.Context(), int64(id), book)
	if err != nil {
		tools.RespondWithError(w, fmt.Errorf("failed updating book: %v", err.Error()), http.StatusBadRequest)
		return
	}

	updatedBookResponse := newResponseBook(updatedBook)
	tools.RespondWithJSON(w, updatedBookResponse, http.StatusOK)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		tools.RespondWithError(w, InvalidId, http.StatusBadRequest)
		return
	}

	err = h.bookService.DeleteBook(r.Context(), int64(id))
	if err != nil {
		tools.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
