package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Connection struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Estrutura para o Usuário Padrão
type DefaultUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`     // Nome identificador (ex: Conta Admin TI)
	Username string `json:"username"`
	Password string `json:"password"` // Hexadecimal criptografado no banco
}

type DB struct {
	conn *sql.DB
}

func InitDB() (*DB, error) {
	localAppData, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(localAppData, "rdpman")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(appDir, "data.db")
	
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Criação das tabelas
	queryConns := `
	CREATE TABLE IF NOT EXISTS connections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		host TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`
	if _, err := conn.Exec(queryConns); err != nil { return nil, err }

	queryUsers := `
	CREATE TABLE IF NOT EXISTS default_users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`
	if _, err := conn.Exec(queryUsers); err != nil { return nil, err }

	return &DB{conn: conn}, nil
}

// --- Funções CRUD de Connections ---
func (db *DB) GetAll() ([]Connection, error) {
	rows, err := db.conn.Query("SELECT id, name, host, username, password FROM connections ORDER BY name COLLATE NOCASE")
	if err != nil { return nil, err }
	defer rows.Close()

	var conns []Connection
	for rows.Next() {
		var c Connection
		if err := rows.Scan(&c.ID, &c.Name, &c.Host, &c.Username, &c.Password); err != nil { return nil, err }
		conns = append(conns, c)
	}
	return conns, nil
}

func (db *DB) GetByID(id int) (*Connection, error) {
	var c Connection
	err := db.conn.QueryRow("SELECT id, name, host, username, password FROM connections WHERE id = ?", id).
		Scan(&c.ID, &c.Name, &c.Host, &c.Username, &c.Password)
	if err != nil { return nil, err }
	return &c, nil
}

func (db *DB) Insert(c Connection) error {
	_, err := db.conn.Exec("INSERT INTO connections (name, host, username, password) VALUES (?, ?, ?, ?)", c.Name, c.Host, c.Username, c.Password)
	return err
}

func (db *DB) Update(c Connection) error {
	_, err := db.conn.Exec("UPDATE connections SET name = ?, host = ?, username = ?, password = ? WHERE id = ?", c.Name, c.Host, c.Username, c.Password, c.ID)
	return err
}

func (db *DB) Delete(id int) error {
	_, err := db.conn.Exec("DELETE FROM connections WHERE id = ?", id)
	return err
}

// --- Funções CRUD de Default Users ---
func (db *DB) GetDefaultUsers() ([]DefaultUser, error) {
	rows, err := db.conn.Query("SELECT id, name, username, password FROM default_users ORDER BY name COLLATE NOCASE")
	if err != nil { return nil, err }
	defer rows.Close()

	var users []DefaultUser
	for rows.Next() {
		var u DefaultUser
		if err := rows.Scan(&u.ID, &u.Name, &u.Username, &u.Password); err != nil { return nil, err }
		users = append(users, u)
	}
	return users, nil
}

func (db *DB) GetDefaultUserByID(id int) (*DefaultUser, error) {
	var u DefaultUser
	err := db.conn.QueryRow("SELECT id, name, username, password FROM default_users WHERE id = ?", id).
		Scan(&u.ID, &u.Name, &u.Username, &u.Password)
	if err != nil { return nil, err }
	return &u, nil
}

func (db *DB) InsertDefaultUser(u DefaultUser) error {
	_, err := db.conn.Exec("INSERT INTO default_users (name, username, password) VALUES (?, ?, ?)", u.Name, u.Username, u.Password)
	return err
}

func (db *DB) UpdateDefaultUser(u DefaultUser) error {
	_, err := db.conn.Exec("UPDATE default_users SET name = ?, username = ?, password = ? WHERE id = ?", u.Name, u.Username, u.Password, u.ID)
	return err
}

func (db *DB) DeleteDefaultUser(id int) error {
	_, err := db.conn.Exec("DELETE FROM default_users WHERE id = ?", id)
	return err
}

func (db *DB) Close() error {
	return db.conn.Close()
}