package psx

import "vk-rest/pkg/models"

type IProfileRepo interface {
	GetUser(login string, password []byte) (*models.UserItem, bool, error)
	FindUser(login string) (bool, error)
	CreateUser(login string, password []byte) error
	GetUserId(login string) (uint64, error)
}
