#!/bin/bash

APP="controlDB"
echo "Testear ${APP} para arquitectura amd64 - linux"
echo "Test de dominio:"
go test ./src/test/nucleo_test/dominio_test
echo "Test de controladores:"
go test ./src/test/nucleo_test/controlador_test
echo "Test realizado con éxito"