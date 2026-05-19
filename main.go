package main

import (
	"controlDB/src/controladores/env"
	"controlDB/src/nucleo/dominio"
	"controlDB/src/repositorios/mysql"
	"fmt"
)

func main() {
	if env.VerificarEnv(".env", "./src/controladores/env/env.copia") {

		host, err := env.GetEnv("host_db", ".env")
		name_db, err := env.GetEnv("name_db", ".env")
		port_db, err := env.GetEnv("port_db", ".env")
		user_db, err := env.GetEnv("user_db", ".env")
		pw_db, err := env.GetEnv("pw_db", ".env")
		if err != nil {
			fmt.Println(err)
			return
		}
		cfg := dominio.CrearConfigDB("mariadb", port_db, user_db, pw_db, name_db, host)
		fmt.Printf("configuración DB: %+v", cfg)
		repo, err := mysql.NewRepositorio(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(repo)
	} else {
		fmt.Println("No se encontró el archivo env, verificar archivo y ajustar las variables de entorno!")

	}

}
