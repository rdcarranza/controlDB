# controlDB
Script para la automatización de operaciones sobre Base de Datos mysql.

# Estructura lógica del proyecto
```
main
  └── cmd/cli/controladores
        ├── cmd/cli/nucleo          (parsea flags)
        ├── src/controladores/env   (lee .env)
        ├── src/nucleo/dominio      (construye cfg)
        ├── src/nucleo/servicios    (ejecuta caso de uso)
        └── src/repositorios/mysql  (inyecta repo)

```


