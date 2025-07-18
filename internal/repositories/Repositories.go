package repositories

import (
	sqlServer "database/sql"
)

type Repositories struct {
	db   *sqlServer.DB
	User UserRepository
	// Outros repositórios aqui
}

func NewRepositories(db *sqlServer.DB) (*Repositories, error) {
	repos := &Repositories{
		db: db,
	}

	repos.User = NewUserRepository(db)

	return repos, nil
}

func (repos *Repositories) GetDB() *sqlServer.DB {
	return repos.db
}
