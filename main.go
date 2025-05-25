package main

import (
	"fmt"
	"libroselectronicos/controllers"
	"libroselectronicos/views"
)

func main() {
	fmt.Println("Bienvenido al sistema de Libros Electrónicos")
	for {
		views.MostrarMenu()
		opcion := views.LeerOpcion()
		switch opcion {
		case 1:
			controllers.ListarLibros()
		case 0:
			fmt.Println("Saliendo...")
			return
		default:
			fmt.Println("Opción inválida, intenta de nuevo.")
		}
	}
}
