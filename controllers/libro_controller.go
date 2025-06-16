package controllers

import (
	"fmt"
	"libroselectronicos/db"
	"libroselectronicos/models"
)

type LibroController struct {
	almacen *db.AlmacenLibros
}

func NewLibroController() *LibroController {
	return &LibroController{almacen: db.NuevoAlmacen()}
}

func (c *LibroController) AgregarLibro(id int, titulo, autor string, anio int) {
	libro := models.NuevoLibro(id, titulo, autor, anio)
	err := c.almacen.AgregarLibro(libro)
	if err != nil {
		fmt.Println("⚠️ Error al agregar:", err)
		return
	}
	fmt.Println("✅ Libro agregado correctamente.")
}

func (c *LibroController) ListarLibros() {
	libros := c.almacen.ListarLibros()
	if len(libros) == 0 {
		fmt.Println("📭 No hay libros.")
		return
	}
	for _, libro := range libros {
		fmt.Printf("ID: %d, Título: %s, Autor: %s, Año: %d\n", libro.ID(), libro.Titulo(), libro.Autor(), libro.Anio())
	}
}
