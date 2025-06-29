package db_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "libroselectronicos/controllers" // El paquete real del controlador
	"libroselectronicos/models"                  // El paquete models

	"github.com/gorilla/mux"
)

// MockApiAlmacen es una implementación mock de db.LibroAlmacenamiento para pruebas de API.
// Este mock DEBE implementar db.LibroAlmacenamiento exactamente como está definida en db/storage.go
type MockApiAlmacen struct {
	// CAMBIO CLAVE AQUÍ: La firma de MockAgregarLibro ahora es *models.Libro
	MockAgregarLibro    func(libro *models.Libro) error
	MockListarLibros    func() []*models.Libro
	MockObtenerLibro    func(id int) (*models.Libro, error)
	MockActualizarLibro func(id int, updates map[string]interface{}) error
	MockEliminarLibro   func(id int) error
	MockClose           func() error
}

// Implementación de los métodos de la interfaz db.LibroAlmacenamiento
// CAMBIO CLAVE AQUÍ: La firma de AgregarLibro ahora es *models.Libro
func (m *MockApiAlmacen) AgregarLibro(libro *models.Libro) error {
	if m.MockAgregarLibro != nil {
		return m.MockAgregarLibro(libro)
	}
	return errors.New("AgregarLibro no implementado en mock")
}

func (m *MockApiAlmacen) ListarLibros() []*models.Libro {
	if m.MockListarLibros != nil {
		return m.MockListarLibros()
	}
	return []*models.Libro{}
}

func (m *MockApiAlmacen) ObtenerLibro(id int) (*models.Libro, error) {
	if m.MockObtenerLibro != nil {
		return m.MockObtenerLibro(id)
	}
	return nil, models.ErrLibroNoEncontrado
}

func (m *MockApiAlmacen) ActualizarLibro(id int, updates map[string]interface{}) error {
	if m.MockActualizarLibro != nil {
		return m.MockActualizarLibro(id, updates)
	}
	return errors.New("ActualizarLibro no implementado en mock")
}

func (m *MockApiAlmacen) EliminarLibro(id int) error {
	if m.MockEliminarLibro != nil {
		return m.MockEliminarLibro(id)
	}
	return errors.New("EliminarLibro no implementado en mock")
}

func (m *MockApiAlmacen) Close() error {
	if m.MockClose != nil {
		return m.MockClose()
	}
	return nil
}

// TestGetLibrosAPI prueba la ruta GET /api/libros
func TestGetLibrosAPI(t *testing.T) {
	mockAlmacen := &MockApiAlmacen{
		MockListarLibros: func() []*models.Libro {
			return []*models.Libro{
				models.NuevoLibroConCaratula(1, "Libro Uno", "Autor A", 2000, "url1.jpg"),
				models.NuevoLibroConCaratula(2, "Libro Dos", "Autor B", 2005, "url2.jpg"),
			}
		},
	}
	controller := controllers.NewLibroApiController(mockAlmacen)

	req, err := http.NewRequest("GET", "/api/libros", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetLibrosAPI)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Manejador devolvió código de estado incorrecto: esperado %d, obtenido %d", http.StatusOK, status)
	}

	// Corrección para el \n
	expectedBody := `[{"id":1,"titulo":"Libro Uno","autor":"Autor A","anio":2000,"caratula_url":"url1.jpg"},{"id":2,"titulo":"Libro Dos","autor":"Autor B","anio":2005,"caratula_url":"url2.jpg"}]` + "\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Manejador devolvió cuerpo incorrecto: esperado %s, obtenido %s", expectedBody, rr.Body.String())
	}
}

