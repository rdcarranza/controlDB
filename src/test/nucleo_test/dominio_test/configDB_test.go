package dominio_test

import (
	"controlDB/src/nucleo/dominio"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrearConfigDB(t *testing.T) {
	cfg := dominio.CrearConfigDB("mariadb", "3306", "root", "secret", "testdb", "localhost",
		dominio.OperacionInicializar, "schema.sql")

	assert.Equal(t, "mariadb", cfg.Tipo)
	assert.Equal(t, "3306", cfg.Puerto)
	assert.Equal(t, "root", cfg.Usuario)
	assert.Equal(t, "testdb", cfg.BD)
	assert.Equal(t, dominio.OperacionInicializar, cfg.Op)
	assert.Equal(t, "schema.sql", cfg.ArchivoSQL)
}

func TestConfigDB_DSN(t *testing.T) {
	cfg := dominio.CrearConfigDB("mariadb", "3306", "root", "secret", "testdb", "localhost",
		dominio.OperacionInicializar, "schema.sql")

	dsn := cfg.DSN()
	assert.Contains(t, dsn, "root:secret@tcp(localhost:3306)/testdb")
	assert.Contains(t, dsn, "multiStatements=true")
}

func TestConfigDB_DSNSinBD(t *testing.T) {
	cfg := dominio.CrearConfigDB("mariadb", "3306", "root", "secret", "testdb", "localhost",
		dominio.OperacionInicializar, "schema.sql")

	dsn := cfg.DSNSinBD()
	assert.Contains(t, dsn, "root:secret@tcp(localhost:3306)/")
	assert.NotContains(t, dsn, "testdb")
}

func TestConfigDB_Validar_OK(t *testing.T) {
	cfg := dominio.CrearConfigDB("mariadb", "3306", "root", "secret", "testdb", "localhost",
		dominio.OperacionRespaldar, "respaldo.sql")

	assert.NoError(t, cfg.Validar())
}

func TestConfigDB_Validar_SinHost(t *testing.T) {
	cfg := dominio.CrearConfigDB("mariadb", "3306", "root", "secret", "testdb", "",
		dominio.OperacionRespaldar, "respaldo.sql")

	assert.Error(t, cfg.Validar())
}

func TestConfigDB_Validar_SinUsuario(t *testing.T) {
	cfg := dominio.CrearConfigDB("mariadb", "3306", "", "secret", "testdb", "localhost",
		dominio.OperacionRespaldar, "respaldo.sql")

	assert.Error(t, cfg.Validar())
}

func TestConfigDB_Validar_SinOperacion(t *testing.T) {
	cfg := dominio.CrearConfigDB("mariadb", "3306", "root", "secret", "testdb", "localhost",
		"", "respaldo.sql")

	assert.Error(t, cfg.Validar())
}

func TestConfigDB_Validar_PuertoVacioDefecto3306(t *testing.T) {
	cfg := dominio.CrearConfigDB("mariadb", "", "root", "secret", "testdb", "localhost",
		dominio.OperacionInicializar, "schema.sql")

	err := cfg.Validar()
	assert.NoError(t, err)
	assert.Equal(t, "3306", cfg.Puerto)
}

func TestOperacionConstantes(t *testing.T) {
	assert.Equal(t, dominio.Operacion("init"), dominio.OperacionInicializar)
	assert.Equal(t, dominio.Operacion("backup"), dominio.OperacionRespaldar)
	assert.Equal(t, dominio.Operacion("restore"), dominio.OperacionRestaurar)
}
