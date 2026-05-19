package dominio

type ConfigDB struct {
	Tipo       string
	Puerto     string
	Usuario    string
	Contraseña string
	BD         string
	Host       string
}

func CrearConfigDB(tipo string, puerto string, usuario string, contraseña string, bd string, host string) *ConfigDB {
	return &ConfigDB{
		Tipo:       tipo,
		Puerto:     puerto,
		Usuario:    usuario,
		Contraseña: contraseña,
		BD:         bd,
		Host:       host,
	}
}
