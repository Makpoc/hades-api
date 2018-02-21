package models

// User represents an user
type User struct {
	Name         string
	BsRole       string
	BsModules    Modules
	TsModules    Modules
	TsCapacity   string
	MinerLevel   string
	MinerModules Modules
	TZ           UserTime
}

// Users is a list of Users
type Users []User
