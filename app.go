package main

import (
	"context"
	"fmt"
	"rdpman/backend/crypto"
	"rdpman/backend/database"
	"rdpman/backend/rdp"
)

type App struct {
	ctx context.Context
	db  *database.DB
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	dbInstance, err := database.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Falha ao iniciar banco de dados: %v", err))
	}
	a.db = dbInstance
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// --- Métodos das Conexões ---

func (a *App) GetConnections() ([]database.Connection, error) {
	conns, err := a.db.GetAll()
	if err != nil {
		return nil, err
	}
	for i := range conns {
		conns[i].Password = ""
		
		decHost, _ := crypto.Decrypt(conns[i].Host)
		if decHost != "" { conns[i].Host = decHost }
		
		decUser, _ := crypto.Decrypt(conns[i].Username)
		if decUser != "" { conns[i].Username = decUser }
	}
	return conns, nil
}

func (a *App) SaveConnection(c database.Connection, rawPassword string, updatePassword bool) error {
	encHost, err := crypto.Encrypt(c.Host)
	if err == nil { c.Host = encHost }
	
	encUser, err := crypto.Encrypt(c.Username)
	if err == nil { c.Username = encUser }

	if updatePassword {
		enc, err := crypto.Encrypt(rawPassword)
		if err != nil {
			return err
		}
		c.Password = enc
	} else if c.ID > 0 {
		oldConn, err := a.db.GetByID(c.ID)
		if err == nil {
			c.Password = oldConn.Password
		}
	}

	if c.ID == 0 {
		return a.db.Insert(c)
	}
	return a.db.Update(c)
}

func (a *App) DeleteConnection(id int) error {
	return a.db.Delete(id)
}

func (a *App) Connect(id int) error {
	c, err := a.db.GetByID(id)
	if err != nil {
		return err
	}
	
	decHost, _ := crypto.Decrypt(c.Host)
	decUser, _ := crypto.Decrypt(c.Username)
	
	return rdp.Connect(decHost, decUser, c.Password)
}

// --- Métodos dos Usuários Padrões ---

func (a *App) GetDefaultUsers() ([]database.DefaultUser, error) {
	users, err := a.db.GetDefaultUsers()
	if err != nil {
		return nil, err
	}
	// Descriptografa tudo para mandar à memória do Frontend e facilitar o preenchimento automático
	for i := range users {
		decUser, _ := crypto.Decrypt(users[i].Username)
		if decUser != "" { users[i].Username = decUser }
		
		decPass, _ := crypto.Decrypt(users[i].Password)
		if decPass != "" { users[i].Password = decPass }
	}
	return users, nil
}

func (a *App) SaveDefaultUser(u database.DefaultUser, rawPassword string, updatePassword bool) error {
	encUser, err := crypto.Encrypt(u.Username)
	if err == nil { u.Username = encUser }

	if updatePassword {
		encPass, err := crypto.Encrypt(rawPassword)
		if err == nil { u.Password = encPass }
	} else if u.ID > 0 {
		oldUser, err := a.db.GetDefaultUserByID(u.ID)
		if err == nil {
			u.Password = oldUser.Password
		}
	}

	if u.ID == 0 {
		return a.db.InsertDefaultUser(u)
	}
	return a.db.UpdateDefaultUser(u)
}

func (a *App) DeleteDefaultUser(id int) error {
	return a.db.DeleteDefaultUser(id)
}