package service

import (
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserService interface {
	Register(user *model.User) (model.User, error)          // Registers a new user
	Login(user *model.User) (token *string, err error)      // Logs in the user and returns a JWT token
	GetUserTaskCategory() ([]model.UserTaskCategory, error) // Fetches the user task categories
}

type userService struct {
	userRepo     repo.UserRepository    // User repository
	sessionsRepo repo.SessionRepository // Session repository
}

// NewUserService initializes and returns a new UserService instance
func NewUserService(userRepository repo.UserRepository, sessionsRepo repo.SessionRepository) UserService {
	return &userService{userRepository, sessionsRepo}
}

// Register creates a new user if the email does not already exist
func (s *userService) Register(user *model.User) (model.User, error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return *user, err
	}

	// Check if the email already exists
	if dbUser.Email != "" || dbUser.ID != 0 {
		return *user, errors.New("email already exists")
	}

	// Set the creation time for the new user
	user.CreatedAt = time.Now()

	// Store the new user in the database
	newUser, err := s.userRepo.CreateUser(*user)
	if err != nil {
		return *user, err
	}

	return newUser, nil
}

// Login checks the user's credentials and generates a JWT token if valid
func (s *userService) Login(user *model.User) (token *string, err error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	// Check if the user exists
	if dbUser.Email == "" || dbUser.ID == 0 {
		return nil, errors.New("user not found")
	}

	// Validate the password
	if user.Password != dbUser.Password {
		return nil, errors.New("wrong email or password")
	}

	// Create the JWT token with expiration time of 20 minutes
	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &model.Claims{
		Email: dbUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(model.JwtKey)
	if err != nil {
		return nil, err
	}

	// Create a new session
	session := model.Session{
		Token:  tokenString,
		Email:  user.Email,
		Expiry: expirationTime,
	}

	// Check if the session already exists and handle accordingly
	_, err = s.sessionsRepo.SessionAvailEmail(session.Email)
	if err != nil {
		// If session does not exist, add a new session
		err = s.sessionsRepo.AddSessions(session)
		if err != nil {
			return nil, err
		}
	} else {
		// If session exists, update the existing session
		err = s.sessionsRepo.UpdateSessions(session)
		if err != nil {
			return nil, err
		}
	}

	return &tokenString, nil
}

// GetUserTaskCategory fetches the user task categories from the repository
func (s *userService) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	// Calling the repository to fetch user task categories
	userTaskCategories, err := s.userRepo.GetUserTaskCategory()
	if err != nil {
		return nil, err
	}
	return userTaskCategories, nil
}
