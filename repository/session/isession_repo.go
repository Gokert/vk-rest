package session

import (
	"context"
	"github.com/sirupsen/logrus"
	"vk-rest/pkg/models"
)

type ISessionRepo interface {
	AddSession(ctx context.Context, active models.Session, log *logrus.Logger) (bool, error)
	CheckActiveSession(ctx context.Context, sid string, lg *logrus.Logger) (bool, error)
	GetUserLogin(ctx context.Context, sid string, lg *logrus.Logger) (string, error)
	DeleteSession(ctx context.Context, sid string, lg *logrus.Logger) (bool, error)
}
