package db

import (
	"errors"
	"libroselectronicos/models"
)

type AlmacenLibros struct {
	libros map[int]models.LibroInterface
}

func NuevoAlmacen() *AlmacenLibros {
	return &AlmacenLibros{libros: make(map[int]models.LibroInterface)}
}

func (a *AlmacenLibros) AgregarLibro(l models.LibroInterface) error {
	if _, existe := a.libros[l.ID()]; existe {
		return errors.New("libro ya existe con ese ID")
	}
	a.libros[l.ID()] = l
	return nil
}

func (a *AlmacenLibros) ObtenerLibro(id int) (models.LibroInterface, error) {
	if libro, ok := a.libros[id]; ok {
		return libro, nil
	}
	return nil, models.ErrLibroNoEncontrado
}

func (a *AlmacenLibros) ListarLibros() []models.LibroInterface {
	lista := make([]models.LibroInterface, 0, len(a.libros))
	for _, libro := range a.libros {
		lista = append(lista, libro)
	}
	return lista
}
