package mysql

import (
	"bufio"
	"controlDB/src/nucleo/dominio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Repositorio struct {
	Cliente *sql.DB
}

func NewRepositorio(cfg *dominio.ConfigDB) (*Repositorio, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %w", err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)
	return &Repositorio{Cliente: db}, nil
}

// ─── Operaciones ──────────────────────────────────────────────────────────────

// Inicializar crea la BD si no existe y ejecuta el archivo SQL indicado.
func (r *Repositorio) Inicializar(cfg *dominio.ConfigDB) error {
	// Conectar sin seleccionar BD para poder crearla
	db, err := sql.Open("mysql", cfg.DSNSinBD())
	if err != nil {
		return fmt.Errorf("no se pudo conectar al servidor: %w", err)
	}
	defer db.Close()

	q := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		cfg.BD,
	)
	if _, err = db.Exec(q); err != nil {
		return fmt.Errorf("no se pudo crear la BD: %w", err)
	}

	if _, err = db.Exec("USE `" + cfg.BD + "`"); err != nil {
		return fmt.Errorf("no se pudo seleccionar la BD: %w", err)
	}

	return ejecutarArchivoSQL(db, cfg.ArchivoSQL)
}

// Respaldar exporta la BD a un archivo SQL.
// Usa mysqldump si está disponible en el PATH; si no, hace exportación nativa Go.
func (r *Repositorio) Respaldar(cfg *dominio.ConfigDB) error {
	if ruta, err := exec.LookPath("mysqldump"); err == nil {
		return respaldarConDump(ruta, cfg)
	}
	return respaldarManual(cfg)
}

// RespaldarNativo usa mariadb-backup (o xtrabackup) para hacer un respaldo
// binario físico del servidor completo en el directorio indicado por ArchivoSQL.
//
// El directorio de destino debe estar vacío o no existir.
// mariadb-backup debe estar instalado y disponible en el PATH.
//
// Equivale a ejecutar:
//
//	mariadb-backup --backup --target-dir=<dir> --user=<u> --password=<p> --host=<h> --port=<p>
func (r *Repositorio) RespaldarNativo(cfg *dominio.ConfigDB) error {
	herramienta, err := buscarHerramientaNativa()
	if err != nil {
		return err
	}

	// mariadb-backup trabaja sobre directorios, no archivos .sql
	destino := cfg.ArchivoSQL
	if err := os.MkdirAll(destino, 0750); err != nil {
		return fmt.Errorf("no se pudo crear el directorio de destino '%s': %w", destino, err)
	}

	args := []string{
		"--backup",
		"--target-dir=" + destino,
		"--host=" + cfg.Host,
		"--port=" + cfg.Puerto,
		"--user=" + cfg.Usuario,
		"--password=" + cfg.Contraseña,
	}

	cmd := exec.Command(herramienta, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s falló: %w", herramienta, err)
	}
	return nil
}

// buscarHerramientaNativa busca mariadb-backup y xtrabackup en el PATH.
// mariadb-backup tiene precedencia por ser la herramienta oficial de MariaDB.
func buscarHerramientaNativa() (string, error) {
	for _, nombre := range []string{"mariadb-backup", "xtrabackup"} {
		if ruta, err := exec.LookPath(nombre); err == nil {
			return ruta, nil
		}
	}
	return "", fmt.Errorf(
		"no se encontró mariadb-backup ni xtrabackup en el PATH\n" +
			"  En Debian/Ubuntu: sudo apt install mariadb-backup\n" +
			"  En RHEL/Rocky:    sudo dnf install MariaDB-backup",
	)
}

// Restaurar importa un archivo SQL sobre la BD existente.
func (r *Repositorio) Restaurar(cfg *dominio.ConfigDB) error {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return fmt.Errorf("no se pudo conectar: %w", err)
	}
	defer db.Close()
	return ejecutarArchivoSQL(db, cfg.ArchivoSQL)
}

// ─── helpers internos ─────────────────────────────────────────────────────────

