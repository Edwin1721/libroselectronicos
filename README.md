Objetivo del Programa

El objetivo principal de este sistema es proporcionar una plataforma web intuitiva y funcional para la gestión de un catálogo de libros electrónicos. Permite a los usuarios navegar por los libros, registrarse, iniciar sesión, y para los usuarios con el rol de "administrador", realizar operaciones completas de gestión (Crear, Leer, Actualizar, Eliminar) sobre los libros. Además, incorpora un módulo de alquiler que permite a los usuarios "prestar" libros del catálogo y gestionarlos.

Funcionalidades Principales

El sistema está organizado en módulos clave que abarcan las siguientes funcionalidades:

### 1. Gestión de Libros (CRUD)
* **Listado de Libros:** Muestra una tabla con todos los libros disponibles en el catálogo, incluyendo Título, Autor, Año, Estado (Disponible/Alquilado) y opciones de acción.
* **Añadir Nuevo Libro:** Permite a los usuarios con rol de `administrador` agregar nuevos libros al catálogo, especificando ID, Título, Autor, Año, URL de la carátula y Sinopsis.
* **Editar Libro:** Facilita la modificación de la información de un libro existente (solo para `administradores`).
* **Eliminar Libro:** Permite la eliminación de libros del catálogo (solo para `administradores`).
* **Ver Sinopsis:** Muestra los detalles completos y la sinopsis de un libro específico.
* **Estado de Disponibilidad:** Cada libro tiene un estado `Disponible` o `Alquilado`, que se actualiza automáticamente con las operaciones de alquiler.

### 2. Gestión de Usuarios y Autenticación
* **Registro de Usuarios:** Permite a nuevos usuarios crear una cuenta con un nombre de usuario, contraseña y correo electrónico. Se asigna un rol por defecto (`lector`).
* **Inicio de Sesión (Login):** Autentica a los usuarios mediante sus credenciales. Las contraseñas se almacenan de forma segura (hasheadas con bcrypt).
* **Cierre de Sesión (Logout):** Permite a los usuarios finalizar su sesión activa.
* **Roles de Usuario:**
    * `lector`: Puede ver el listado de libros, ver sinopsis, alquilar y devolver libros, y ver sus alquileres.
    * `administrador`: Posee todas las funcionalidades del `lector`, además de poder añadir, editar y eliminar libros.

### 3. Gestión de Alquileres
* **Alquilar Libro:** Los usuarios logueados pueden alquilar un libro que esté `Disponible`. Al alquilar, el estado del libro cambia a `Alquilado`.
* **Mis Alquileres:** Un usuario puede ver una lista de todos los libros que ha alquilado, incluyendo la fecha de alquiler y la fecha de devolución (si ya fue devuelto).
* **Devolver Libro:** Permite a un usuario marcar un libro como devuelto. Al devolver, el estado del libro vuelve a `Disponible`.

---

## ⚙️ Tecnologías Utilizadas

* **Lenguaje de Programación:** Go (Golang)
* **Base de Datos:** SQLite (archivo `libros.db` para persistencia de datos local)
* **Framework Web:** `net/http` (librería estándar de Go)
* **Enrutador HTTP:** `github.com/gorilla/mux`
* **Gestión de Sesiones:** `github.com/gorilla/sessions`
* **Hashing de Contraseñas:** `golang.org/x/crypto/bcrypt`
* **Controlador de Base de Datos SQLite:** `github.com/mattn/go-sqlite3`
* **Plantillas HTML:** `html/template` (librería estándar de Go)
* **Estilos:** CSS (con archivo `static/css/style.css`)

---

## 🚀 Cómo Ejecutar el Proyecto

Sigue estos pasos para poner en marcha el servidor localmente:

1.  **Clonar el Repositorio:**
    ```bash
    git clone [URL_DE_TU_REPOSITORIO]
    cd [nombre-de-la-carpeta-del-proyecto]
    ```

2.  **Instalar Dependencias de Go:**
    Asegúrate de tener Go instalado. Luego, ejecuta en la raíz del proyecto para descargar las dependencias:
    ```bash
    go mod tidy
    ```

3.  **Eliminar la Base de Datos Existente (Opcional, pero recomendado para el primer uso):**
    Para asegurar que la base de datos se crea con el esquema más reciente (incluyendo las tablas de usuarios y alquileres, y el campo `disponible` en libros), puedes eliminar el archivo `libros.db` si ya existe:
    ```bash
    del libros.db # En Windows
    # rm libros.db # En Linux/macOS
    ```

4.  **Ejecutar la Aplicación:**
    ```bash
    go run main.go
    ```
    El servidor se iniciará y estará accesible en `http://localhost:8080`.

5.  **Primer Uso - Registro y Administración (Recomendado):**
    * Abre tu navegador y navega a `http://localhost:8080/registro`.
    * Registra un nuevo usuario (ej., `admin`, `contraseña123`).
    * **Importante:** Detén la aplicación (`Ctrl+C` en la terminal).
    * Abre `DB Browser for SQLite` y abre el archivo `libros.db` en la raíz de tu proyecto.
    * Ve a la pestaña `Browse Data`, selecciona la tabla `usuarios` y cambia el `rol` de tu usuario recién creado (`admin`) a `administrador`. Haz clic en "Write Changes".
    * Cierra `DB Browser for SQLite`.
    * Vuelve a ejecutar la aplicación: `go run main.go`.
    * Inicia sesión con tu usuario `admin` en `http://localhost:8080/login`. Ahora tendrás acceso a las funcionalidades de administrador (Añadir, Editar, Eliminar libros).

---

## 📂 Estructura del Proyecto

.
├── main.go               # Punto de entrada y configuración de rutas
├── go.mod                # Módulos de Go y dependencias
├── go.sum                # Sumas de verificación de dependencias
├── db/                   # Lógica de interacción con la base de datos
│   └── storage.go        # Implementación del almacenamiento (SQLite)
├── models/               # Definiciones de estructuras de datos (modelos)
│   ├── libro.go          # Estructura y métodos para Libro
│   ├── usuario.go        # Estructura y métodos para Usuario
│   └── alquiler.go       # Estructura y métodos para Alquiler
├── views/                # Controladores HTTP y lógica de negocio
│   └── menu.go           # Manejadores de rutas y renderizado de plantillas
├── templates/            # Archivos HTML (vistas)
│   ├── index.html
│   ├── listar.html
│   ├── crear.html
│   ├── editar.html
│   ├── sinopsis.html
│   ├── registro.html
│   ├── login.html
│   └── mis_alquileres.html # Nueva plantilla para alquileres
└── static/               # Archivos estáticos (CSS, JS, imágenes)
└── css/
└── style.css     # Estilos CSS de la aplicación

