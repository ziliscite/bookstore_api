package books

import (
	"bookstore_api/tools"
	"errors"
	"regexp"
	"strings"
	"time"
)

// Custom errors for domain rules
var (
	ErrInvalidID     = errors.New("invalid ID")
	ErrNegativePrice = errors.New("price cannot be negative")
	ErrNegativeStock = errors.New("stock cannot be negative")
)

// Value Objects

type ID int64

func (i ID) Get() int64 {
	return int64(i)
}

type Title string

func (t Title) Get() string {
	return string(t)
}

type Slug string

func (s Slug) Get() string {
	return string(s)
}

type CoverImage string

func (c CoverImage) Get() string {
	return string(c)
}

type Synopsis string

func (s Synopsis) Get() string {
	return string(s)
}

type Price struct {
	value float64
}

func NewPrice(value float64) (Price, error) {
	if value < 0 {
		return Price{}, ErrNegativePrice
	}
	return Price{value: value}, nil
}

func (p Price) Get() float64 {
	return p.value
}

type Stock struct {
	value int64
}

func NewStock(value int64) (Stock, error) {
	if value < 0 {
		return Stock{}, ErrNegativeStock
	}
	return Stock{value: value}, nil
}

func (s Stock) Get() int64 {
	return s.value
}

// Book Entity - Aggregate Root
type Book struct {
	ID         ID
	Title      Title
	Slug       Slug
	CoverImage CoverImage
	Synopsis   Synopsis
	Price      Price
	Stock      Stock
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

// Whatever, I'm gonna use this for the "BookReq"

// NewBook Factory Method to create a new book
func NewBook(title string, coverImage string, synopsis string, price float64, stock int64) (*Book, error) {
	newPrice, err := NewPrice(price)
	if err != nil {
		return nil, err
	}
	newStock, err := NewStock(stock)
	if err != nil {
		return nil, err
	}
	return &Book{
		Title:      Title(title),
		Slug:       generateSlug(title),
		CoverImage: CoverImage(coverImage),
		Synopsis:   Synopsis(synopsis),
		Price:      newPrice,
		Stock:      newStock,
	}, nil
}

func (b *Book) Create(id int64) error {
	if id < 0 {
		return ErrInvalidID
	}

	b.ID = ID(id)

	defer b.MarkUpdated()

	if b.CreatedAt != nil {
		return nil
	}

	now := time.Now().UTC()
	b.CreatedAt = &now

	return nil
}

// UpdateTitle Business Logic inside the Aggregate Root
func (b *Book) UpdateTitle(newTitle string) {
	b.Title = Title(newTitle)
	b.Slug = generateSlug(newTitle)
	b.MarkUpdated()
}

func (b *Book) UpdateStock(newStock int64) error {
	stock, err := NewStock(newStock)
	if err != nil {
		return err
	}
	b.Stock = stock
	b.MarkUpdated()
	return nil
}

func (b *Book) UpdatePrice(newPrice float64) error {
	price, err := NewPrice(newPrice)
	if err != nil {
		return err
	}
	b.Price = price
	b.MarkUpdated()
	return nil
}

// MarkUpdated Helper method to mark when the entity is updated
func (b *Book) MarkUpdated() {
	if b.UpdatedAt != nil {
		return
	}
	now := time.Now().UTC()
	b.UpdatedAt = &now
}

// Slug generation logic, encapsulated in a helper function
func generateSlug(title string) Slug {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")
	re := regexp.MustCompile(`[^\w-]+`)
	s = re.ReplaceAllString(s, "")
	re = regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return Slug(s)
}

// EncryptCover Encryption logic for CoverImage
func (b *Book) EncryptCover(aesKey []byte) error {
	encryptedURL, err := tools.Encrypt([]byte(b.CoverImage.Get()), aesKey)
	if err != nil {
		return err
	}
	b.CoverImage = CoverImage(encryptedURL)
	return nil
}

// DecryptCover Decryption logic for CoverImage
func (b *Book) DecryptCover(aesKey []byte) error {
	decryptedURL, err := tools.Decrypt(b.CoverImage.Get(), aesKey)
	if err != nil {
		return err
	}
	b.CoverImage = CoverImage(decryptedURL)
	return nil
}
