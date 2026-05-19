package controlador_error

import "fmt"

// Error controlado
type ErrorC struct {
	Archivo string
	Modulo  string
	Mensaje string
}

// implementación de interfaz
func controlador_Error() error {
	return &ErrorC{}
}

// constructor
func NuevoControladorError() *ErrorC {
	return controlador_Error().(*ErrorC)
}

func (c *ErrorC) Error() string {
	c.Mensaje = fmt.Sprintf("Error en el archivo: %s,  módulo -> %s", c.Archivo, c.Modulo)
	return c.Mensaje
}
