package dominio

import "fmt"

// Operacion define las acciones soportadas sobre la base de datos.
type Operacion string

const (
	OperacionInicializar Operacion = "init"
	OperacionRespaldar   Operacion = "backup"
	OperacionRestaurar   Operacion = "restore"
)

// ConfigDB contiene los parámetros de conexión y la operación a ejecutar.
// Las credenciales vienen del .env; Op y ArchivoSQL vienen de los flags -o y -a.
type ConfigDB struct {
	Tipo       string
	Puerto     string
	Usuario    string
	Contraseña string
	BD         string
	Host       string
	Op         Operacion
	ArchivoSQL string
}

func CrearConfigDB(tipo, puerto, usuario, contraseña, bd, host string, op Operacion, archivoSQL string) *ConfigDB {
	return &ConfigDB{
		Tipo:       tipo,
		Puerto:     puerto,
		Usuario:    usuario,
		Contraseña: contraseña,
		BD:         bd,
		Host:       host,
		Op:         op,
		ArchivoSQL: archivoSQL,
	}
}

// DSN construye el Data Source Name para el driver MySQL/MariaDB.
func (c *ConfigDB) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&charset=utf8mb4&parseTime=true",
		c.Usuario, c.Contraseña, c.Host, c.Puerto, c.BD)
}

// DSNSinBD construye el DSN sin base de datos (para poder crearla si no existe).
func (c *ConfigDB) DSNSinBD() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?multiStatements=true&charset=utf8mb4",
		c.Usuario, c.Contraseña, c.Host, c.Puerto)
}

// Validar verifica que los campos de conexión del .env estén completos.
// Op y ArchivoSQL son validados antes por main (vienen de los flags).
func (c *ConfigDB) Validar() error {
	if c.Host == "" {
		return fmt.Errorf("host_db no puede estar vacío en el .env")
	}
	if c.Usuario == "" {
		return fmt.Errorf("user_db no puede estar vacío en el .env")
	}
	if c.BD == "" {
		return fmt.Errorf("name_db no puede estar vacío en el .env")
	}
	if c.Puerto == "" {
		c.Puerto = "3306"
	}
	if c.Op == "" {
		return fmt.Errorf("operación no especificada (flag -o)")
	}
	if c.ArchivoSQL == "" {
		return fmt.Errorf("archivo SQL no especificado (flag -a)")
	}
	return nil
}
