package views

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"

	// Para manejar fechas en el futuro (alquileres)
	"libroselectronicos/db"
	"libroselectronicos/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt" // Para hashear contraseñas
)

const (
	sessionName = "session-name"
	// Clave secreta para las sesiones. ¡DEBE CAMBIARSE EN PRODUCCIÓN!
	// Para desarrollo, una cadena aleatoria está bien.
	sessionKey = "super-secret-key-that-should-be-long-and-random"
)

var store = sessions.NewCookieStore([]byte(sessionKey))

type MenuController struct {
	almacen     db.LibroAlmacenamiento
	indexTpl    templateExecutor
	listTpl     templateExecutor
	createTpl   templateExecutor
	editTpl     templateExecutor
	sinopsisTpl templateExecutor
	registerTpl templateExecutor // Nueva plantilla para registro
	loginTpl    templateExecutor // Nueva plantilla para login
}

type templateExecutor interface {
	Execute(wr http.ResponseWriter, data interface{}) error
}

type htmlTemplateWrapper struct {
	tpl *template.Template
}

func (w *htmlTemplateWrapper) Execute(wr http.ResponseWriter, data interface{}) error {
	return w.tpl.Execute(wr, data)
}

func NewMenuController(almacen db.LibroAlmacenamiento) *MenuController {
	return &MenuController{
		almacen:     almacen,
		indexTpl:    &htmlTemplateWrapper{template.Must(template.ParseFiles("templates/index.html"))},
		listTpl:     &htmlTemplateWrapper{template.Must(template.ParseFiles("templates/listar.html"))},
		createTpl:   &htmlTemplateWrapper{template.Must(template.ParseFiles("templates/crear.html"))},
		editTpl:     &htmlTemplateWrapper{template.Must(template.ParseFiles("templates/editar.html"))},
		sinopsisTpl: &htmlTemplateWrapper{template.Must(template.ParseFiles("templates/sinopsis.html"))},
		registerTpl: &htmlTemplateWrapper{template.Must(template.ParseFiles("templates/registro.html"))},
		loginTpl:    &htmlTemplateWrapper{template.Must(template.ParseFiles("templates/login.html"))},
	}
}

// Estructura para pasar datos a las plantillas que necesitan información del usuario
type TemplateData struct {
	Libros  []*models.Libro
	Usuario *models.Usuario // nil si no está logueado
	Error   string
}

// Helper para obtener el usuario logueado
func (vc *MenuController) getLoggedInUser(r *http.Request) *models.Usuario {
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Printf("Error al obtener sesión: %v", err)
		return nil
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok || userID == 0 {
		return nil
	}

	user, err := vc.almacen.ObtenerUsuarioPorID(userID)
	if err != nil {
		log.Printf("Error al obtener usuario de la DB: %v", err)
		return nil
	}
	return user
}

