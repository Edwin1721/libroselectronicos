package db

import (
	"database/sql"
	"log" // Asegúrate de que esta importación esté aquí

	"libroselectronicos/models"

	_ "github.com/mattn/go-sqlite3" // Driver SQLite
)

// Extendemos la interfaz para incluir operaciones de Usuario y Alquiler
type LibroAlmacenamiento interface {
	AgregarLibro(libro *models.Libro) error
	ObtenerLibro(id int) (*models.Libro, error)
	ListarLibros() []*models.Libro
	ActualizarLibro(id int, updates map[string]interface{}) error
	EliminarLibro(id int) error

	// --- Nuevas operaciones para Usuarios ---
	AgregarUsuario(usuario *models.Usuario) error
	ObtenerUsuarioPorID(id int) (*models.Usuario, error)
	ObtenerUsuarioPorUsername(username string) (*models.Usuario, error)

	// --- Nuevas operaciones para Alquileres (en un paso posterior) ---
	// Esto lo añadiremos después de definir el modelo Alquiler

	Close() error // Método para cerrar la conexión a la base de datos
}

type sqliteAlmacenamiento struct {
	db *sql.DB
}

// NuevoAlmacen crea una nueva instancia de sqliteAlmacenamiento.
func NuevoAlmacen() LibroAlmacenamiento {
	db, err := sql.Open("sqlite3", "./libros.db")
	if err != nil {
		log.Printf("Error al abrir la base de datos: %v", err)
		return nil
	}

	// Crear la tabla 'libros' si no existe
	createLibrosTableSQL := `
	CREATE TABLE IF NOT EXISTS libros (
		id INTEGER PRIMARY KEY,
		titulo TEXT NOT NULL,
		autor TEXT NOT NULL,
		anio INTEGER NOT NULL,
		caratula_url TEXT,
		sinopsis TEXT -- Campo para la sinopsis
	);`
	_, err = db.Exec(createLibrosTableSQL)
	if err != nil {
		log.Printf("Error al crear la tabla 'libros': %v", err)
		return nil
	}

	// Crear la tabla 'usuarios' si no existe
	createUsuariosTableSQL := `
    CREATE TABLE IF NOT EXISTS usuarios (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        email TEXT,
        rol TEXT NOT NULL DEFAULT 'lector'
    );`
	_, err = db.Exec(createUsuariosTableSQL)
	if err != nil {
		log.Printf("Error al crear la tabla 'usuarios': %v", err)
		return nil
	}

	return &sqliteAlmacenamiento{db: db}
}

// NuevoAlmacenForTest existe para que los tests puedan usar una DB separada
func NewAlmacenForTest(dbPath string) LibroAlmacenamiento {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error al abrir la base de datos de prueba: %v", err)
		return nil
	}

	// Crear la tabla 'libros' si no existe (para tests)
	createLibrosTableSQL := `
    CREATE TABLE IF NOT EXISTS libros (
        id INTEGER PRIMARY KEY,
        titulo TEXT NOT NULL,
        autor TEXT NOT NULL,
        anio INTEGER NOT NULL,
        caratula_url TEXT,
        sinopsis TEXT
    );`
	_, err = db.Exec(createLibrosTableSQL)
	if err != nil {
		log.Printf("Error al crear la tabla 'libros' para tests: %v", err)
		return nil
	}

	// Crear la tabla 'usuarios' si no existe (para tests)
	createUsuariosTableSQL := `
    CREATE TABLE IF NOT EXISTS usuarios (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        email TEXT,
        rol TEXT NOT NULL DEFAULT 'lector'
    );`
	_, err = db.Exec(createUsuariosTableSQL)
	if err != nil {
		log.Printf("Error al crear la tabla 'usuarios' para tests: %v", err)
		return nil
	}

	return &sqliteAlmacenamiento{db: db}
}

func (s *sqliteAlmacenamiento) Close() error {
	return s.db.Close()
}

