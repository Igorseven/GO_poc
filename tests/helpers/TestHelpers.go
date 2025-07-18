package helpers

import (
	entity "PocGo/internal/domain/entities"
	"fmt"
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s\nEsperado: %v\nAtual: %v", message, expected, actual)
	}
}

func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Errorf("%s\nValor é nulo", message)
		return
	}

	// Handle pointer types
	if reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil() {
		t.Errorf("%s\nValor é ponteiro nulo", message)
	}
}

func AssertNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value != nil {
		// Handle non-pointer types
		if reflect.ValueOf(value).Kind() != reflect.Ptr {
			t.Errorf("%s\nValor não é nulo: %v", message, value)
			return
		}

		// Handle pointer types
		if !reflect.ValueOf(value).IsNil() {
			t.Errorf("%s\nValor não é nulo: %v", message, value)
		}
	}
}

func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s\nErro esperado, mas obteve nulo", message)
	}
}

func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s\nErro inesperado: %v", message, err)
	}
}

func CreateTestUser(id string) *entity.User {
	return &entity.User{
		ID:     id,
		Name:   fmt.Sprintf("Test User %s", id),
		Email:  fmt.Sprintf("user%s@example.com", id),
		Status: 1,
	}
}

func CreateTestUsers(count int) *[]entity.User {
	if count <= 0 {
		return &[]entity.User{}
	}

	users := make([]entity.User, count)

	for i := 0; i < count; i++ {
		id := fmt.Sprintf("%d", i+1)
		users[i] = *CreateTestUser(id)
	}
	return &users
}