// Index muestra la página principal.
func (vc *MenuController) Index(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Usuario: vc.getLoggedInUser(r),
	}
	err := vc.indexTpl.Execute(w, data)
	if err != nil {
		log.Printf("Error al renderizar plantilla index.html: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// ListarLibrosHTML lista todos los libros en HTML y muestra el usuario logueado.
func (vc *MenuController) ListarLibrosHTML(w http.ResponseWriter, r *http.Request) {
	libros := vc.almacen.ListarLibros()
	data := TemplateData{
		Libros:  libros,
		Usuario: vc.getLoggedInUser(r),
	}
	err := vc.listTpl.Execute(w, data)
	if err != nil {
		log.Printf("Error al renderizar plantilla listar.html: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// --- Manajadores de Autenticación ---

// RegistrarUsuarioHTML muestra el formulario de registro.
func (vc *MenuController) RegistrarUsuarioHTML(w http.ResponseWriter, r *http.Request) {
	err := vc.registerTpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error al renderizar plantilla registro.html: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// RegistrarUsuarioSubmit maneja el envío del formulario de registro.
func (vc *MenuController) RegistrarUsuarioSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	rol := r.FormValue("rol") // Podrías tener un selector o dejarlo por defecto 'lector'

	if username == "" || password == "" {
		http.Error(w, "Nombre de usuario y contraseña son requeridos.", http.StatusBadRequest)
		return
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al hashear contraseña: %v", err)
		http.Error(w, "Error interno del servidor al procesar contraseña", http.StatusInternalServerError)
		return
	}

	// rol por defecto si no se selecciona (ej. para admins)
	if rol == "" {
		rol = "lector"
	}
	nuevoUsuario := models.NuevoUsuario(0, username, string(hashedPassword), email, rol)

	// --- INICIO DE NUEVOS LOGS PARA DEBUGGING ---
	log.Printf("DEBUG: Intentando agregar usuario: %s (Rol: %s)", nuevoUsuario.GetUsername(), nuevoUsuario.GetRol())
	// --- FIN DE NUEVOS LOGS PARA DEBUGGING ---

	err = vc.almacen.AgregarUsuario(nuevoUsuario)
	if err != nil {
		// --- LOG MEJORADO ---
		log.Printf("ERROR: Fallo al agregar usuario %s a la base de datos: %v", nuevoUsuario.GetUsername(), err)
		// --- FIN LOG MEJORADO ---
		if errors.Is(err, models.ErrUsuarioYaExiste) {
			http.Error(w, "El nombre de usuario ya está en uso.", http.StatusConflict)
		} else {
			log.Printf("Error al registrar usuario: %v", err) // Este ya lo tenías
			http.Error(w, "Error interno del servidor al registrar usuario", http.StatusInternalServerError)
		}
		return
	}
	// --- NUEVO LOG PARA ÉXITO ---
	log.Printf("DEBUG: Usuario %s agregado exitosamente.", nuevoUsuario.GetUsername())
	// --- FIN NUEVO LOG PARA ÉXITO ---

	http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirigir al login después del registro exitoso
}

// LoginHTML muestra el formulario de inicio de sesión.
func (vc *MenuController) LoginHTML(w http.ResponseWriter, r *http.Request) {
	err := vc.loginTpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error al renderizar plantilla login.html: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// LoginSubmit maneja el envío del formulario de inicio de sesión.
func (vc *MenuController) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Nombre de usuario y contraseña son requeridos.", http.StatusBadRequest)
		return
	}

	usuario, err := vc.almacen.ObtenerUsuarioPorUsername(username)
	if err != nil {
		if errors.Is(err, models.ErrUsuarioNoEncontrado) {
			http.Error(w, "Usuario o contraseña incorrectos.", http.StatusUnauthorized)
		} else {
			log.Printf("Error al obtener usuario para login: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	// Comparar la contraseña ingresada con la contraseña hasheada
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(password))
	if err != nil {
		http.Error(w, "Usuario o contraseña incorrectos.", http.StatusUnauthorized)
		return
	}

	// Iniciar sesión (establecer cookie de sesión)
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Printf("Error al obtener sesión para login: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
	session.Values["user_id"] = usuario.ID
	session.Values["username"] = usuario.Username // Guardamos también el username para conveniencia
	session.Save(r, w)                            // Guardar la sesión

	http.Redirect(w, r, "/libros", http.StatusSeeOther) // Redirigir a la lista de libros
}

// Logout cierra la sesión del usuario.
func (vc *MenuController) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Printf("Error al obtener sesión para logout: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1 // Expira la cookie de sesión
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther) // Redirigir a la página de inicio
}

// CrearLibroHTML muestra el formulario para crear un nuevo libro.
func (vc *MenuController) CrearLibroHTML(w http.ResponseWriter, r *http.Request) {
	err := vc.createTpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error al renderizar plantilla crear.html: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// CrearLibroHTMLSubmit maneja el envío del formulario para crear un nuevo libro.
func (vc *MenuController) CrearLibroHTMLSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	titulo := r.FormValue("titulo")
	autor := r.FormValue("autor")
	anioStr := r.FormValue("anio")
	caratulaURL := r.FormValue("caratula_url")
	sinopsis := r.FormValue("sinopsis") // Captura la sinopsis

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	anio, err := strconv.Atoi(anioStr)
	if err != nil {
		http.Error(w, "Año de publicación inválido", http.StatusBadRequest)
		return
	}

	// Usar el constructor NuevoLibroCompleto para incluir la sinopsis
	nuevoLibro := models.NuevoLibroCompleto(id, titulo, autor, anio, caratulaURL, sinopsis)

	err = vc.almacen.AgregarLibro(nuevoLibro)
	if err != nil {
		if errors.Is(err, models.ErrLibroYaExiste) {
			http.Error(w, "Un libro con este ID ya existe.", http.StatusConflict)
		} else {
			log.Printf("Error al agregar libro: %v", err)
			http.Error(w, "Error interno del servidor al guardar libro", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

// EditarLibroHTML muestra el formulario para editar un libro existente.
func (vc *MenuController) EditarLibroHTML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	libro, err := vc.almacen.ObtenerLibro(id)
	if err != nil {
		if errors.Is(err, models.ErrLibroNoEncontrado) {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			log.Printf("Error al obtener libro para edición: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	// --- NUEVA LÍNEA DE LOGGING ---
	if libro == nil {
		log.Printf("ERROR: Libro es NIL después de ObtenerLibro para ID: %d. Esto no debería ocurrir si no hubo error.", id)
		http.Error(w, "Error interno del servidor: el libro es nulo inesperadamente.", http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Preparando para editar libro con ID %d y Título '%s'", libro.GetID(), libro.GetTitulo())
	// --- FIN NUEVA LÍNEA DE LOGGING ---

	err = vc.editTpl.Execute(w, libro)
	if err != nil {
		log.Printf("Error al renderizar plantilla editar.html: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// EditarLibroHTMLSubmit maneja el envío del formulario para actualizar un libro.
func (vc *MenuController) EditarLibroHTMLSubmit(w http.ResponseWriter, r *http.Request) {
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
	if val := r.FormValue("titulo"); val != "" {
		updates["titulo"] = val
	}
	if val := r.FormValue("autor"); val != "" {
		updates["autor"] = val
	}
	if val := r.FormValue("anio"); val != "" {
		if anio, err := strconv.Atoi(val); err == nil {
			updates["anio"] = anio
		} else {
			http.Error(w, "Año inválido", http.StatusBadRequest)
			return
		}
	}
	if val := r.FormValue("caratula_url"); val != "" {
		updates["caratula_url"] = val
	}
	// Campo para la sinopsis
	if val := r.FormValue("sinopsis"); val != "" {
		updates["sinopsis"] = val
	}

	if len(updates) == 0 {
		http.Error(w, "No se proporcionaron campos para actualizar", http.StatusBadRequest)
		return
	}

	err = vc.almacen.ActualizarLibro(id, updates)
	if err != nil {
		if errors.Is(err, models.ErrLibroNoEncontrado) {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			log.Printf("Error al actualizar libro: %v", err)
			http.Error(w, "Error interno del servidor al actualizar libro", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

// EliminarLibroHTMLSubmit maneja la eliminación de un libro.
func (vc *MenuController) EliminarLibroHTMLSubmit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	err = vc.almacen.EliminarLibro(id)
	if err != nil {
		if errors.Is(err, models.ErrLibroNoEncontrado) {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			log.Printf("Error al eliminar libro: %v", err)
			http.Error(w, "Error interno del servidor al eliminar libro", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

// VerSinopsisHTML muestra la sinopsis de un libro.
func (vc *MenuController) VerSinopsisHTML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	libro, err := vc.almacen.ObtenerLibro(id)
	if err != nil {
		if errors.Is(err, models.ErrLibroNoEncontrado) {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			log.Printf("Error al obtener libro para sinopsis: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	if libro == nil {
		log.Println("Error: ObtenerLibro devolvió nil sin error explícito para sinopsis")
		http.Error(w, "Error interno del servidor: libro es nulo para sinopsis", http.StatusInternalServerError)
		return
	}

	err = vc.sinopsisTpl.Execute(w, libro) // Pasa el libro a la plantilla de sinopsis
	if err != nil {
		log.Printf("Error al renderizar plantilla sinopsis.html: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}
