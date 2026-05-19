package mysql

import (
	"controlDB/src/nucleo/dominio"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Repositorio struct {
	Cliente *sql.DB
}

func NewRepositorio(cfg *dominio.ConfigDB) (*Repositorio, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Usuario, cfg.Contraseña, cfg.Host, cfg.Puerto, cfg.BD)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &Repositorio{Cliente: db}, nil
}
