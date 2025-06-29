package models

import "errors"

// ErrUsuarioNoEncontrado es un error que se devuelve cuando un usuario no se encuentra.
var ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")

// ErrUsuarioYaExiste es un error que se devuelve cuando un usuario con el mismo nombre de usuario ya existe.
var ErrUsuarioYaExiste = errors.New("nombre de usuario ya existe")

// Usuario representa la estructura de un usuario en la aplicación.
type Usuario struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // Aquí se almacenará la contraseña HASHED
	Email    string `json:"email"`
	Rol      string `json:"rol"` // Ej. "lector", "administrador"
}

// NuevoUsuario crea una nueva instancia de Usuario.
// La contraseña debe ser un hash, no texto plano.
func NuevoUsuario(id int, username, password, email, rol string) *Usuario {
	return &Usuario{
		ID:       id,
		Username: username,
		Password: password,
		Email:    email,
		Rol:      rol,
	}
}

// Getters para campos de Usuario (opcional, pero buena práctica)
func (u *Usuario) GetID() int {
	return u.ID
}

func (u *Usuario) GetUsername() string {
	return u.Username
}

func (u *Usuario) GetPassword() string {
	return u.Password
}

func (u *Usuario) GetEmail() string {
	return u.Email
}

func (u *Usuario) GetRol() string {
	return u.Rol
}
