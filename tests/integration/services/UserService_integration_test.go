package services_test

import (
	entity "PocGo/internal/domain/entities"
	repository "PocGo/internal/repositories"
	service "PocGo/internal/services"
	"PocGo/tests/helpers"
	"PocGo/tests/integration/testutils"
	"testing"
)

func TestUserService_Integration_GetById(t *testing.T) {
	testutils.SkipIfNoDatabase(t)

	testutils.WithTestDB(t, func(t *testing.T, db *testutils.TestDB) {
		// Arrange
		testUser := db.GetTestUser(t, "90FFA97D-110F-4BCE-C6EB-08DDB9C2DAB7")

		userRepo := repository.NewUserRepository(db.DB)
		userService := service.NewUserService(userRepo)

		// Act
		user, err := userService.GetById(testUser.ID)

		// Assert
		helpers.AssertNoError(t, err, "Não deve retornar um erro")
		helpers.AssertNotNil(t, user, "Usuário não deve ser nulo")
		helpers.AssertEqual(t, testUser.ID, user.ID, "ID do usuário deve corresponder")
	})
}

func TestUserService_Integration_GetAll(t *testing.T) {
	testutils.SkipIfNoDatabase(t)

	testutils.WithTestDB(t, func(t *testing.T, db *testutils.TestDB) {
		// Arrange
		testUserKeys := []string{"standard", "admin", "inactive"}
		var testUsers []*entity.User

		for _, key := range testUserKeys {
			if user, exists := db.ExistingUsers[key]; exists {
				testUsers = append(testUsers, user)
			}
		}

		if len(testUsers) < 1 {
			t.Skip("Pulando teste porque não foram encontrados usuários de teste suficientes no banco de dados")
			return
		}

		userRepo := repository.NewUserRepository(db.DB)
		userService := service.NewUserService(userRepo)

		// Act
		users, err := userService.GetAll("")

		// Assert
		helpers.AssertNoError(t, err, "Não deve retornar um erro")
		helpers.AssertNotNil(t, users, "Usuários não devem ser nulos")
		helpers.AssertEqual(t, true, len(*users) >= len(testUsers),
			"Deve ter pelo menos o número de usuários de teste que estamos verificando")

		foundCount := 0
		for _, testUser := range testUsers {
			for _, user := range *users {
				if user.ID == testUser.ID {
					helpers.AssertEqual(t, testUser.Name, user.Name, "Nome do usuário deve corresponder")
					helpers.AssertEqual(t, testUser.Email, user.Email, "Email do usuário deve corresponder")
					helpers.AssertEqual(t, testUser.Status, user.Status, "Status do usuário deve corresponder")
					foundCount++
					break
				}
			}
		}
		helpers.AssertEqual(t, len(testUsers), foundCount, "Todos os usuários de teste devem ser encontrados nos resultados")
	})
}

func TestUserService_Integration_Update(t *testing.T) {
	testutils.SkipIfNoDatabase(t)

	testutils.WithTestDB(t, func(t *testing.T, db *testutils.TestDB) {
		// Arrange
		testUser := db.GetTestUser(t, "standard")

		originalName := testUser.Name
		originalEmail := testUser.Email
		originalStatus := testUser.Status

		userRepo := repository.NewUserRepository(db.DB)
		userService := service.NewUserService(userRepo)

		updateUser := &entity.User{
			ID:     testUser.ID,
			Name:   originalName + " (Updated)",
			Email:  "updated." + originalEmail,
			Status: originalStatus,
		}

		// Act
		err := userService.Update(updateUser)

		// Assert
		helpers.AssertNoError(t, err, "Não deve retornar um erro")

		updatedUser, err := userService.GetById(testUser.ID)
		helpers.AssertNoError(t, err, "Não deve retornar um erro ao recuperar o usuário atualizado")
		helpers.AssertNotNil(t, updatedUser, "Usuário atualizado não deve ser nulo")
		helpers.AssertEqual(t, updateUser.Name, updatedUser.Name, "Nome do usuário deve ser atualizado")
		helpers.AssertEqual(t, updateUser.Email, updatedUser.Email, "Email do usuário deve ser atualizado")
		helpers.AssertEqual(t, updateUser.Status, updatedUser.Status, "Status do usuário deve ser atualizado")

		restoreUser := &entity.User{
			ID:     testUser.ID,
			Name:   originalName,
			Email:  originalEmail,
			Status: originalStatus,
		}

		err = userService.Update(restoreUser)
		if err != nil {
			t.Logf("Aviso: Falha ao restaurar valores originais do usuário: %v", err)
		}
	})
}
