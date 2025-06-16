package views

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func MostrarMenu() {
	fmt.Println("\n🔸 Menú Principal")
	fmt.Println("1. Agregar Libro")
	fmt.Println("2. Listar Libros")
	fmt.Println("0. Salir")
	fmt.Print("Selecciona una opción: ")
}

func LeerOpcion() int {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	opcion, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	return opcion
}
