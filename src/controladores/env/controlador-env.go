package env

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func VerificarEnv(dir_env string, c_dir_env string) bool {
	/*
		ex, _ := os.Executable()
		fmt.Println("La ruta del ejecutable es: " + ex)
		exPath := filepath.Dir(ex)
		fmt.Println("El directorio del ejecutable es:" + exPath)
	*/

	env := dir_env
	copia_env := c_dir_env
	if !EnvExiste(env) {
		fmt.Println("El archivo env NO existe!")
		err := CrearEnv(env, copia_env)
		if err != nil {
			log.Fatal(err)
		}

	}

	if EnvExiste(env) {
		return true
	}

	return false

}

func CrearEnv(dir_env string, c_dir_env string) error {
	copia_env, err := os.Open(c_dir_env)
	if err != nil {
		log.Fatal(err)
	}
	_env, err := os.OpenFile(dir_env, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
		//log.Fatal(err)
	}
	_, err = io.Copy(_env, copia_env)
	if err != nil {
		return err
		//log.Fatal(err)

	} else {
		fmt.Println("Se genera env exitosamente!")
	}

	return nil
}

func EnvExiste(arch_env string) bool {
	if _, err := os.Stat(arch_env); os.IsNotExist(err) {
		return false
	}
	return true

}

func GetEnv(v string, dir_env string) (string, error) {
	env := dir_env

	file, err := os.Open(env)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	variable := v
	scanner := bufio.NewScanner(file)
	linea := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, variable) {
			fmt.Printf("Se encontró la variable '%s' en la línea: %s\n", variable, line)
			linea = line
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return "", err
	}

	return strings.Replace(linea, v+"=", "", -1), nil
}
