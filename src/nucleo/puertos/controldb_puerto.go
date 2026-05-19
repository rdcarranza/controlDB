package puertos

import "controlDB/src/nucleo/dominio"

// RepositorioDB define las operaciones de bajo nivel contra la base de datos.
// Es el puerto de salida que deben implementar los adaptadores (MySQL, etc.).
type RepositorioDB interface {
	Inicializar(cfg *dominio.ConfigDB) error
	Respaldar(cfg *dominio.ConfigDB) error
	Restaurar(cfg *dominio.ConfigDB) error
}
