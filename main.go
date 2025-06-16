package main

import (
	"fmt"
	"libroselectronicos/controllers"
	"libroselectronicos/views"
)

func main() {
	controlador := controllers.NewLibroController()
	fmt.Println("📚 Bienvenido al Sistema de Libros Electrónicos")

	for {
		views.MostrarMenu()
		opcion := views.LeerOpcion()

		switch opcion {
		case 1:
			var id, anio int
			var titulo, autor string
			fmt.Print("ID: ")
			fmt.Scan(&id)
			fmt.Print("Título: ")
			fmt.Scan(&titulo)
			fmt.Print("Autor: ")
			fmt.Scan(&autor)
			fmt.Print("Año: ")
			fmt.Scan(&anio)
			controlador.AgregarLibro(id, titulo, autor, anio)

		case 2:
			controlador.ListarLibros()

		case 0:
			fmt.Println("👋 Saliendo...")
			return

		default:
			fmt.Println("❌ Opción inválida")
		}
	}
}
