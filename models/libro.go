package models

import "errors"

var ErrLibroNoEncontrado = errors.New("libro no encontrado")

// LibroInterface define los métodos que cualquier tipo de libro debe implementar.
type LibroInterface interface {
	GetID() int
	GetTitulo() string
	GetAutor() string
	GetAnio() int
	SetID(id int)
	SetTitulo(titulo string)
	SetAutor(autor string)
	SetAnio(anio int)
}

// Libro representa la estructura de un libro.
type Libro struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
	Anio   int    `json:"anio"`
}

func NuevoLibro(id int, titulo, autor string, anio int) *Libro {
	return &Libro{
		ID:     id,
		Titulo: titulo,
		Autor:  autor,
		Anio:   anio,
	}
}

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

func (l *Libro) SetID(id int) {
	l.ID = id
}

func (l *Libro) SetTitulo(titulo string) {
	l.Titulo = titulo
}

func (l *Libro) SetAutor(autor string) {
	l.Autor = autor
}

func (l *Libro) SetAnio(anio int) {
	l.Anio = anio
}
