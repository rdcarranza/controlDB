package controladores

import (
	"controlDB/cmd/cli/nucleo"
	"controlDB/src/controladores/env"
	"controlDB/src/nucleo/dominio"
	"controlDB/src/nucleo/servicios"
	"controlDB/src/repositorios/mysql"
	"fmt"
	"os"
)

const (
	DirEnv   = ".env"
	CopiaEnv = "./src/controladores/env/env.copia"
)

// ControladorCLI traduce el request de consola, consume el servicio
// y traduce el response hacia stdout/stderr.
type ControladorCLI struct{}

func NewControladorCLI() *ControladorCLI {
	return &ControladorCLI{}
}

// Ejecutar es el punto de entrada del adaptador CLI:
//  1. Traduce el request (flags → dominio)
//  2. Consume el servicio
//  3. Traduce el response (error → stderr con código de salida)
func (c *ControladorCLI) Ejecutar() {
	// 1a. Parsear flags de consola
	params := nucleo.ParsearFlags()

	// 1b. Verificar / crear el .env
	if !env.VerificarEnv(DirEnv, CopiaEnv) {
		c.salir(fmt.Errorf("no se pudo encontrar ni crear el archivo .env"))
	}

	// 1c. Leer credenciales desde el .env y construir la configuración
	cfg, err := c.construirConfig(params)
	if err != nil {
		c.salir(err)
	}

	// 2. Construir repositorio → servicio → ejecutar
	repo, err := mysql.NewRepositorio(cfg)
	if err != nil {
		c.salir(fmt.Errorf("no se pudo inicializar el repositorio: %w", err))
	}

	svc := servicios.NewServicioDB(repo)

	if err := svc.Ejecutar(cfg); err != nil {
		c.salir(err)
	}
}

// construirConfig lee las credenciales del .env y las combina con los flags.
func (c *ControladorCLI) construirConfig(params *nucleo.ParametrosCLI) (*dominio.ConfigDB, error) {
	host, err := env.GetEnv("host_db", DirEnv)
	if err != nil {
		return nil, err
	}

	nameDB, err := env.GetEnv("name_db", DirEnv)
	if err != nil {
		return nil, err
	}

	portDB, err := env.GetEnv("port_db", DirEnv)
	if err != nil {
		return nil, err
	}

	userDB, err := env.GetEnv("user_db", DirEnv)
	if err != nil {
		return nil, err
	}

	pwDB, err := env.GetEnv("pw_db", DirEnv)
	if err != nil {
		return nil, err
	}

	cfg := dominio.CrearConfigDB(
		"mariadb",
		portDB,
		userDB,
		pwDB,
		nameDB,
		host,
		dominio.Operacion(params.Operacion),
		params.ArchivoSQL,
	)

	return cfg, nil
}

// salir traduce cualquier error a stderr y termina el proceso.
func (c *ControladorCLI) salir(err error) {
	fmt.Fprintf(os.Stderr, "✖ %v\n", err)
	os.Exit(1)
}
