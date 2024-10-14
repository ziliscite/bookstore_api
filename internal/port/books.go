package port

import (
	"bookstore_api/internal/core/domain/books"
	"context"
)

// Ini harusnya dipisah buat bookrequest sama bookresponse
// maksudnya, kan kalau create request gitu, ga perlu id

type BookRepository interface {
	Create(ctx context.Context, book *books.Book) (*books.Book, error)
	GetById(ctx context.Context, id int64) (*books.Book, error)
	GetAll(ctx context.Context, page int64) ([]*books.Book, error)
	Update(ctx context.Context, book *books.Book) (*books.Book, error)
	Delete(ctx context.Context, id int64) error
}

//type BookReq struct {
//	Title      books.Title
//	CoverImage books.CoverImage
//	Synopsis   books.Synopsis
//	Price      books.Price
//	Stock      books.Stock
//}