// TestGetLibroByIDAPI prueba la ruta GET /api/libros/{id}
func TestGetLibroByIDAPI(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockBook       *models.Libro
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Libro encontrado",
			id:             "1",
			mockBook:       models.NuevoLibroConCaratula(1, "Libro de Prueba", "Autor Prueba", 2020, "url.jpg"),
			mockError:      nil,
			expectedStatus: http.StatusOK,
			// Corrección para el \n
			expectedBody: `{"id":1,"titulo":"Libro de Prueba","autor":"Autor Prueba","anio":2020,"caratula_url":"url.jpg"}` + "\n",
		},
		{
			name:           "Libro no encontrado",
			id:             "99",
			mockBook:       nil,
			mockError:      models.ErrLibroNoEncontrado,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Libro no encontrado\n",
		},
		{
			name:           "ID inválido",
			id:             "abc",
			mockBook:       nil,
			mockError:      nil, // No se llama al mock si el ID es inválido
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "ID de libro inválido\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAlmacen := &MockApiAlmacen{
				// Lógica corregida para el mock:
				MockObtenerLibro: func(id int) (*models.Libro, error) {
					if tt.mockError != nil { // Si se espera un error (como no encontrado)
						return nil, tt.mockError
					}
					// Si no hay error, entonces se espera un libro.
					// Asegúrate de que el ID solicitado coincida con el libro mock
					if tt.mockBook != nil && id == tt.mockBook.ID {
						return tt.mockBook, nil // No hay error aquí
					}
					// Esto es un fallback, en teoría no debería alcanzarse si los tests están bien definidos.
					return nil, models.ErrLibroNoEncontrado
				},
			}
			controller := controllers.NewLibroApiController(mockAlmacen)

			req, err := http.NewRequest("GET", "/api/libros/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/api/libros/{id}", controller.GetLibroByIDAPI).Methods("GET")
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Manejador devolvió código de estado incorrecto: esperado %d, obtenido %d. Cuerpo: %s", tt.expectedStatus, status, rr.Body.String())
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("Manejador devolvió cuerpo incorrecto: esperado %q, obtenido %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

// TestCreateLibroAPI prueba la ruta POST /api/libros
func TestCreateLibroAPI(t *testing.T) {
	tests := []struct {
		name           string
		inputJSON      string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Creación exitosa",
			inputJSON:      `{"id":1,"titulo":"Nuevo Libro","autor":"Nuevo Autor","anio":2024,"caratula_url":"new_url.jpg"}`,
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			// Corrección para el \n AQUÍ
			expectedBody: `{"id":1,"titulo":"Nuevo Libro","autor":"Nuevo Autor","anio":2024,"caratula_url":"new_url.jpg"}` + "\n",
		},
		{
			name:           "ID duplicado",
			inputJSON:      `{"id":1,"titulo":"Libro Duplicado","autor":"Autor Duplicado","anio":2020,"caratula_url":"dup_url.jpg"}`,
			mockError:      models.ErrLibroYaExiste, // Ahora esperamos este error centinela
			expectedStatus: http.StatusConflict,     // Esperamos 409 Conflict
			expectedBody:   "libro con este ID ya existe\n",
		},
		{
			name:           "JSON inválido",
			inputJSON:      `{"id":1,"titulo":"Nuevo Libro"`, // JSON mal formado
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Solicitud JSON inválida: unexpected EOF\n",
		},
		{
			name:           "ID es 0",
			inputJSON:      `{"id":0,"titulo":"Libro ID 0","autor":"Autor","anio":2020,"caratula_url":"url.jpg"}`,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "El ID del libro no puede ser 0\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAlmacen := &MockApiAlmacen{
				// CAMBIO CLAVE AQUÍ: La función del mock acepta *models.Libro
				MockAgregarLibro: func(libro *models.Libro) error {
					return tt.mockError
				},
			}
			controller := controllers.NewApiLibroController(mockAlmacen)

			req, err := http.NewRequest("POST", "/api/libros", bytes.NewBufferString(tt.inputJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.CreateLibroAPI)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Manejador devolvió código de estado incorrecto: esperado %d, obtenido %d. Cuerpo: %s", tt.expectedStatus, status, rr.Body.String())
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("Manejador devolvió cuerpo incorrecto: esperado %q, obtenido %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

// TestUpdateLibroAPI prueba la ruta PUT /api/libros/{id}
func TestUpdateLibroAPI(t *testing.T) {
	tests := []struct {
		name                  string
		id                    string
		inputJSON             string
		mockUpdateErr         error
		mockGetAfterUpdate    *models.Libro
		mockGetAfterUpdateErr error
		expectedStatus        int
		expectedBody          string
	}{
		{
			name:                  "Actualización exitosa",
			id:                    "1",
			inputJSON:             `{"titulo":"Título Actualizado","anio":2025,"caratula_url":"updated_url.jpg"}`,
			mockUpdateErr:         nil,
			mockGetAfterUpdate:    models.NuevoLibroConCaratula(1, "Título Actualizado", "Autor Existente", 2025, "updated_url.jpg"),
			mockGetAfterUpdateErr: nil,
			expectedStatus:        http.StatusOK,
			// Corrección para el \n AQUÍ
			expectedBody: `{"id":1,"titulo":"Título Actualizado","autor":"Autor Existente","anio":2025,"caratula_url":"updated_url.jpg"}` + "\n",
		},
		{
			name:                  "Libro no encontrado para actualizar",
			id:                    "99",
			inputJSON:             `{"titulo":"Inexistente"}`,
			mockUpdateErr:         models.ErrLibroNoEncontrado,
			mockGetAfterUpdate:    nil,
			mockGetAfterUpdateErr: models.ErrLibroNoEncontrado,
			expectedStatus:        http.StatusNotFound,
			expectedBody:          "Libro no encontrado\n",
		},
		{
			name:                  "ID inválido",
			id:                    "abc",
			inputJSON:             `{"titulo":"Test"}`,
			mockUpdateErr:         nil,
			mockGetAfterUpdate:    nil,
			mockGetAfterUpdateErr: nil,
			expectedStatus:        http.StatusBadRequest,
			expectedBody:          "ID de libro inválido\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAlmacen := &MockApiAlmacen{
				MockActualizarLibro: func(id int, updates map[string]interface{}) error {
					return tt.mockUpdateErr
				},
				MockObtenerLibro: func(id int) (*models.Libro, error) {
					return tt.mockGetAfterUpdate, tt.mockGetAfterUpdateErr
				},
			}
			controller := controllers.NewApiLibroController(mockAlmacen)

			req, err := http.NewRequest("PUT", "/api/libros/"+tt.id, bytes.NewBufferString(tt.inputJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/api/libros/{id}", controller.UpdateLibroAPI).Methods("PUT")
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Manejador devolvió código de estado incorrecto: esperado %d, obtenido %d. Cuerpo: %s", tt.expectedStatus, status, rr.Body.String())
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("Manejador devolvió cuerpo incorrecto: esperado %q, obtenido %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

// TestDeleteLibroAPI prueba la ruta DELETE /api/libros/{id}
func TestDeleteLibroAPI(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Eliminación exitosa",
			id:             "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent, // Esperamos 204 No Content
			expectedBody:   "",                   // No hay cuerpo para 204
		},
		{
			name:           "Libro no encontrado para eliminar",
			id:             "99",
			mockError:      models.ErrLibroNoEncontrado,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Libro no encontrado\n",
		},
		{
			name:           "ID inválido",
			id:             "abc",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "ID de libro inválido\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAlmacen := &MockApiAlmacen{
				MockEliminarLibro: func(id int) error {
					return tt.mockError
				},
			}
			controller := controllers.NewApiLibroController(mockAlmacen)

			req, err := http.NewRequest("DELETE", "/api/libros/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/api/libros/{id}", controller.DeleteLibroAPI).Methods("DELETE")
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Manejador devolvió código de estado incorrecto: esperado %d, obtenido %d. Cuerpo: %s", tt.expectedStatus, status, rr.Body.String())
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("Manejador devolvió cuerpo incorrecto: esperado %q, obtenido %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
