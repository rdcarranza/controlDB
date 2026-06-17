package main

import "controlDB/cmd/cli/controladores"

func main() {
	controladores.NewControladorCLI().Ejecutar()
}
