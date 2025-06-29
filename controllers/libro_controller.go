package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"libroselectronicos/db"
	"libroselectronicos/models"

	"github.com/gorilla/mux"
)

type LibroController struct {
	Almacen db.LibroAlmacenamiento
}

func NuevoLibroController(almacen db.LibroAlmacenamiento) *LibroController {
	return &LibroController{Almacen: almacen}
}

func (lc *LibroController) CrearLibro(w http.ResponseWriter, r *http.Request) {
	var nuevoLibro models.Libro // Aquí recibimos la estructura completa, incluyendo CaratulaURL
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

	// ¡CORRECCIÓN AQUÍ!
	// Utiliza NuevoLibroConCaratula y pasa todos los campos, incluyendo CaratulaURL
	// La función AgregarLibro de la interfaz espera *models.Libro
	libroParaAlmacenar := models.NuevoLibroConCaratula(
		nuevoLibro.ID,
		nuevoLibro.Titulo,
		nuevoLibro.Autor,
		nuevoLibro.Anio,
		nuevoLibro.CaratulaURL, // Añadido CaratulaURL
	)

	err = lc.Almacen.AgregarLibro(libroParaAlmacenar)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") { // Cambio a un error más específico de SQLite
			http.Error(w, fmt.Sprintf("Error: El libro con ID %d ya existe.", nuevoLibro.ID), http.StatusConflict)
		} else {
			http.Error(w, "Error al agregar libro: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(libroParaAlmacenar)
}

func (lc *LibroController) ObtenerLibros(w http.ResponseWriter, r *http.Request) {
	libros := lc.Almacen.ListarLibros()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libros)
}

func (lc *LibroController) ObtenerLibroPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	libro, err := lc.Almacen.ObtenerLibro(id)
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

func (lc *LibroController) ActualizarLibro(w http.ResponseWriter, r *http.Request) {
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

	err = lc.Almacen.ActualizarLibro(id, updates)
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

func (lc *LibroController) EliminarLibro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	err = lc.Almacen.EliminarLibro(id)
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
