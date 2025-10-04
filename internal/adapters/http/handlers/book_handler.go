package http

import (
	"context"
	"strconv"
	"time"

	domains "example.com/practice/fiber/internal/domain"
	usecase "example.com/practice/fiber/internal/usecase/book"
	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	service *usecase.BookService
}

func NewBookHandler(service *usecase.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) RegisterRoutes(app fiber.Router) {
	app.Get("/books", h.GetAllBooks)
	app.Get("/books/:id", h.GetBookByID)
	app.Post("/books", h.CreateBook)
	app.Put("/books/:id", h.UpdateBook)
	app.Delete("/books/:id", h.DeleteBook)
}

func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}
	book := new(domains.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	book.ID = id
	if err := h.service.UpdateBook(ctx, book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}
	return c.Status(fiber.StatusOK).JSON(book)
}

func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.service.DeleteBook(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book deleted successfully",
	})
}

func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	books, err := h.service.GetAllBooks(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get all books",
		})
	}
	return c.Status(fiber.StatusOK).JSON(books)
}

func (h *BookHandler) GetBookByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	book, err := h.service.GetBookByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get book by ID",
		})
	}
	return c.Status(fiber.StatusOK).JSON(book)
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	book := new(domains.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.service.CreateBook(ctx, book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create book",
		})
	}
	return c.Status(fiber.StatusOK).JSON(book)
}
