package nucleo

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// ParametrosCLI contiene los argumentos parseados desde la consola.
type ParametrosCLI struct {
	Operacion  string
	ArchivoSQL string
}

// ParsearFlags lee los flags -o y -a, valida que estén presentes
// y devuelve los parámetros listos para usar.
func ParsearFlags() *ParametrosCLI {
	params := &ParametrosCLI{}

	flag.StringVar(&params.Operacion, "o", "", "Operación a ejecutar: init | backup | backup2 | restore ")
	flag.StringVar(&params.ArchivoSQL, "a", "", "Ruta al archivo .sql")
	flag.Usage = MostrarUso
	flag.Parse()

	if params.Operacion == "" || params.ArchivoSQL == "" {
		MostrarUso()
		os.Exit(1)
	}

	params.ArchivoSQL = expandirRuta(params.ArchivoSQL)
	return params
}

// MostrarUso imprime la ayuda del comando en stderr.
func MostrarUso() {
	fmt.Fprintf(os.Stderr, `
Uso: controlDB -o <operacion> -a <archivo.sql>

Flags:
  -o    Operación a ejecutar:
          init     Crea la BD si no existe y ejecuta el SQL.
          backup   Exporta la BD al archivo indicado.
		  backup2  Respaldo binario físico con mariadb-backup (requiere instalación).
          restore  Importa el archivo SQL sobre la BD existente.
  -a    Ruta al archivo .sql (soporta ~). [backup2 recibe un directorio vacío, no un archivo .sql]  

Ejemplos:
  controlDB -o init    -a ~/scripts/schema.sql
  controlDB -o backup  -a ~/respaldos/backup.sql
  controlDB -o backup2 -a ~/respaldos/
  controlDB -o restore -a ~/respaldos/backup.sql

Las credenciales de conexión se leen del archivo .env:
  host_db, name_db, port_db, user_db, pw_db
`)
}

// expandirRuta reemplaza el ~ inicial por el home del usuario.
func expandirRuta(ruta string) string {
	if len(ruta) > 0 && ruta[0] == '~' {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, ruta[1:])
		}
	}
	return ruta
}
