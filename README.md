🧠 Sistema de Gestión de Libros Electrónicos
🗂️ Estructura del Proyecto
main.go: Punto de entrada de la aplicación.

controllers/: Lógica para manejar operaciones como listar libros.

views/: Interfaz por consola (menú de usuario).

models/: Estructuras de datos, como el modelo Libro, con encapsulación y métodos asociados.

db/: Preparado para futura conexión a una base de datos.

🛠️ Tecnologías y Herramientas
Lenguaje: Go (Golang)

Paquetes estándar utilizados:

fmt

bufio

os

errors

🧩 Paradigmas Aplicados
Programación funcional (en parte):

Uso de funciones puras para separar responsabilidades.

Flujo claro sin variables globales innecesarias.

Modularización de la lógica.

Programación orientada a objetos:

Encapsulación de atributos del modelo Libro usando métodos Get y Set.

Organización por paquetes coherentes.

Uso de interfaces para separar el comportamiento de impresión de libros (Imprimible).

Comentarios descriptivos en las funciones más complejas para facilitar el mantenimiento y comprensión del código.

⚙️ Funcionalidades Implementadas
Mostrar menú principal por consola.

Listar libros electrónicos disponibles.

Salir del sistema de forma segura.

🔐 Características de Buen Diseño
Encapsulación: Los atributos del modelo Libro están protegidos mediante funciones Get y Set.

Manejo de errores: Validaciones y reportes de error cuando se intenta acceder o modificar valores inválidos.

Interfaces: Abstracción para representar comportamientos comunes sin acoplarse a una implementación específica.

Comentarios: Agregados en los bloques de código complejos y en funciones clave, como el menú, iteraciones y validaciones.

🚀 Futuras Mejoras
Agregar, eliminar y buscar libros.

Conectar con base de datos real (MySQL, PostgreSQL o SQLite).

Agregar autenticación de usuario.

Interfaz gráfica con Go o web con HTML y JavaScript.

👨‍💻 Autor
Edwin Bermeo
Proyecto universitario – Programación Orientada a Objetos (GoLand)




