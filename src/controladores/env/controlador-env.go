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
	if !EnvExiste(dir_env) {
		fmt.Println("El archivo env NO existe, generando desde la plantilla...")
		if err := CrearEnv(dir_env, c_dir_env); err != nil {
			log.Fatal(err)
		}
	}
	return EnvExiste(dir_env)

}

func CrearEnv(dir_env string, c_dir_env string) error {
	copia_env, err := os.Open(c_dir_env)
	if err != nil {
		return fmt.Errorf("no se pudo abrir la plantilla '%s': %w", c_dir_env, err)
	}
	defer copia_env.Close()

	_env, err := os.OpenFile(dir_env, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return fmt.Errorf("no se pudo crear '%s': %w", dir_env, err)
	}
	defer _env.Close()

	if _, err = io.Copy(_env, copia_env); err != nil {
		return fmt.Errorf("error al copiar plantilla: %w", err)
	}

	fmt.Printf("Archivo '%s' generado exitosamente. Completá las variables y volvé a ejecutar.\n", dir_env)
	return nil
}

func EnvExiste(arch_env string) bool {
	if _, err := os.Stat(arch_env); os.IsNotExist(err) {
		return false
	}
	return true

}

func GetEnv(v string, dir_env string) (string, error) {
	file, err := os.Open(dir_env)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir '%s': %w", dir_env, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// ignorar comentarios y líneas vacías
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		partes := strings.SplitN(line, "=", 2)
		if len(partes) == 2 && strings.TrimSpace(partes[0]) == v {
			return strings.TrimSpace(partes[1]), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error al leer '%s': %w", dir_env, err)
	}

	return "", fmt.Errorf("variable '%s' no encontrada en '%s'", v, dir_env)
}
