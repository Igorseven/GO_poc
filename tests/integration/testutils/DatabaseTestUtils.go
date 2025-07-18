package testutils

import (
	config "PocGo/internal/configuration"
	entity "PocGo/internal/domain/entities"
	repository "PocGo/internal/repositories"
	dataBaseConnection "PocGo/pkg/database"
	dbProvider "database/sql"
	"log"
	"testing"
)

type TestUser struct {
	ID     string
	Name   string
	Email  string
	Status int
}

var TestUserRegistry = map[string]TestUser{
	"user 1": {
		ID:     "90FFA97D-110F-4BCE-C6EB-08DDB9C2DAB7",
		Name:   "TESTE1@GMAIL.COM",
		Email:  "Teste1@gmail.com",
		Status: 1,
	},
	"user 2": {
		ID:     "D64DF1EC-E446-440C-C6EC-08DDB9C2DAB7",
		Name:   "TESTE2@GMAIL.COM",
		Email:  "Teste2@gmail.com",
		Status: 1,
	},
	"user 3": {
		ID:     "F8E852CA-1D6A-4253-C6ED-08DDB9C2DAB7",
		Name:   "TESTE3@GMAIL.COM",
		Email:  "Teste3@gmail.com",
		Status: 1,
	},
}

type TestDB struct {
	DB            *dbProvider.DB
	Config        *config.Config
	UserRepo      repository.UserRepository
	ExistingUsers map[string]*entity.User
}

func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	configuration := config.LoadConfig("test")

	dbInstance, err := dataBaseConnection.NewConnection(configuration)
	if err != nil {
		t.Fatalf("Falha ao conectar ao banco de dados de teste: %v", err)
	}

	testDB := &TestDB{
		DB:            dbInstance.GetConnection(),
		Config:        configuration,
		ExistingUsers: make(map[string]*entity.User),
	}

	testDB.UserRepo = repository.NewUserRepository(testDB.DB)

	testDB.VerifyTestUsers(t)

	return testDB
}

func (tdb *TestDB) CleanupTestDB(t *testing.T) {
	t.Helper()
	if tdb.DB != nil {
		if err := tdb.DB.Close(); err != nil {
			t.Logf("Aviso: Falha ao fechar conexão com o banco de dados de teste: %v", err)
		}
	}
}

func (tdb *TestDB) ExecuteSQL(t *testing.T, sql string) {
	t.Helper()
	_, err := tdb.DB.Exec(sql)
	if err != nil {
		t.Fatalf("Falha ao executar SQL: %v\nSQL: %s", err, sql)
	}
}

func (tdb *TestDB) VerifyTestUsers(t *testing.T) {
	t.Helper()

	for key, testUser := range TestUserRegistry {
		user, err := tdb.UserRepo.FindById(testUser.ID)

		if err != nil {
			t.Logf("Aviso: Usuário de teste '%s' com ID '%s' não encontrado no banco de dados. Alguns testes podem falhar.", key, testUser.ID)
			continue
		}

		tdb.ExistingUsers[key] = user
		t.Logf("Usuário de teste encontrado '%s' com ID '%s'", key, testUser.ID)
	}

	if len(tdb.ExistingUsers) == 0 {
		t.Logf("Aviso: Nenhum usuário de teste encontrado no banco de dados. Os testes provavelmente falharão.")
	}
}

func (tdb *TestDB) GetTestUser(t *testing.T, key string) *entity.User {
	t.Helper()

	user, exists := tdb.ExistingUsers[key]
	if !exists {
		t.Fatalf("Usuário de teste '%s' não encontrado no mapa ExistingUsers. Certifique-se de que ele existe no banco de dados e foi carregado durante VerifyTestUsers.", key)
	}

	return user
}

func WithTestDB(t *testing.T, testFunc func(*testing.T, *TestDB)) {
	t.Helper()

	testDB := SetupTestDB(t)

	defer func() {
		testDB.CleanupTestDB(t)
	}()

	testFunc(t, testDB)
}

func SkipIfNoDatabase(t *testing.T) {
	t.Helper()

	configuration := config.LoadConfig("test")

	db, err := dataBaseConnection.NewConnection(configuration)
	if err != nil {
		t.Skip("Pulando teste devido a erro de conexão com o banco de dados:", err)
		return
	}

	defer func() {
		if db != nil {
			err := db.Close()
			if err != nil {
				log.Printf("Aviso: Falha ao fechar conexão com o banco de dados: %v", err)
			}
		}
	}()

	err = db.GetConnection().Ping()
	if err != nil {
		t.Skip("Pulando teste devido a erro de ping no banco de dados:", err)
	}

	userRepo := repository.NewUserRepository(db.GetConnection())

	for key, testUser := range TestUserRegistry {
		_, err := userRepo.FindById(testUser.ID)
		if err != nil {
			t.Logf("Aviso: Usuário de teste '%s' com ID '%s' não encontrado no banco de dados.", key, testUser.ID)
		}
	}
}
