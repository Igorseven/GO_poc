package services_test

import (
	entity "PocGo/internal/domain/entities"
	"PocGo/tests/helpers"
	"PocGo/tests/integration/testutils"
	"testing"
)

func TestUserService_TestBootstrap_GetById(t *testing.T) {
	testutils.SkipIfNoDatabase(t)

	testutils.WithTestApplication(t, func(t *testing.T, app *testutils.TestApplication) {
		// Arrange
		testUser, err := app.GetTestUser(t, "90FFA97D-110F-4BCE-C6EB-08DDB9C2DAB7")
		if err != nil {
			t.Fatalf("Failed to get test user: %v", err)
		}

		// Act
		user, err := app.Services.User.GetById(testUser.ID)

		// Assert
		helpers.AssertNoError(t, err, "Should not return an error")
		helpers.AssertNotNil(t, user, "User should not be nil")
		helpers.AssertEqual(t, testUser.ID, user.ID, "User ID should match")
	})
}

func TestUserService_TestBootstrap_GetAll(t *testing.T) {
	testutils.SkipIfNoDatabase(t)

	testutils.WithTestApplication(t, func(t *testing.T, app *testutils.TestApplication) {
		// Arrange
		existingUsers := app.VerifyTestUsers(t)
		if len(existingUsers) < 1 {
			t.Skip("Skipping test because not enough test users were found in the database")
			return
		}

		// Act
		users, err := app.Services.User.GetAll("")

		// Assert
		helpers.AssertNoError(t, err, "Should not return an error")
		helpers.AssertNotNil(t, users, "Users should not be nil")
		helpers.AssertEqual(t, true, len(*users) >= len(existingUsers),
			"Should have at least the number of test users we're checking")

		foundCount := 0
		for _, testUser := range existingUsers {
			for _, user := range *users {
				if user.ID == testUser.ID {
					helpers.AssertEqual(t, testUser.Name, user.Name, "User name should match")
					helpers.AssertEqual(t, testUser.Email, user.Email, "User email should match")
					helpers.AssertEqual(t, testUser.Status, user.Status, "User status should match")
					foundCount++
					break
				}
			}
		}
		helpers.AssertEqual(t, len(existingUsers), foundCount, "All test users should be found in the results")
	})
}

func TestUserService_TestBootstrap_Update(t *testing.T) {
	testutils.SkipIfNoDatabase(t)

	testutils.WithTestApplication(t, func(t *testing.T, app *testutils.TestApplication) {
		// Arrange
		existingUsers := app.VerifyTestUsers(t)

		var testUser *entity.User
		for _, user := range existingUsers {
			if user.Status == 1 {
				testUser = user
				break
			}
		}

		if testUser == nil {
			t.Skip("Skipping test because no suitable test user was found")
			return
		}

		originalName := testUser.Name
		originalEmail := testUser.Email
		originalStatus := testUser.Status

		updateUser := &entity.User{
			ID:     testUser.ID,
			Name:   originalName + " (Updated)",
			Email:  "updated." + originalEmail,
			Status: originalStatus,
		}

		// Act
		err := app.Services.User.Update(updateUser)

		// Assert
		helpers.AssertNoError(t, err, "Should not return an error")

		updatedUser, err := app.Services.User.GetById(testUser.ID)
		helpers.AssertNoError(t, err, "Should not return an error when retrieving the updated user")
		helpers.AssertNotNil(t, updatedUser, "Updated user should not be nil")
		helpers.AssertEqual(t, updateUser.Name, updatedUser.Name, "User name should be updated")
		helpers.AssertEqual(t, updateUser.Email, updatedUser.Email, "User email should be updated")
		helpers.AssertEqual(t, updateUser.Status, updatedUser.Status, "User status should be updated")

		restoreUser := &entity.User{
			ID:     testUser.ID,
			Name:   originalName,
			Email:  originalEmail,
			Status: originalStatus,
		}

		err = app.Services.User.Update(restoreUser)
		if err != nil {
			t.Logf("Warning: Failed to restore original user values: %v", err)
		}
	})
}
