package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"libroselectronicos/controllers"
	"libroselectronicos/db"
	"libroselectronicos/models"

	"github.com/gorilla/mux"
)

// MockApiAlmacen es una implementación de mock de db.LibroAlmacenamiento para pruebas.
type MockApiAlmacen struct {
	Libros map[int]models.LibroInterface
	Err    error // Para simular errores de la base de datos
}

var _ db.LibroAlmacenamiento = (*MockApiAlmacen)(nil)

func NewMockApiAlmacen(initialBooks ...models.LibroInterface) *MockApiAlmacen {
	m := &MockApiAlmacen{
		Libros: make(map[int]models.LibroInterface),
	}
	for _, book := range initialBooks {
		m.Libros[book.GetID()] = book
	}
	return m
}

func (m *MockApiAlmacen) AgregarLibro(libro models.LibroInterface) error {
	if m.Err != nil {
		return m.Err
	}
	if _, exists := m.Libros[libro.GetID()]; exists {
		return fmt.Errorf("libro con ID %d ya existe", libro.GetID())
	}
	m.Libros[libro.GetID()] = libro
	return nil
}

func (m *MockApiAlmacen) ListarLibros() []models.LibroInterface {
	if m.Err != nil {
		return nil
	}
	var libros []models.LibroInterface
	for _, libro := range m.Libros {
		libros = append(libros, libro)
	}
	return libros
}

func (m *MockApiAlmacen) ObtenerLibro(id int) (models.LibroInterface, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if libro, ok := m.Libros[id]; ok {
		return libro, nil
	}
	return nil, models.ErrLibroNoEncontrado
}

func (m *MockApiAlmacen) ActualizarLibro(id int, updates map[string]interface{}) error {
	if m.Err != nil {
		return m.Err
	}
	libro, ok := m.Libros[id]
	if !ok {
		return models.ErrLibroNoEncontrado
	}

	if titulo, ok := updates["titulo"].(string); ok {
		libro.SetTitulo(titulo)
	}
	if autor, ok := updates["autor"].(string); ok {
		libro.SetAutor(autor)
	}
	if anioFloat, ok := updates["anio"].(float64); ok {
		libro.SetAnio(int(anioFloat))
	} else if anioInt, ok := updates["anio"].(int); ok { // Por si viene como int
		libro.SetAnio(anioInt)
	}
	m.Libros[id] = libro
	return nil
}

func (m *MockApiAlmacen) EliminarLibro(id int) error {
	if m.Err != nil {
		return m.Err
	}
	if _, ok := m.Libros[id]; !ok {
		return models.ErrLibroNoEncontrado
	}
	delete(m.Libros, id)
	return nil
}

func (m *MockApiAlmacen) Close() error {
	return nil
}

// ClearTable implementa el método ClearTable de la interfaz db.LibroAlmacenamiento.
func (m *MockApiAlmacen) ClearTable() error {
	m.Libros = make(map[int]models.LibroInterface)
	return nil
}

// --- Pruebas para ApiLibroController ---

func TestCreateLibroAPI(t *testing.T) {
	mockAlmacen := NewMockApiAlmacen()
	apiController := controllers.NewApiLibroController(mockAlmacen)

	router := mux.NewRouter()
	router.HandleFunc("/api/libros", apiController.CreateLibroAPI).Methods("POST")

	libroJSON := `{"id":1, "titulo":"Libro de Prueba", "autor":"Autor Test", "anio":2024}`
	req, err := http.NewRequest("POST", "/api/libros", bytes.NewBufferString(libroJSON))
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Manejador devolvió código de estado incorrecto: esperado %v, obtenido %v. Cuerpo: %s",
			http.StatusCreated, status, rr.Body.String())
	}

	var fetchedLibro models.Libro // Deserializar a la struct concreta Libro
	err = json.Unmarshal(rr.Body.Bytes(), &fetchedLibro)
	if err != nil {
		t.Fatalf("No se pudo deserializar la respuesta JSON: %v", err)
	}

	// Accede directamente a los campos de la struct Libro
	if fetchedLibro.ID != 1 {
		t.Errorf("Datos del libro incorrectos: esperado ID 1, obtenido ID %d", fetchedLibro.ID)
	}
	if fetchedLibro.Titulo != "Libro de Prueba" {
		t.Errorf("Datos del libro incorrectos: esperado Título 'Libro de Prueba', obtenido '%s'", fetchedLibro.Titulo)
	}

	// Intentar crear un libro con ID existente
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Manejador devolvió código de estado incorrecto para duplicado: esperado %v, obtenido %v. Cuerpo: %s",
			http.StatusConflict, status, rr.Body.String())
	}
}

