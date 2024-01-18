package user

type UserService struct {
	userDAO *DAO
}

func NewUserService(userDAO *DAO) *UserService {
	return &UserService{userDAO: userDAO}
}

func (us *UserService) CreateUser(user User) (*User, error) {
	err := us.userDAO.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) GetUser(userID string) (*User, error) {
	return us.userDAO.GetUser(userID)
}