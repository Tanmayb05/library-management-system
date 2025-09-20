package database

import (
	"database/sql"
	"fmt"
	"library-management/internal/logger"
	"library-management/internal/models"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB() (*DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "librarydb")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"host": host,
		"port": port,
		"user": user,
		"dbname": dbname,
		"sslmode": sslmode,
	}).Info("Successfully connected to database")
	return &DB{conn: conn}, nil
}

func (db *DB) Close() {
	if db.conn != nil {
		db.conn.Close()
	}
}

// Ping verifies the database connection is still alive
func (db *DB) Ping() error {
	if db.conn == nil {
		return fmt.Errorf("database connection is nil")
	}
	return db.conn.Ping()
}

// Book CRUD operations
func (db *DB) CreateBook(book *models.CreateBookRequest) (*models.Book, error) {
	query := `
		INSERT INTO books (title, author, isbn, publication_year, available)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, author, isbn, publication_year, available, created_at, updated_at
	`

	var result models.Book
	err := db.conn.QueryRow(
		query,
		book.Title,
		book.Author,
		book.ISBN,
		book.PublicationYear,
		true,
	).Scan(
		&result.ID,
		&result.Title,
		&result.Author,
		&result.ISBN,
		&result.PublicationYear,
		&result.Available,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	return &result, nil
}

func (db *DB) GetBookByID(id int) (*models.Book, error) {
	query := `
		SELECT id, title, author, isbn, publication_year, available, created_at, updated_at
		FROM books WHERE id = $1
	`

	var book models.Book
	err := db.conn.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.PublicationYear,
		&book.Available,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return &book, nil
}

func (db *DB) GetAllBooks() ([]*models.Book, error) {
	query := `
		SELECT id, title, author, isbn, publication_year, available, created_at, updated_at
		FROM books ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.PublicationYear,
			&book.Available,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, &book)
	}

	return books, nil
}

func (db *DB) UpdateBook(id int, updates *models.UpdateBookRequest) (*models.Book, error) {
	// Build dynamic query based on provided fields
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *updates.Title)
		argIndex++
	}
	if updates.Author != nil {
		setParts = append(setParts, fmt.Sprintf("author = $%d", argIndex))
		args = append(args, *updates.Author)
		argIndex++
	}
	if updates.ISBN != nil {
		setParts = append(setParts, fmt.Sprintf("isbn = $%d", argIndex))
		args = append(args, *updates.ISBN)
		argIndex++
	}
	if updates.PublicationYear != nil {
		setParts = append(setParts, fmt.Sprintf("publication_year = $%d", argIndex))
		args = append(args, *updates.PublicationYear)
		argIndex++
	}
	if updates.Available != nil {
		setParts = append(setParts, fmt.Sprintf("available = $%d", argIndex))
		args = append(args, *updates.Available)
		argIndex++
	}

	if len(setParts) == 0 {
		return db.GetBookByID(id) // No updates, return current book
	}

	query := fmt.Sprintf(`
		UPDATE books SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d
		RETURNING id, title, author, isbn, publication_year, available, created_at, updated_at
	`, fmt.Sprintf("%s", setParts[0]), argIndex)

	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf(`
			UPDATE books SET %s, %s, updated_at = CURRENT_TIMESTAMP
			WHERE id = $%d
			RETURNING id, title, author, isbn, publication_year, available, created_at, updated_at
		`, setParts[0], setParts[i], argIndex)
	}

	// Rebuild query properly
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query = fmt.Sprintf(`
		UPDATE books SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d
		RETURNING id, title, author, isbn, publication_year, available, created_at, updated_at
	`, setClause, argIndex)

	args = append(args, id)

	var book models.Book
	err := db.conn.QueryRow(query, args...).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.PublicationYear,
		&book.Available,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return &book, nil
}

func (db *DB) DeleteBook(id int) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book with id %d not found", id)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