func TestGetLibrosAPI(t *testing.T) {
	initialBooks := []models.LibroInterface{
		models.NuevoLibro(1, "Libro A", "Autor X", 2000),
		models.NuevoLibro(2, "Libro B", "Autor Y", 2001),
	}
	mockAlmacen := NewMockApiAlmacen(initialBooks...)
	apiController := controllers.NewApiLibroController(mockAlmacen)

	router := mux.NewRouter()
	router.HandleFunc("/api/libros", apiController.GetLibrosAPI).Methods("GET")

	req, err := http.NewRequest("GET", "/api/libros", nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Manejador devolvió código de estado incorrecto: esperado %v, obtenido %v. Cuerpo: %s",
			http.StatusOK, status, rr.Body.String())
	}

	var libros []models.Libro // Deserializar a un slice de la struct concreta Libro
	err = json.Unmarshal(rr.Body.Bytes(), &libros)
	if err != nil {
		t.Fatalf("No se pudo deserializar la respuesta JSON: %v", err)
	}

	if len(libros) != 2 {
		t.Errorf("Se esperaban 2 libros, obtenidos %d. Cuerpo: %s", len(libros), rr.Body.String())
	}
	// Asegúrate de que los IDs y títulos son correctos
	// Como el orden de un map no está garantizado, podrías ordenar `libros` por ID si el orden es importante,
	// o iterar sobre ambos slices para verificar la existencia.
	// Para simplicidad, asumo un orden si los agregaste en orden y el mock los devuelve así.
	// Sin embargo, una forma más robusta es verificar cada libro individualmente.
	expectedBooks := map[int]struct {
		Title  string
		Author string
		Year   int
	}{
		1: {Title: "Libro A", Author: "Autor X", Year: 2000},
		2: {Title: "Libro B", Author: "Autor Y", Year: 2001},
	}

	if len(libros) != len(expectedBooks) {
		t.Fatalf("Número de libros inesperado. Esperado %d, obtenido %d", len(expectedBooks), len(libros))
	}

	for _, libro := range libros {
		expected, ok := expectedBooks[libro.ID]
		if !ok {
			t.Errorf("Libro con ID %d inesperado", libro.ID)
		}
		if libro.Titulo != expected.Title {
			t.Errorf("Titulo incorrecto para ID %d. Esperado '%s', obtenido '%s'", libro.ID, expected.Title, libro.Titulo)
		}
		if libro.Autor != expected.Author {
			t.Errorf("Autor incorrecto para ID %d. Esperado '%s', obtenido '%s'", libro.ID, expected.Author, libro.Autor)
		}
		if libro.Anio != expected.Year {
			t.Errorf("Año incorrecto para ID %d. Esperado %d, obtenido %d", libro.ID, expected.Year, libro.Anio)
		}
	}
}

