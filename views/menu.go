package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"libroselectronicos/db"
	"libroselectronicos/models"

	"github.com/gorilla/mux"
)

type ViewsController struct {
	Almacen   db.LibroAlmacenamiento
	Templates *template.Template
}

func NewViewsController(almacen db.LibroAlmacenamiento) *ViewsController {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error al cargar templates: %v", err)
	}
	return &ViewsController{
		Almacen:   almacen,
		Templates: tmpl,
	}
}

func (vc *ViewsController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if err := vc.Templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Error al cargar la página principal", http.StatusInternalServerError)
	}
}

func (vc *ViewsController) ListarLibrosHTML(w http.ResponseWriter, r *http.Request) {
	libros := vc.Almacen.ListarLibros()
	if err := vc.Templates.ExecuteTemplate(w, "listar.html", libros); err != nil {
		http.Error(w, "Error al cargar la lista de libros", http.StatusInternalServerError)
	}
}

func (vc *ViewsController) CrearLibroHTMLForm(w http.ResponseWriter, r *http.Request) {
	if err := vc.Templates.ExecuteTemplate(w, "crear.html", nil); err != nil {
		http.Error(w, "Error al cargar el formulario de creación", http.StatusInternalServerError)
	}
}

func (vc *ViewsController) CrearLibroHTMLSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	titulo := r.FormValue("titulo")
	autor := r.FormValue("autor")
	anioStr := r.FormValue("anio")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	anio, err := strconv.Atoi(anioStr)
	if err != nil {
		http.Error(w, "Año inválido", http.StatusBadRequest)
		return
	}

	nuevoLibro := models.NuevoLibro(id, titulo, autor, anio)
	err = vc.Almacen.AgregarLibro(nuevoLibro)
	if err != nil {
		http.Redirect(w, r, "/libros/crear?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

func (vc *ViewsController) EditarLibroHTMLForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	libro, err := vc.Almacen.ObtenerLibro(id)
	if err != nil {
		if err == models.ErrLibroNoEncontrado {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener libro para edición: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := vc.Templates.ExecuteTemplate(w, "editar.html", libro); err != nil {
		http.Error(w, "Error al cargar el formulario de edición", http.StatusInternalServerError)
	}
}

func (vc *ViewsController) EditarLibroHTMLSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}

	updates := make(map[string]interface{})
	if titulo := r.FormValue("titulo"); titulo != "" {
		updates["titulo"] = titulo
	}
	if autor := r.FormValue("autor"); autor != "" {
		updates["autor"] = autor
	}
	if anioStr := r.FormValue("anio"); anioStr != "" {
		if anio, err := strconv.Atoi(anioStr); err == nil {
			updates["anio"] = anio
		}
	}

	if len(updates) == 0 {
		http.Error(w, "No hay datos para actualizar", http.StatusBadRequest)
		return
	}

	err = vc.Almacen.ActualizarLibro(id, updates)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/libros/%d/editar?error=%s", id, err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

func (vc *ViewsController) EliminarLibroHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	err = vc.Almacen.EliminarLibro(id)
	if err != nil {
		if err == models.ErrLibroNoEncontrado {
			http.Redirect(w, r, "/libros?error=Libro no encontrado para eliminar", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/libros?error="+err.Error(), http.StatusSeeOther)
		}
		return
	}

	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}
