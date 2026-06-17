package dominio_test

import (
	"controlDB/src/nucleo/dominio"
	"controlDB/src/nucleo/servicios"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─── Mock del repositorio ─────────────────────────────────────────────────────

type repoMock struct {
	llamadas  []string
	errorFijo error
}

func (m *repoMock) Inicializar(cfg *dominio.ConfigDB) error {
	m.llamadas = append(m.llamadas, "inicializar")
	return m.errorFijo
}

func (m *repoMock) Respaldar(cfg *dominio.ConfigDB) error {
	m.llamadas = append(m.llamadas, "respaldar")
	return m.errorFijo
}

func (m *repoMock) RespaldarNativo(cfg *dominio.ConfigDB) error {
	m.llamadas = append(m.llamadas, "respaldar_nativo")
	return m.errorFijo
}

func (m *repoMock) Restaurar(cfg *dominio.ConfigDB) error {
	m.llamadas = append(m.llamadas, "restaurar")
	return m.errorFijo
}

func cfgValida(op dominio.Operacion) *dominio.ConfigDB {
	return dominio.CrearConfigDB("mariadb", "3306", "root", "secret", "testdb", "localhost",
		op, "archivo.sql")
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func TestServicio_Init(t *testing.T) {
	mock := &repoMock{}
	svc := servicios.NewServicioDB(mock)

	err := svc.Ejecutar(cfgValida(dominio.OperacionInicializar))

	assert.NoError(t, err)
	assert.Equal(t, []string{"inicializar"}, mock.llamadas)
}

func TestServicio_Backup(t *testing.T) {
	mock := &repoMock{}
	svc := servicios.NewServicioDB(mock)

	err := svc.Ejecutar(cfgValida(dominio.OperacionRespaldar))

	assert.NoError(t, err)
	assert.Equal(t, []string{"respaldar"}, mock.llamadas)
}

func TestServicio_Restore(t *testing.T) {
	mock := &repoMock{}
	svc := servicios.NewServicioDB(mock)

	err := svc.Ejecutar(cfgValida(dominio.OperacionRestaurar))

	assert.NoError(t, err)
	assert.Equal(t, []string{"restaurar"}, mock.llamadas)
}

func TestServicio_OperacionDesconocida(t *testing.T) {
	mock := &repoMock{}
	svc := servicios.NewServicioDB(mock)

	err := svc.Ejecutar(cfgValida("borrar"))
	assert.Error(t, err)
	assert.Empty(t, mock.llamadas)
}

func TestServicio_PropagaErrorRepositorio(t *testing.T) {
	mock := &repoMock{errorFijo: errors.New("fallo de red")}
	svc := servicios.NewServicioDB(mock)

	err := svc.Ejecutar(cfgValida(dominio.OperacionRespaldar))
	assert.Error(t, err)
}

func TestServicio_BackupNativo(t *testing.T) {
	mock := &repoMock{}
	svc := servicios.NewServicioDB(mock)

	err := svc.Ejecutar(cfgValida(dominio.OperacionRespaldarNativo))

	assert.NoError(t, err)
	assert.Equal(t, []string{"respaldar_nativo"}, mock.llamadas)
}

func TestServicio_ConfigInvalida(t *testing.T) {
	mock := &repoMock{}
	svc := servicios.NewServicioDB(mock)

	cfg := dominio.CrearConfigDB("mariadb", "3306", "", "", "", "", "", "")
	err := svc.Ejecutar(cfg)
	assert.Error(t, err)
	assert.Empty(t, mock.llamadas)
}
