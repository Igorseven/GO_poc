package testutils

import (
	config "PocGo/internal/configuration"
	entity "PocGo/internal/domain/entities"
	repository "PocGo/internal/repositories"
	service "PocGo/internal/services"
	dataBaseConnection "PocGo/pkg/database"
	dbProvider "database/sql"
	"sync"
	"testing"
)

type TestApplication struct {
	Configuration *config.Config
	Database      *dbProvider.DB
	Repositories  *repository.Repositories
	Services      *service.Services
	t             *testing.T
	dbInstance    *dataBaseConnection.Database
	cleanupOnce   sync.Once
}

func NewTestApplication(t *testing.T) *TestApplication {
	t.Helper()

	configuration := config.LoadConfig("test")

	dbInstance, err := dataBaseConnection.NewConnection(configuration)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	app := &TestApplication{
		Configuration: configuration,
		Database:      dbInstance.GetConnection(),
		t:             t,
		dbInstance:    dbInstance,
	}

	app.setupRepositories()
	app.setupServices()

	return app
}

func (app *TestApplication) setupRepositories() {
	app.t.Helper()

	repos, err := repository.NewRepositories(app.Database)
	if err != nil {
		app.t.Fatalf("Failed to create repositories: %v", err)
	}
	app.Repositories = repos
}

func (app *TestApplication) setupServices() {
	app.t.Helper()

	app.Services = service.NewServices(app.Repositories)
}

func (app *TestApplication) Cleanup() {
	app.cleanupOnce.Do(func() {
		if app.dbInstance != nil {
			err := app.dbInstance.Close()
			if err != nil {
				app.t.Logf("Warning: Failed to close test database connection: %v", err)
			}
		}
	})
}

func WithTestApplication(t *testing.T, testFunc func(*testing.T, *TestApplication)) {
	t.Helper()

	app := NewTestApplication(t)
	defer app.Cleanup()

	testFunc(t, app)
}

func (app *TestApplication) GetTestUser(t *testing.T, key string) (*entity.User, error) {
	t.Helper()

	testDB := &TestDB{
		DB:            app.Database,
		Config:        app.Configuration,
		UserRepo:      app.Repositories.User,
		ExistingUsers: make(map[string]*entity.User),
	}

	if len(key) == 36 {
		return app.Repositories.User.FindById(key)
	}

	testDB.VerifyTestUsers(t)

	user := testDB.GetTestUser(t, key)
	return user, nil
}

func (app *TestApplication) VerifyTestUsers(t *testing.T) map[string]*entity.User {
	t.Helper()

	testDB := &TestDB{
		DB:            app.Database,
		Config:        app.Configuration,
		UserRepo:      app.Repositories.User,
		ExistingUsers: make(map[string]*entity.User),
	}

	testDB.VerifyTestUsers(t)

	return testDB.ExistingUsers
}
