package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings" // Asegúrate de que strings esté importado si se usa

	"libroselectronicos/db"
	"libroselectronicos/models"

	"github.com/gorilla/mux"
)

// ApiLibroController DEBE EMPEZAR CON MAYÚSCULA 'A' para ser exportable y visible.
type ApiLibroController struct {
	Almacen db.LibroAlmacenamiento
}

// NewApiLibroController DEBE EMPEZAR CON MAYÚSCULA 'N' y 'A' para ser exportable y visible.
func NewApiLibroController(almacen db.LibroAlmacenamiento) *ApiLibroController {
	return &ApiLibroController{Almacen: almacen}
}

// CreateLibroAPI maneja la solicitud HTTP POST para crear un nuevo libro.
func (ac *ApiLibroController) CreateLibroAPI(w http.ResponseWriter, r *http.Request) {
	var nuevoLibro models.Libro // Deserializar a la struct concreta Libro
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &nuevoLibro); err != nil {
		http.Error(w, "Error al decodificar el libro: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Asegurarse de que el ID sea proporcionado para la creación con ID específico
	if nuevoLibro.ID == 0 {
		http.Error(w, "El ID del libro es requerido para la creación.", http.StatusBadRequest)
		return
	}

	// Agregar el libro usando la interfaz
	libroParaAlmacenar := models.NuevoLibro(nuevoLibro.ID, nuevoLibro.Titulo, nuevoLibro.Autor, nuevoLibro.Anio)
	err = ac.Almacen.AgregarLibro(libroParaAlmacenar)
	if err != nil {
		if strings.Contains(err.Error(), "ya existe") {
			http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusConflict)
		} else {
			http.Error(w, "Error al agregar libro: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(libroParaAlmacenar) // Retornar el libro agregado
}

// GetLibrosAPI maneja la solicitud HTTP GET para listar todos los libros.
func (ac *ApiLibroController) GetLibrosAPI(w http.ResponseWriter, r *http.Request) {
	libros := ac.Almacen.ListarLibros()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libros)
}

// GetLibroByIDAPI maneja la solicitud HTTP GET para obtener un libro por su ID.
func (ac *ApiLibroController) GetLibroByIDAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	libro, err := ac.Almacen.ObtenerLibro(id)
	if err != nil {
		if err == models.ErrLibroNoEncontrado {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener libro: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libro)
}

// UpdateLibroAPI maneja la solicitud HTTP PUT para actualizar un libro.
func (ac *ApiLibroController) UpdateLibroAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &updates); err != nil {
		http.Error(w, "Error al decodificar los datos de actualización: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = ac.Almacen.ActualizarLibro(id, updates)
	if err != nil {
		if err == models.ErrLibroNoEncontrado {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar libro: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Libro con ID %d actualizado exitosamente", id)
}

// DeleteLibroAPI maneja la solicitud HTTP DELETE para eliminar un libro.
func (ac *ApiLibroController) DeleteLibroAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	err = ac.Almacen.EliminarLibro(id)
	if err != nil {
		if err == models.ErrLibroNoEncontrado {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar libro: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Libro con ID %d eliminado exitosamente", id)
}
