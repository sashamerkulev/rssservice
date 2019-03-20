package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
)

type UserUpdateRepository interface {
	UpdateUser(userId int64, name string, phone string, logger logger.Logger) error
}

type UserUpdate struct {
	UserId     int64
	Name       string
	Phone      string
	Logger     logger.Logger
	Repository UserUpdateRepository
}

func (uu UserUpdate) UpdateUser() (user model.User, err error) {
	err = uu.Repository.UpdateUser(uu.UserId, uu.Name, uu.Phone, uu.Logger)
	return model.User{UserId: uu.UserId, Name: uu.Name, Phone: uu.Phone}, err
}
