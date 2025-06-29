package models

import "errors"

// ErrLibroNoEncontrado es un error que se devuelve cuando un libro no se encuentra.
var ErrLibroNoEncontrado = errors.New("libro no encontrado")

// ErrLibroYaExiste es un error que se devuelve cuando un libro con el mismo ID ya existe.
var ErrLibroYaExiste = errors.New("libro ya existe")

// Libro representa la estructura de un libro en la aplicación.
type Libro struct {
	ID          int    `json:"id"`
	Titulo      string `json:"titulo"`
	Autor       string `json:"autor"`
	Anio        int    `json:"anio"`
	CaratulaURL string `json:"caratula_url"` // URL a la imagen de la carátula
	Sinopsis    string `json:"sinopsis"`     // ¡NUEVO CAMPO PARA LA SINOPSIS!
}

// NuevoLibro crea una nueva instancia de Libro.
func NuevoLibro(id int, titulo, autor string, anio int) *Libro {
	return &Libro{
		ID:     id,
		Titulo: titulo,
		Autor:  autor,
		Anio:   anio,
	}
}

// NuevoLibroConCaratula crea una nueva instancia de Libro con URL de carátula.
func NuevoLibroConCaratula(id int, titulo, autor string, anio int, caratulaURL string) *Libro {
	return &Libro{
		ID:          id,
		Titulo:      titulo,
		Autor:       autor,
		Anio:        anio,
		CaratulaURL: caratulaURL,
	}
}

// NuevoLibroCompleto (Renombrado para incluir Sinopsis) crea una nueva instancia de Libro con todos los campos.
func NuevoLibroCompleto(id int, titulo, autor string, anio int, caratulaURL, sinopsis string) *Libro {
	return &Libro{
		ID:          id,
		Titulo:      titulo,
		Autor:       autor,
		Anio:        anio,
		CaratulaURL: caratulaURL,
		Sinopsis:    sinopsis, // Incluir la sinopsis
	}
}

// Métodos Getter para acceder a los campos de forma segura
func (l *Libro) GetID() int {
	return l.ID
}

func (l *Libro) GetTitulo() string {
	return l.Titulo
}

func (l *Libro) GetAutor() string {
	return l.Autor
}

func (l *Libro) GetAnio() int {
	return l.Anio
}

func (l *Libro) GetCaratulaURL() string {
	return l.CaratulaURL
}

func (l *Libro) GetSinopsis() string { // ¡NUEVO GETTER!
	return l.Sinopsis
}
