package main

import (
	"log"
	"net/http"

	"libroselectronicos/db"
	"libroselectronicos/views"

	"github.com/gorilla/mux"
)

func main() {
	almacen := db.NuevoAlmacen()
	if almacen == nil {
		log.Fatalf("No se pudo inicializar la base de datos.")
	}
	defer almacen.Close()

	viewsController := views.NewMenuController(almacen)

	router := mux.NewRouter()

	// Rutas de Autenticación
	router.HandleFunc("/registro", viewsController.RegistrarUsuarioHTML).Methods("GET")
	router.HandleFunc("/registro", viewsController.RegistrarUsuarioSubmit).Methods("POST")
	router.HandleFunc("/login", viewsController.LoginHTML).Methods("GET")
	router.HandleFunc("/login", viewsController.LoginSubmit).Methods("POST")
	router.HandleFunc("/logout", viewsController.Logout).Methods("POST") // Normalmente un POST para logout

	// Rutas para las vistas HTML de libros
	router.HandleFunc("/", viewsController.Index).Methods("GET")
	router.HandleFunc("/libros", viewsController.ListarLibrosHTML).Methods("GET")
	router.HandleFunc("/libros/crear", viewsController.CrearLibroHTML).Methods("GET")
	router.HandleFunc("/libros/crear", viewsController.CrearLibroHTMLSubmit).Methods("POST")
	router.HandleFunc("/libros/{id}/editar", viewsController.EditarLibroHTML).Methods("GET")
	router.HandleFunc("/libros/{id}/editar", viewsController.EditarLibroHTMLSubmit).Methods("POST")
	router.HandleFunc("/libros/{id}/eliminar", viewsController.EliminarLibroHTMLSubmit).Methods("POST")
	router.HandleFunc("/libros/{id}/sinopsis", viewsController.VerSinopsisHTML).Methods("GET")

	// Servir archivos estáticos (CSS, JS, imágenes)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	port := ":8080"
	log.Printf("Servidor iniciado en http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