// --- Operaciones de Libros (sin cambios, solo se incluyen para la referencia completa) ---
func (s *sqliteAlmacenamiento) AgregarLibro(libro *models.Libro) error {
	existingBook, err := s.ObtenerLibro(libro.ID)
	if err == nil && existingBook != nil {
		return models.ErrLibroYaExiste
	}
	if err != nil && err != models.ErrLibroNoEncontrado {
		return err
	}

	stmt, err := s.db.Prepare("INSERT INTO libros(id, titulo, autor, anio, caratula_url, sinopsis) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(libro.ID, libro.Titulo, libro.Autor, libro.Anio, libro.CaratulaURL, libro.Sinopsis)
	return err
}

func (s *sqliteAlmacenamiento) ObtenerLibro(id int) (*models.Libro, error) {
	row := s.db.QueryRow("SELECT id, titulo, autor, anio, caratula_url, sinopsis FROM libros WHERE id = ?", id)
	libro := &models.Libro{}
	err := row.Scan(&libro.ID, &libro.Titulo, &libro.Autor, &libro.Anio, &libro.CaratulaURL, &libro.Sinopsis)
	if err == sql.ErrNoRows {
		return nil, models.ErrLibroNoEncontrado
	}
	return libro, err
}

func (s *sqliteAlmacenamiento) ListarLibros() []*models.Libro {
	rows, err := s.db.Query("SELECT id, titulo, autor, anio, caratula_url, sinopsis FROM libros")
	if err != nil {
		log.Printf("Error al listar libros: %v", err)
		return nil
	}
	defer rows.Close()

	libros := []*models.Libro{}
	for rows.Next() {
		libro := &models.Libro{}
		if err := rows.Scan(&libro.ID, &libro.Titulo, &libro.Autor, &libro.Anio, &libro.CaratulaURL, &libro.Sinopsis); err != nil {
			log.Printf("Error al escanear libro: %v", err)
			continue
		}
		libros = append(libros, libro)
	}
	return libros
}

func (s *sqliteAlmacenamiento) ActualizarLibro(id int, updates map[string]interface{}) error {
	query := "UPDATE libros SET "
	args := []interface{}{}
	i := 0
	for key, value := range updates {
		query += key + " = ?"
		args = append(args, value)
		if i < len(updates)-1 {
			query += ", "
		}
		i++
	}
	query += " WHERE id = ?"
	args = append(args, id)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrLibroNoEncontrado // Si no se actualizó ninguna fila, es porque no existe
	}
	return nil
}

func (s *sqliteAlmacenamiento) EliminarLibro(id int) error {
	stmt, err := s.db.Prepare("DELETE FROM libros WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrLibroNoEncontrado // Si no se eliminó ninguna fila, es porque no existe
	}
	return nil
}

// --- NUEVAS Operaciones de Usuario ---

// AgregarUsuario añade un nuevo usuario a la base de datos.
func (s *sqliteAlmacenamiento) AgregarUsuario(usuario *models.Usuario) error {
	// Verificar si el usuario ya existe por nombre de usuario
	existingUser, err := s.ObtenerUsuarioPorUsername(usuario.Username)
	if err == nil && existingUser != nil {
		log.Printf("DEBUG: Usuario %s ya existe en la DB.", usuario.Username) // NUEVO LOG
		return models.ErrUsuarioYaExiste                                      // El usuario ya existe
	}
	if err != nil && err != models.ErrUsuarioNoEncontrado {
		log.Printf("ERROR: Fallo al verificar existencia de usuario %s: %v", usuario.Username, err) // NUEVO LOG
		return err                                                                                  // Otro tipo de error
	}

	log.Printf("DEBUG: Preparando sentencia SQL para insertar usuario %s.", usuario.Username) // NUEVO LOG
	stmt, err := s.db.Prepare("INSERT INTO usuarios(username, password, email, rol) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Printf("ERROR: Fallo al preparar sentencia INSERT para usuario %s: %v", usuario.Username, err) // NUEVO LOG
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(usuario.Username, usuario.Password, usuario.Email, usuario.Rol)
	if err != nil {
		log.Printf("ERROR: Fallo al ejecutar INSERT para usuario %s: %v", usuario.Username, err) // NUEVO LOG
	} else {
		log.Printf("DEBUG: Sentencia INSERT ejecutada para usuario %s.", usuario.Username) // NUEVO LOG
	}
	return err
}

// ObtenerUsuarioPorID recupera un usuario por su ID.
func (s *sqliteAlmacenamiento) ObtenerUsuarioPorID(id int) (*models.Usuario, error) {
	row := s.db.QueryRow("SELECT id, username, password, email, rol FROM usuarios WHERE id = ?", id)
	usuario := &models.Usuario{}
	err := row.Scan(&usuario.ID, &usuario.Username, &usuario.Password, &usuario.Email, &usuario.Rol)
	if err == sql.ErrNoRows {
		return nil, models.ErrUsuarioNoEncontrado
	}
	return usuario, err
}

// ObtenerUsuarioPorUsername recupera un usuario por su nombre de usuario.
func (s *sqliteAlmacenamiento) ObtenerUsuarioPorUsername(username string) (*models.Usuario, error) {
	row := s.db.QueryRow("SELECT id, username, password, email, rol FROM usuarios WHERE username = ?", username)
	usuario := &models.Usuario{}
	err := row.Scan(&usuario.ID, &usuario.Username, &usuario.Password, &usuario.Email, &usuario.Rol)
	if err == sql.ErrNoRows {
		return nil, models.ErrUsuarioNoEncontrado
	}
	return usuario, err
}