// ejecutarArchivoSQL lee el .sql y ejecuta cada sentencia individualmente.
func ejecutarArchivoSQL(db *sql.DB, ruta string) error {
	f, err := os.Open(ruta)
	if err != nil {
		return fmt.Errorf("no se pudo abrir '%s': %w", ruta, err)
	}
	defer f.Close()

	var sentencia strings.Builder
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	nLinea := 0
	for scanner.Scan() {
		nLinea++
		texto := strings.TrimSpace(scanner.Text())

		if texto == "" || strings.HasPrefix(texto, "--") || strings.HasPrefix(texto, "#") {
			continue
		}

		sentencia.WriteString(texto)
		sentencia.WriteString("\n")

		if strings.HasSuffix(texto, ";") {
			sql := strings.TrimSpace(sentencia.String())
			if sql != "" && sql != ";" {
				if _, err := db.Exec(sql); err != nil {
					return fmt.Errorf("error en línea %d: %w\nSQL: %s", nLinea, err, sql)
				}
			}
			sentencia.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error al leer el archivo: %w", err)
	}

	// sentencia final sin ; (si la hay)
	if resto := strings.TrimSpace(sentencia.String()); resto != "" {
		if _, err := db.Exec(resto); err != nil {
			return fmt.Errorf("error al ejecutar sentencia final: %w", err)
		}
	}

	return nil
}

// respaldarConDump usa mysqldump para la exportación.
func respaldarConDump(dumpPath string, cfg *dominio.ConfigDB) error {
	args := []string{
		"--host=" + cfg.Host,
		"--port=" + cfg.Puerto,
		"--user=" + cfg.Usuario,
		"--password=" + cfg.Contraseña,
		"--single-transaction",
		"--routines",
		"--triggers",
		"--add-drop-table",
		cfg.BD,
	}

	out, err := os.Create(cfg.ArchivoSQL)
	if err != nil {
		return fmt.Errorf("no se pudo crear '%s': %w", cfg.ArchivoSQL, err)
	}
	defer out.Close()

	cmd := exec.Command(dumpPath, args...)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysqldump falló: %w", err)
	}
	return nil
}

// respaldarManual exporta estructura + datos usando solo el driver Go.
func respaldarManual(cfg *dominio.ConfigDB) error {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return err
	}
	defer db.Close()

	out, err := os.Create(cfg.ArchivoSQL)
	if err != nil {
		return fmt.Errorf("no se pudo crear '%s': %w", cfg.ArchivoSQL, err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	fmt.Fprintf(w, "-- Respaldo generado por controlDB\n-- Fecha: %s\n-- BD: %s\n\n",
		time.Now().Format(time.RFC3339), cfg.BD)
	fmt.Fprintln(w, "SET FOREIGN_KEY_CHECKS=0;")

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return fmt.Errorf("error al listar tablas: %w", err)
	}
	defer rows.Close()

	var tablas []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return err
		}
		tablas = append(tablas, t)
	}

	for _, tabla := range tablas {
		if err := volcarTabla(db, w, tabla); err != nil {
			return fmt.Errorf("error al volcar '%s': %w", tabla, err)
		}
	}

	fmt.Fprintln(w, "\nSET FOREIGN_KEY_CHECKS=1;")
	return w.Flush()
}

func volcarTabla(db *sql.DB, w *bufio.Writer, tabla string) error {
	var nombre, ddl string
	if err := db.QueryRow("SHOW CREATE TABLE `"+tabla+"`").Scan(&nombre, &ddl); err != nil {
		return err
	}
	fmt.Fprintf(w, "\nDROP TABLE IF EXISTS `%s`;\n%s;\n", tabla, ddl)

	rows, err := db.Query("SELECT * FROM `" + tabla + "`")
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	vals := make([]interface{}, len(cols))
	pvals := make([]interface{}, len(cols))
	for i := range vals {
		pvals[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(pvals...); err != nil {
			return err
		}
		partes := make([]string, len(cols))
		for i, v := range vals {
			partes[i] = valorSQL(v)
		}
		fmt.Fprintf(w, "INSERT INTO `%s` VALUES (%s);\n", tabla, strings.Join(partes, ", "))
	}
	return rows.Err()
}

func valorSQL(v interface{}) string {
	if v == nil {
		return "NULL"
	}
	switch val := v.(type) {
	case []byte:
		return "'" + escaparSQL(string(val)) + "'"
	case string:
		return "'" + escaparSQL(val) + "'"
	case time.Time:
		return "'" + val.Format("2006-01-02 15:04:05") + "'"
	default:
		return fmt.Sprintf("%v", val)
	}
}

func escaparSQL(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	return s
}
