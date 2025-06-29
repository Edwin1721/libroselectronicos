package db_test

import (
	"errors"
	"libroselectronicos/db"
	"libroselectronicos/models"
	"os"
	"testing"
)

// Nombre de la base de datos de prueba
const testDBPath = "test_libros.db"

// setupTestDB inicializa una base de datos de prueba limpia.
// Ahora devuelve la interfaz db.LibroAlmacenamiento
func setupTestDB(t *testing.T) db.LibroAlmacenamiento { // CAMBIO: de *db.AlmacenSQLite a db.LibroAlmacenamiento
	// Asegúrate de que no haya un archivo de base de datos anterior.
	os.Remove(testDBPath)

	// Llama a NewAlmacenForTest para crear una instancia de la base de datos de prueba.
	almacen := db.NewAlmacenForTest(testDBPath) // Usar el constructor correcto de la base de datos para pruebas
	if almacen == nil {
		t.Fatalf("No se pudo inicializar el almacén de prueba")
	}
	return almacen // Esto devolverá la implementación concreta que satisface la interfaz
}

// teardownTestDB cierra la base de datos y elimina el archivo.
// Ahora acepta la interfaz db.LibroAlmacenamiento
func teardownTestDB(almacen db.LibroAlmacenamiento) { // CAMBIO: de *db.AlmacenSQLite a db.LibroAlmacenamiento
	almacen.Close()
	os.Remove(testDBPath)
}

// TestAgregarLibro
func TestAgregarLibro(t *testing.T) {
	almacen := setupTestDB(t)
	defer teardownTestDB(almacen)

	libro := models.NuevoLibroConCaratula(1, "El Gran Go", "Gopher", 2023, "http://example.com/gopher.jpg")
	err := almacen.AgregarLibro(libro)
	if err != nil {
		t.Fatalf("Error al agregar libro: %v", err)
	}

	// Verificar que el libro fue agregado
	libros := almacen.ListarLibros()
	if len(libros) != 1 {
		t.Errorf("Se esperaba 1 libro, se obtuvieron %d", len(libros))
	}
	if libros[0].GetTitulo() != "El Gran Go" {
		t.Errorf("Título incorrecto: esperado 'El Gran Go', obtenido '%s'", libros[0].GetTitulo())
	}

	// Intentar agregar un libro con el mismo ID (debería fallar)
	err = almacen.AgregarLibro(libro)
	if !errors.Is(err, models.ErrLibroYaExiste) { // Uso de errors.Is para comparar errores centinela
		t.Errorf("Se esperaba un error de libro duplicado, obtenido: %v", err)
	}
}

// TestObtenerLibro
func TestObtenerLibro(t *testing.T) {
	almacen := setupTestDB(t)
	defer teardownTestDB(almacen)

	libro1 := models.NuevoLibroConCaratula(1, "Libro A", "Autor X", 2000, "urlA.jpg")
	almacen.AgregarLibro(libro1)

	// Obtener libro existente
	obtenido, err := almacen.ObtenerLibro(1)
	if err != nil {
		t.Fatalf("Error al obtener libro: %v", err)
	}
	if obtenido.GetTitulo() != "Libro A" {
		t.Errorf("Título incorrecto: esperado 'Libro A', obtenido '%s'", obtenido.GetTitulo())
	}

	// Obtener libro no existente
	_, err = almacen.ObtenerLibro(99)
	if !errors.Is(err, models.ErrLibroNoEncontrado) {
		t.Errorf("Se esperaba ErrLibroNoEncontrado, obtenido: %v", err)
	}
}

// TestListarLibros
func TestListarLibros(t *testing.T) {
	almacen := setupTestDB(t)
	defer teardownTestDB(almacen)

	// Lista vacía inicialmente
	libros := almacen.ListarLibros()
	if len(libros) != 0 {
		t.Errorf("Se esperaba lista vacía, se obtuvieron %d libros", len(libros))
	}

	libro1 := models.NuevoLibroConCaratula(1, "Libro Uno", "Autor Uno", 2020, "url1.jpg")
	libro2 := models.NuevoLibroConCaratula(2, "Libro Dos", "Autor Dos", 2021, "url2.jpg")
	almacen.AgregarLibro(libro1)
	almacen.AgregarLibro(libro2)

	libros = almacen.ListarLibros()
	if len(libros) != 2 {
		t.Errorf("Se esperaban 2 libros, se obtuvieron %d", len(libros))
	}
}

// TestActualizarLibro
func TestActualizarLibro(t *testing.T) {
	almacen := setupTestDB(t)
	defer teardownTestDB(almacen)

	libro := models.NuevoLibroConCaratula(1, "Titulo Original", "Autor Original", 2020, "original.jpg")
	almacen.AgregarLibro(libro)

	updates := map[string]interface{}{
		"titulo":       "Titulo Actualizado",
		"anio":         2022,
		"caratula_url": "updated.jpg",
	}
	err := almacen.ActualizarLibro(1, updates)
	if err != nil {
		t.Fatalf("Error al actualizar libro: %v", err)
	}

	actualizado, err := almacen.ObtenerLibro(1)
	if err != nil {
		t.Fatalf("Error al obtener libro actualizado: %v", err)
	}

	if actualizado.GetTitulo() != "Titulo Actualizado" {
		t.Errorf("Título no actualizado: esperado 'Titulo Actualizado', obtenido '%s'", actualizado.GetTitulo())
	}
	if actualizado.GetAnio() != 2022 {
		t.Errorf("Año no actualizado: esperado 2022, obtenido %d", actualizado.GetAnio())
	}
	if actualizado.GetCaratulaURL() != "updated.jpg" {
		t.Errorf("Carátula no actualizada: esperado 'updated.jpg', obtenido '%s'", actualizado.GetCaratulaURL())
	}

	// Intentar actualizar un libro no existente
	err = almacen.ActualizarLibro(99, updates)
	if !errors.Is(err, models.ErrLibroNoEncontrado) {
		t.Errorf("Se esperaba ErrLibroNoEncontrado al actualizar, obtenido: %v", err)
	}
}

// TestEliminarLibro
func TestEliminarLibro(t *testing.T) {
	almacen := setupTestDB(t)
	defer teardownTestDB(almacen)

	libro := models.NuevoLibroConCaratula(1, "Libro a Eliminar", "Autor", 2020, "url.jpg")
	almacen.AgregarLibro(libro)

	err := almacen.EliminarLibro(1)
	if err != nil {
		t.Fatalf("Error al eliminar libro: %v", err)
	}

	// Verificar que fue eliminado
	_, err = almacen.ObtenerLibro(1)
	if !errors.Is(err, models.ErrLibroNoEncontrado) {
		t.Errorf("Se esperaba ErrLibroNoEncontrado después de eliminar, obtenido: %v", err)
	}

	// Intentar eliminar un libro no existente
	err = almacen.EliminarLibro(99)
	if !errors.Is(err, models.ErrLibroNoEncontrado) {
		t.Errorf("Se esperaba ErrLibroNoEncontrado al eliminar no existente, obtenido: %v", err)
	}
}
