package services

import (
	"github.com/roblieblang/luthien-core-server/internal/dao"
	"github.com/roblieblang/luthien-core-server/internal/models"
)

type UserService struct {
	userDAO *dao.UserDAO
}

func NewUserService(userDAO *dao.UserDAO) *UserService {
	return &UserService{userDAO: userDAO}
}

func (us *UserService) CreateUser(user models.User) (*models.User, error) {
	err := us.userDAO.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) GetUser(userID int) (*models.User, error) {
	return us.userDAO.GetUser(userID)
}