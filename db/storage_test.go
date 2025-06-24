package db_test // Se asume que este test está en el paquete db_test para aislamiento.

import (
	"os"
	"strings"
	"testing"

	"libroselectronicos/db"     // Importa el paquete db
	"libroselectronicos/models" // Importa el paquete models
)

// setupTestDB crea una base de datos SQLite temporal para pruebas
// y asegura que la tabla 'libros' esté limpia.
func setupTestDB(t *testing.T) *db.AlmacenLibros { // <-- Retorna el tipo CONCRETO para que el método ClearTable sea visible
	dbPath := "./test_libros.db"
	_ = os.Remove(dbPath) // Asegurarse de que el archivo de la DB no exista al inicio del test

	// Usar la función pública NewAlmacenWithDB para crear la instancia de AlmacenLibros
	almacen := db.NewAlmacenWithDB(dbPath)

	// Limpiar la tabla antes de cada test para asegurar un estado inicial vacío
	// ¡CORRECCIÓN APLICADA AQUÍ! Usar el nuevo método público ClearTable()
	err := almacen.ClearTable() // <-- ESTA ES LA LÍNEA QUE DEBE FUNCIONAR AHORA
	if err != nil {
		almacen.Close()
		t.Fatalf("Error al limpiar la tabla de libros en DB de prueba: %v", err)
	}

	return almacen
}

// tearDownTestDB cierra la base de datos de prueba y elimina el archivo.
func tearDownTestDB(_ *testing.T, almacen *db.AlmacenLibros) {
	if almacen == nil {
		return
	}
	almacen.Close()
	dbPath := "./test_libros.db"
	_ = os.Remove(dbPath) // Elimina el archivo después de que el test termina
}

func TestAgregarLibro(t *testing.T) {
	almacen := setupTestDB(t)
	defer tearDownTestDB(t, almacen)

	libro := models.NuevoLibro(1, "Título Test", "Autor Test", 2023)
	err := almacen.AgregarLibro(libro)
	if err != nil {
		t.Fatalf("Error al agregar libro: %v", err)
	}

	// Verificar si el libro fue realmente agregado
	found, err := almacen.ObtenerLibro(1)
	if err != nil {
		t.Fatalf("No se pudo obtener el libro agregado: %v", err)
	}
	if found.GetID() != 1 || found.GetTitulo() != "Título Test" {
		t.Errorf("Libro agregado incorrectamente. Obtenido: ID:%d, Título:%s", found.GetID(), found.GetTitulo())
	}

	// Intentar agregar el mismo libro (debería fallar)
	err = almacen.AgregarLibro(libro)
	if err == nil || !strings.Contains(err.Error(), "ya existe") {
		t.Errorf("Se esperaba error de libro existente, obtenido: %v", err)
	}
}

func TestListarLibros(t *testing.T) {
	almacen := setupTestDB(t)
	defer tearDownTestDB(t, almacen)

	// Al inicio, la lista debería estar vacía
	libros := almacen.ListarLibros()
	if len(libros) != 0 {
		t.Errorf("Se esperaba lista vacía al inicio, obtenido %d elementos: %v", len(libros), libros)
	}

	// Agregar algunos libros
	almacen.AgregarLibro(models.NuevoLibro(1, "Libro A", "Autor X", 2001))
	almacen.AgregarLibro(models.NuevoLibro(2, "Libro B", "Autor Y", 2002))

	libros = almacen.ListarLibros()
	if len(libros) != 2 {
		t.Errorf("Se esperaba 2 libros, obtenido %d", len(libros))
	}
}

func TestObtenerLibro(t *testing.T) {
	almacen := setupTestDB(t)
	defer tearDownTestDB(t, almacen)

	libro := models.NuevoLibro(10, "Libro Diez", "Autor Diez", 2010)
	almacen.AgregarLibro(libro)

	found, err := almacen.ObtenerLibro(10)
	if err != nil {
		t.Fatalf("Error al obtener libro existente: %v", err)
	}
	if found.GetID() != 10 {
		t.Errorf("ID de libro incorrecto, esperado 10, obtenido %d", found.GetID())
	}

	// Intentar obtener libro no existente
	_, err = almacen.ObtenerLibro(999)
	if err == nil || err != models.ErrLibroNoEncontrado {
		t.Errorf("Se esperaba ErrLibroNoEncontrado, obtenido: %v", err)
	}
}

func TestActualizarLibro(t *testing.T) {
	almacen := setupTestDB(t)
	defer tearDownTestDB(t, almacen)

	almacen.AgregarLibro(models.NuevoLibro(1, "Original", "Auth", 2000))

	updates := map[string]interface{}{
		"titulo": "Actualizado",
		"anio":   2023,
	}
	err := almacen.ActualizarLibro(1, updates)
	if err != nil {
		t.Fatalf("Error al actualizar libro: %v", err)
	}

	updated, err := almacen.ObtenerLibro(1)
	if err != nil {
		t.Fatalf("Error al obtener libro actualizado: %v", err)
	}
	if updated.GetTitulo() != "Actualizado" || updated.GetAnio() != 2023 || updated.GetAutor() != "Auth" {
		t.Errorf("Libro no actualizado correctamente. Obtenido: ID:%d, Título:%s, Autor:%s, Año:%d",
			updated.GetID(), updated.GetTitulo(), updated.GetAutor(), updated.GetAnio())
	}

	// Intentar actualizar libro no existente
	err = almacen.ActualizarLibro(999, updates)
	if err == nil || err != models.ErrLibroNoEncontrado {
		t.Errorf("Se esperaba ErrLibroNoEncontrado al actualizar libro no existente, obtenido: %v", err)
	}
}

func TestEliminarLibro(t *testing.T) {
	almacen := setupTestDB(t)
	defer tearDownTestDB(t, almacen)

	almacen.AgregarLibro(models.NuevoLibro(1, "A Eliminar", "X", 2000))

	err := almacen.EliminarLibro(1)
	if err != nil {
		t.Fatalf("Error al eliminar libro: %v", err)
	}

	// Verificar que el libro fue eliminado
	_, err = almacen.ObtenerLibro(1)
	if err == nil || err != models.ErrLibroNoEncontrado {
		t.Errorf("Libro no eliminado, se esperaba ErrLibroNoEncontrado, obtenido: %v", err)
	}

	// Intentar eliminar libro no existente
	err = almacen.EliminarLibro(999)
	if err == nil || err != models.ErrLibroNoEncontrado {
		t.Errorf("Se esperaba ErrLibroNoEncontrado al eliminar libro no existente, obtenido: %v", err)
	}
}
