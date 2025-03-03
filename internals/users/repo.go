package users

import (
	"concurrent_money_transfer_system/utils"
	"sync"
	"time"
)

type UserRepo interface {
	CreateUser(user User) (User, error)
	GetUser(id string) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id string) error
	GetAllUsers() ([]User, error)
}

type userRepo struct {
	users sync.Map
}

func (r *userRepo) CreateUser(user User) (User, error) {
	if user.ID == "" {
		user.ID = utils.GenerateUniqueEntityId()
	}
	if _, ok := r.users.Load(user.ID); ok {
		return User{}, utils.NewError(utils.ErrUserAlreadyExists)
	}
	r.users.Store(user.ID, user)
	return user, nil	
}

// TODO: Handle error if user not found
func (r *userRepo) GetUser(id string) (User, error) {
	user, ok := r.users.Load(id)
	if !ok {
		return User{}, utils.NewError(utils.ErrUserNotFound)
	}
	if user.(User).DeletedAt != nil {
		return User{}, utils.NewError(utils.ErrUserNotFound)
	}
	return user.(User), nil
}

func (r *userRepo) UpdateUser(user User) (User, error) {
	_, err := r.GetUser(user.ID)
	if err != nil {
		return User{}, err
	}
	r.users.Store(user.ID, user)
	return user, nil
}

func (r *userRepo) DeleteUser(id string) error {
	user, ok := r.users.Load(id)
	if !ok {
		return utils.NewError(utils.ErrUserNotFound)
	}

	deletedAt := time.Now()

	u := user.(User)
	u.DeletedAt = &deletedAt
	r.users.Store(id, u)
	return nil
}

func (r *userRepo) GetUserByEmail(email string) (User, error) {
	var foundUser User
	var found bool
	
	r.users.Range(func(key, value interface{}) bool {
		if user, ok := value.(User); ok && user.Email == email {
			foundUser = user
			found = true
			return false // Stop iteration
		}
		return true // Continue iteration
	})
	
	if !found {
		return User{}, utils.NewError(utils.ErrUserNotFound)
	}
	
	return foundUser, nil
}

func (r *userRepo) GetAllUsers() ([]User, error) {
	users := make([]User, 0)
	r.users.Range(func(key, value interface{}) bool {
		if user, ok := value.(User); ok {
			users = append(users, user)
		}
		return true
	})
	return users, nil
}


var userRepoInstance *userRepo

func NewUserRepo() UserRepo {
	if userRepoInstance == nil {
		userRepoInstance = &userRepo{
			users: sync.Map{},
		}
	}
	return userRepoInstance
}
