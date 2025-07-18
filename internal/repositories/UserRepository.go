package repositories

import (
	converter "PocGo/internal/configuration/converters"
	entity "PocGo/internal/domain/entities"
	notify "PocGo/internal/domain/notification"
	dbProvider "database/sql"
	"errors"
)

const (
	findByIdQuery     = `SELECT [id], [normalized_login], [login], [status] FROM [Auth].[User] WHERE id = @p1`
	findAllQuery      = `SELECT [id], [normalized_login], [login], [status] FROM [Auth].[User] WHERE creation_date > @p1`
	updateQuery       = `UPDATE [Auth].[User] SET name = @p1, email = @p2, status = @p3 WHERE id = @p4`
	findOldUsersQuery = `SELECT [id], [normalized_login], [login], [status] FROM [Auth].[User] WHERE creation_date < DATEADD(month, -5, GETDATE())`
	updateStatusQuery = `UPDATE [Auth].[User] SET status = @p1 WHERE id = @p2`
)

type UserRepository interface {
	Update(user *entity.User) error
	FindById(id string) (*entity.User, error)
	FindAll(date string) (*[]entity.User, error)
	FindOldUsers() (*[]entity.User, error)
	UpdateStatus(id string, status int) error
}

type userRepository struct {
	dataBase *dbProvider.DB
}

func NewUserRepository(dbProvider *dbProvider.DB) UserRepository {
	return &userRepository{
		dataBase: dbProvider,
	}
}

type userTemp struct {
	ID     []byte
	Name   string
	Email  string
	Status int
}

func mapToUser(temp userTemp) entity.User {
	return entity.User{
		ID:     converter.BytesToString(temp.ID),
		Name:   temp.Name,
		Email:  temp.Email,
		Status: temp.Status,
	}
}

func (r *userRepository) FindById(id string) (*entity.User, error) {
	var temp userTemp

	err := r.dataBase.QueryRow(findByIdQuery, id).Scan(
		&temp.ID,
		&temp.Name,
		&temp.Email,
		&temp.Status,
	)

	if err != nil {
		if errors.Is(err, dbProvider.ErrNoRows) {
			return nil, notify.CreateSimpleNotification(notify.NotFound, err)
		}
		return nil, notify.CreateSimpleNotification(notify.FindErrorRepository, err)
	}

	user := mapToUser(temp)
	return &user, nil
}

func (r *userRepository) FindAll(date string) (*[]entity.User, error) {
	rows, err := r.dataBase.Query(findAllQuery, date)
	if err != nil {
		return nil, notify.CreateSimpleNotification(notify.FindErrorRepository, err)
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, notify.CreateSimpleNotification(notify.FindAllErrorRepository, err)
	}

	safeUsers := converter.ListSafe(users)
	return &safeUsers, nil
}

func (r *userRepository) Update(user *entity.User) error {
	result, err := r.dataBase.Exec(updateQuery, user.Name, user.Email, user.Status, user.ID)
	if err != nil {
		return notify.CreateSimpleNotification(notify.InvalidData, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return notify.CreateSimpleNotification(notify.InvalidData, err)
	}

	if rowsAffected == 0 {
		return notify.CreateSimpleNotification(notify.NotFound, nil)
	}

	return nil
}

func (r *userRepository) FindOldUsers() (*[]entity.User, error) {
	rows, err := r.dataBase.Query(findOldUsersQuery)
	if err != nil {
		return nil, notify.CreateSimpleNotification(notify.FindErrorRepository, err)
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, notify.CreateSimpleNotification(notify.FindAllErrorRepository, err)
	}

	safeUsers := converter.ListSafe(users)
	return &safeUsers, nil
}

func (r *userRepository) UpdateStatus(id string, status int) error {
	result, err := r.dataBase.Exec(updateStatusQuery, status, id)
	if err != nil {
		return notify.CreateSimpleNotification(notify.InvalidData, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return notify.CreateSimpleNotification(notify.InvalidData, err)
	}

	if rowsAffected == 0 {
		return notify.CreateSimpleNotification(notify.NotFound, nil)
	}

	return nil
}

func scanUsers(rows *dbProvider.Rows) ([]entity.User, error) {
	var users []entity.User

	for rows.Next() {
		var temp userTemp

		if err := rows.Scan(
			&temp.ID,
			&temp.Name,
			&temp.Email,
			&temp.Status,
		); err != nil {
			return nil, notify.CreateSimpleNotification(notify.ScanErrorRepository, err)
		}

		user := mapToUser(temp)
		users = append(users, user)
	}

	return users, nil
}
