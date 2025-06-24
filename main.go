package main

import (
	"log"
	"net/http"

	"libroselectronicos/controllers"
	"libroselectronicos/db"
	"libroselectronicos/views"

	"github.com/gorilla/mux"
)

func main() {
	almacen := db.NuevoAlmacen()
	defer almacen.Close()

	// CORRECCIÓN: Asegúrate que NewApiLibroController tiene 'New' y 'Api' capitalizados correctamente
	// y que ApiLibroController también tiene 'Api' capitalizado.
	apiController := controllers.NewApiLibroController(almacen)
	viewsController := views.NewViewsController(almacen)

	router := mux.NewRouter()

	// --- Rutas API (JSON) ---
	// CORRECCIÓN: Los nombres de los métodos deben ser los definidos en api_libro_controller.go
	router.HandleFunc("/api/libros", apiController.GetLibrosAPI).Methods("GET")
	router.HandleFunc("/api/libros/{id}", apiController.GetLibroByIDAPI).Methods("GET")
	router.HandleFunc("/api/libros", apiController.CreateLibroAPI).Methods("POST")
	router.HandleFunc("/api/libros/{id}", apiController.UpdateLibroAPI).Methods("PUT")
	router.HandleFunc("/api/libros/{id}", apiController.DeleteLibroAPI).Methods("DELETE")

	// --- Rutas de Vista (HTML) ---
	router.HandleFunc("/", viewsController.IndexHandler).Methods("GET")
	router.HandleFunc("/libros", viewsController.ListarLibrosHTML).Methods("GET")
	router.HandleFunc("/libros/crear", viewsController.CrearLibroHTMLForm).Methods("GET")
	router.HandleFunc("/libros/crear", viewsController.CrearLibroHTMLSubmit).Methods("POST")
	router.HandleFunc("/libros/{id}/editar", viewsController.EditarLibroHTMLForm).Methods("GET")
	router.HandleFunc("/libros/{id}/editar", viewsController.EditarLibroHTMLSubmit).Methods("POST")
	router.HandleFunc("/libros/{id}/eliminar", viewsController.EliminarLibroHTML).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
