package dominio_test

import (
	"controlDB/src/controladores/env"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tmpEnv(t *testing.T, contenido string) string {
	t.Helper()
	dir := t.TempDir()
	ruta := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(ruta, []byte(contenido), 0644))
	return ruta
}

func TestGetEnv_ValorSimple(t *testing.T) {
	ruta := tmpEnv(t, "host_db=localhost\nport_db=3306\n")

	val, err := env.GetEnv("host_db", ruta)
	assert.NoError(t, err)
	assert.Equal(t, "localhost", val)
}

func TestGetEnv_VariableNoEncontrada(t *testing.T) {
	ruta := tmpEnv(t, "host_db=localhost\n")

	_, err := env.GetEnv("user_db", ruta)
	assert.Error(t, err)
}

func TestGetEnv_IgnoraComentarios(t *testing.T) {
	ruta := tmpEnv(t, "# esto es un comentario\nhost_db=192.168.1.1\n")

	val, err := env.GetEnv("host_db", ruta)
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", val)
}

func TestGetEnv_ValorConEspacios(t *testing.T) {
	ruta := tmpEnv(t, "name_db = mi_base \n")

	val, err := env.GetEnv("name_db", ruta)
	assert.NoError(t, err)
	assert.Equal(t, "mi_base", val)
}

func TestEnvExiste_Existe(t *testing.T) {
	ruta := tmpEnv(t, "host_db=localhost\n")
	assert.True(t, env.EnvExiste(ruta))
}

func TestEnvExiste_NoExiste(t *testing.T) {
	assert.False(t, env.EnvExiste("/tmp/este_archivo_no_existe_jamas.env"))
}

func TestCrearEnv(t *testing.T) {
	dir := t.TempDir()
	copia := filepath.Join(dir, "env.copia")
	destino := filepath.Join(dir, ".env")

	require.NoError(t, os.WriteFile(copia, []byte("host_db=\nport_db=3306\n"), 0644))

	err := env.CrearEnv(destino, copia)
	assert.NoError(t, err)
	assert.True(t, env.EnvExiste(destino))
}
