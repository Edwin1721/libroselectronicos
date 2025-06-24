package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"libroselectronicos/models"

	_ "github.com/mattn/go-sqlite3" // Driver de SQLite
)

type LibroAlmacenamiento interface {
	AgregarLibro(libro models.LibroInterface) error
	ListarLibros() []models.LibroInterface
	ObtenerLibro(id int) (models.LibroInterface, error)
	ActualizarLibro(id int, updates map[string]interface{}) error
	EliminarLibro(id int) error
	Close() error
	ClearTable() error // ¡DEBE ESTAR AQUÍ!
}

type AlmacenLibros struct {
	db *sql.DB
}

var _ LibroAlmacenamiento = (*AlmacenLibros)(nil)

func NuevoAlmacen() *AlmacenLibros {
	return NewAlmacenWithDB("./libros.db")
}

func NewAlmacenWithDB(dbPath string) *AlmacenLibros {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error al abrir la base de datos %s: %v", dbPath, err)
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS libros (
        id INTEGER PRIMARY KEY,
        titulo TEXT NOT NULL,
        autor TEXT NOT NULL,
        anio INTEGER NOT NULL
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		log.Fatalf("Error al crear la tabla de libros en %s: %v", dbPath, err)
	}

	return &AlmacenLibros{db: db}
}

func (a *AlmacenLibros) Close() error {
	return a.db.Close()
}

func (a *AlmacenLibros) ClearTable() error { // ¡DEBE ESTAR AQUÍ!
	_, err := a.db.Exec("DELETE FROM libros")
	if err != nil {
		return fmt.Errorf("error al limpiar la tabla libros: %w", err)
	}
	return nil
}

func (a *AlmacenLibros) AgregarLibro(libro models.LibroInterface) error {
	var existingID int
	err := a.db.QueryRow("SELECT id FROM libros WHERE id = ?", libro.GetID()).Scan(&existingID)
	if err == nil {
		return fmt.Errorf("libro con ID %d ya existe", libro.GetID())
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error al verificar existencia del libro: %w", err)
	}

	insertSQL := `INSERT INTO libros(id, titulo, autor, anio) VALUES(?, ?, ?, ?)`
	_, err = a.db.Exec(insertSQL, libro.GetID(), libro.GetTitulo(), libro.GetAutor(), libro.GetAnio())
	if err != nil {
		return fmt.Errorf("error al insertar libro: %w", err)
	}
	return nil
}

func (a *AlmacenLibros) ListarLibros() []models.LibroInterface {
	rows, err := a.db.Query("SELECT id, titulo, autor, anio FROM libros")
	if err != nil {
		log.Printf("Error al consultar libros: %v", err)
		return nil
	}
	defer rows.Close()

	var libros []models.LibroInterface
	for rows.Next() {
		var l models.Libro
		if err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Anio); err != nil {
			log.Printf("Error al escanear libro: %v", err)
			continue
		}
		libros = append(libros, &l)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error después de iterar filas: %v", err)
	}
	return libros
}

func (a *AlmacenLibros) ObtenerLibro(id int) (models.LibroInterface, error) {
	var l models.Libro
	err := a.db.QueryRow("SELECT id, titulo, autor, anio FROM libros WHERE id = ?", id).Scan(&l.ID, &l.Titulo, &l.Autor, &l.Anio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrLibroNoEncontrado
		}
		return nil, fmt.Errorf("error al obtener libro por ID: %w", err)
	}
	return &l, nil
}

func (a *AlmacenLibros) ActualizarLibro(id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no hay datos para actualizar")
	}

	query := "UPDATE libros SET "
	args := []interface{}{}
	parts := []string{}

	for key, value := range updates {
		parts = append(parts, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}
	query += strings.Join(parts, ", ") + " WHERE id = ?"
	args = append(args, id)

	result, err := a.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error al ejecutar actualización: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrLibroNoEncontrado
	}
	return nil
}

func (a *AlmacenLibros) EliminarLibro(id int) error {
	result, err := a.db.Exec("DELETE FROM libros WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error al ejecutar eliminación: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrLibroNoEncontrado
	}
	return nil
}
