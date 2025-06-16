package models

import "errors"

// Libro con campos privados (encapsulación)
type Libro struct {
	id     int
	titulo string
	autor  string
	anio   int
}

// Constructor
func NuevoLibro(id int, titulo, autor string, anio int) *Libro {
	return &Libro{id: id, titulo: titulo, autor: autor, anio: anio}
}

// Getters
func (l *Libro) ID() int        { return l.id }
func (l *Libro) Titulo() string { return l.titulo }
func (l *Libro) Autor() string  { return l.autor }
func (l *Libro) Anio() int      { return l.anio }

// Setters
func (l *Libro) SetTitulo(t string) { l.titulo = t }
func (l *Libro) SetAutor(a string)  { l.autor = a }
func (l *Libro) SetAnio(y int)      { l.anio = y }

// Interfaz que define el comportamiento del libro
type LibroInterface interface {
	ID() int
	Titulo() string
	Autor() string
	Anio() int
	SetTitulo(string)
	SetAutor(string)
	SetAnio(int)
}

// Error para libro no encontrado
var ErrLibroNoEncontrado = errors.New("libro no encontrado")
