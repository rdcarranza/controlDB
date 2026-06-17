package servicios

import (
	"controlDB/src/nucleo/dominio"
	"controlDB/src/nucleo/puertos"
	"fmt"
)

// ServicioDB orquesta los casos de uso de gestión de la base de datos.
type ServicioDB struct {
	repo puertos.RepositorioDB
}

func NewServicioDB(repo puertos.RepositorioDB) *ServicioDB {
	return &ServicioDB{repo: repo}
}

// Ejecutar valida la configuración y despacha la operación correspondiente.
func (s *ServicioDB) Ejecutar(cfg *dominio.ConfigDB) error {
	if err := cfg.Validar(); err != nil {
		return fmt.Errorf("configuración inválida: %w", err)
	}

	switch cfg.Op {
	case dominio.OperacionInicializar:
		fmt.Printf("▶ Inicializando BD '%s' con '%s'...\n", cfg.BD, cfg.ArchivoSQL)
		if err := s.repo.Inicializar(cfg); err != nil {
			return fmt.Errorf("error al inicializar: %w", err)
		}
		fmt.Println("✔ BD inicializada correctamente.")

	case dominio.OperacionRespaldar:
		fmt.Printf("▶ Respaldando BD '%s' → '%s'...\n", cfg.BD, cfg.ArchivoSQL)
		if err := s.repo.Respaldar(cfg); err != nil {
			return fmt.Errorf("error al respaldar: %w", err)
		}
		fmt.Println("✔ Respaldo completado correctamente.")

	case dominio.OperacionRespaldarNativo:
		fmt.Printf("▶ Respaldando BD '%s' con mariadb-backup → '%s'...\n", cfg.BD, cfg.ArchivoSQL)
		if err := s.repo.RespaldarNativo(cfg); err != nil {
			return fmt.Errorf("error al respaldar con mariadb-backup: %w", err)
		}
		fmt.Println("✔ Respaldo nativo completado correctamente.")

	case dominio.OperacionRestaurar:
		fmt.Printf("▶ Restaurando BD '%s' desde '%s'...\n", cfg.BD, cfg.ArchivoSQL)
		if err := s.repo.Restaurar(cfg); err != nil {
			return fmt.Errorf("error al restaurar: %w", err)
		}
		fmt.Println("✔ Restauración completada correctamente.")

	default:
		return fmt.Errorf("operación desconocida: '%s' (valores válidos: init | backup | backup2  | restore)", cfg.Op)
	}

	return nil
}
