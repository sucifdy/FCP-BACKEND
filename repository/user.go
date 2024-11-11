package repository

import (
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/model"
	"fmt"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
}

type userRepository struct {
	filebasedDb *filebased.Data
}

func NewUserRepo(filebasedDb *filebased.Data) *userRepository {
	return &userRepository{filebasedDb}
}

// GetUserByEmail retrieves a user by their email address.
func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	user, err := r.filebasedDb.GetUserByEmail(email)
	if err != nil {
		return model.User{}, fmt.Errorf("error finding user with email %s: %v", email, err)
	}
	return user, nil
}

// CreateUser creates a new user in the database.
func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	createdUser, err := r.filebasedDb.CreateUser(user)
	if err != nil {
		return model.User{}, fmt.Errorf("error creating user: %v", err)
	}
	return createdUser, nil
}

// GetUserTaskCategory retrieves a list of user task categories from the database.
func (r *userRepository) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	userTaskCategories, err := r.filebasedDb.GetUserTaskCategory()
	if err != nil {
		return nil, fmt.Errorf("error retrieving user task categories: %v", err)
	}
	return userTaskCategories, nil
}