func TestGetLibroByIDAPI(t *testing.T) {
	initialBooks := []models.LibroInterface{
		models.NuevoLibro(5, "Libro Buscado", "Autor Z", 1990),
	}
	mockAlmacen := NewMockApiAlmacen(initialBooks...)
	apiController := controllers.NewApiLibroController(mockAlmacen)

	router := mux.NewRouter()
	router.HandleFunc("/api/libros/{id}", apiController.GetLibroByIDAPI).Methods("GET")

	// Prueba 1: Obtener libro existente
	req, err := http.NewRequest("GET", "/api/libros/5", nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Manejador devolvió código de estado incorrecto: esperado %v, obtenido %v. Cuerpo: %s",
			http.StatusOK, status, rr.Body.String())
	}

	var fetchedLibro models.Libro // Usar Libro directamente para deserializar
	err = json.Unmarshal(rr.Body.Bytes(), &fetchedLibro)
	if err != nil {
		t.Fatalf("No se pudo deserializar la respuesta JSON: %v", err)
	}
	// Acceder directamente a los campos de la struct Libro
	if fetchedLibro.ID != 5 { // CORRECCIÓN: Verifica solo el ID primero
		t.Errorf("Datos del libro incorrectos: esperado ID 5, obtenido ID %d", fetchedLibro.ID)
	}
	if fetchedLibro.Titulo != "Libro Buscado" { // CORRECCIÓN: Verifica el título por separado
		t.Errorf("Datos del libro incorrectos: esperado Título 'Libro Buscado', obtenido Título '%s'", fetchedLibro.Titulo)
	}

	// Prueba 2: Obtener un libro que no existe
	req2, err := http.NewRequest("GET", "/api/libros/999", nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusNotFound {
		t.Errorf("Manejador devolvió código de estado incorrecto para no encontrado: esperado %v, obtenido %v. Cuerpo: %s",
			http.StatusNotFound, status, rr2.Body.String())
	}
}

func TestUpdateLibroAPI(t *testing.T) {
	initialBook := models.NuevoLibro(1, "Original", "Autor Original", 2000)
	mockAlmacen := NewMockApiAlmacen(initialBook)
	apiController := controllers.NewApiLibroController(mockAlmacen)

	router := mux.NewRouter()
	router.HandleFunc("/api/libros/{id}", apiController.UpdateLibroAPI).Methods("PUT")

	// Prueba 1: Actualizar libro existente
	updatesJSON := `{"titulo":"Titulo Actualizado", "anio":2025}`
	req, err := http.NewRequest("PUT", "/api/libros/1", bytes.NewBufferString(updatesJSON))
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Manejador devolvió código de estado incorrecto: esperado %v, obtenido %v. Cuerpo: %s",
			http.StatusOK, status, rr.Body.String())
	}

	// Verificar en el mock
	updatedBook, _ := mockAlmacen.ObtenerLibro(1)
	if updatedBook.GetTitulo() != "Titulo Actualizado" || updatedBook.GetAnio() != 2025 {
		t.Errorf("Libro no actualizado correctamente en el mock. Titulo: %s, Anio: %d",
			updatedBook.GetTitulo(), updatedBook.GetAnio())
	}

	// Prueba 2: Intentar actualizar libro no existente
	req2, err := http.NewRequest("PUT", "/api/libros/999", bytes.NewBufferString(updatesJSON))
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}
	req2.Header.Set("Content-Type", "application/json")

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusNotFound {
		t.Errorf("Manejador devolvió código de estado incorrecto para no encontrado: esperado %v, obtenido %v. Cuerpo: %s",
			http.StatusNotFound, status, rr2.Body.String())
	}
}

func TestDeleteLibroAPI(t *testing.T) {
	initialBook := models.NuevoLibro(1, "Libro a Eliminar", "Autor", 2000)
	mockAlmacen := NewMockApiAlmacen(initialBook)
	apiController := controllers.NewApiLibroController(mockAlmacen)

	router := mux.NewRouter()
	router.HandleFunc("/api/libros/{id}", apiController.DeleteLibroAPI).Methods("DELETE")

	// Prueba 1: Eliminar libro existente
	req, err := http.NewRequest("DELETE", "/api/libros/1", nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Manejador devolvió código de estado incorrecto: esperado %v, obtenido %v. Cuerpo: %s",
			http.StatusOK, status, rr.Body.String())
	}

	// Verificar que el libro fue eliminado del mock
	_, err = mockAlmacen.ObtenerLibro(1)
	if err == nil || err != models.ErrLibroNoEncontrado {
		t.Errorf("Libro no eliminado, se esperaba ErrLibroNoEncontrado")
	}

	// Prueba 2: Intentar eliminar libro no existente
	req2, err := http.NewRequest("DELETE", "/api/libros/999", nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusNotFound {
		t.Errorf("Manejador devolvió código de estado incorrecto para no encontrado: esperado %v, obtenido %v. Cuerpo: %s",
			http.StatusNotFound, status, rr2.Body.String())
	}
}
