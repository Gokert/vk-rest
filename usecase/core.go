package usecase

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
	"vk-rest/configs"
	utils "vk-rest/pkg"
	"vk-rest/pkg/models"
	"vk-rest/repository/psx"
	"vk-rest/repository/session"
)

type Core struct {
	log      *logrus.Logger
	mutex    sync.RWMutex
	profiles psx.IProfileRepo
	sessions session.ISessionRepo
	question psx.IQuestionRepo
}

func GetCore(psxCfg *configs.DbPsxConfig, redisCfg *configs.DbRedisCfg, log *logrus.Logger) (*Core, error) {
	repo, err := psx.GetPsxRepo(psxCfg, log)
	if err != nil {
		log.Error("Get GetFilmRepo error: ", err)
		return nil, err
	}

	authRepo, err := session.GetAuthRepo(redisCfg, log)
	if err != nil {
		log.Error("Get GetAuthRepo error: ", err)
		return nil, err
	}

	core := &Core{
		log:      log,
		sessions: authRepo,
		profiles: repo,
		question: repo,
	}

	return core, nil
}

func (c *Core) GetUserBalance(userId uint64) (uint64, error) {
	balance, err := c.question.GetUserBalance(userId)
	if err != nil {
		c.log.Errorf("get user balance  error: %s", err.Error())
		return 0, fmt.Errorf("get user balance  error: %s", err.Error())
	}

	return balance, nil
}

func (c *Core) QuestionEvent(event *models.EventItem) error {
	err := c.question.QuestionEvent(event)
	if err != nil {
		c.log.Errorf("question event  error: %s", err.Error())
		return fmt.Errorf("question event error: %s", err.Error())
	}

	return nil
}

func (c *Core) GetUserStat(userId uint64) (*models.UserStat, error) {
	userStat, err := c.question.GetUserStat(userId)
	if err != nil {
		c.log.Errorf("get user stat  error: %s", err.Error())
		return nil, fmt.Errorf("get user stat  error: %s", err.Error())
	}

	balance, err := c.question.GetUserBalance(userId)
	if err != nil {
		c.log.Errorf("get user balance  error: %s", err.Error())
		return nil, fmt.Errorf("get user balance  error: %s", err.Error())
	}

	userStat.Balance = balance

	return userStat, nil
}

func (c *Core) QuestionAdd(quest *models.Quest) (uint64, error) {
	if quest.Name == "" || quest.Cost == 0 {
		return 0, fmt.Errorf("low params")
	}

	questId, err := c.question.QuestionAdd(quest)
	if err != nil {
		c.log.Errorf("question add error: %s", err.Error())
		return 0, fmt.Errorf("question add error: %s", err.Error())
	}

	return questId, nil
}

func (c *Core) GetUserId(ctx context.Context, sid string) (uint64, error) {
	c.mutex.RLock()
	login, err := c.sessions.GetUserLogin(ctx, sid, c.log)
	c.mutex.RUnlock()

	if err != nil {
		c.log.Errorf("get user login error: %s", err.Error())
		return 0, fmt.Errorf("get user login error: %s", err.Error())
	}

	id, err := c.profiles.GetUserId(login)
	if err != nil {
		c.log.Errorf("get user id error: %s", err.Error())
		return 0, fmt.Errorf("get user id error: %s", err.Error())
	}

	return id, nil
}

func (c *Core) GetUserName(ctx context.Context, sid string) (string, error) {
	c.mutex.RLock()
	login, err := c.sessions.GetUserLogin(ctx, sid, c.log)
	c.mutex.RUnlock()

	if err != nil {
		c.log.Errorf("get user name error: %s", err.Error())
		return "", fmt.Errorf("get user name error: %s", err.Error())
	}

	return login, nil
}

func (c *Core) CreateSession(ctx context.Context, login string) (models.Session, error) {
	sid := utils.RandStringRunes(32)

	newSession := models.Session{
		Login:     login,
		SID:       sid,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	c.mutex.Lock()
	sessionAdded, err := c.sessions.AddSession(ctx, newSession, c.log)
	c.mutex.Unlock()

	if !sessionAdded && err != nil {
		return models.Session{}, err
	}

	if !sessionAdded {
		return models.Session{}, nil
	}

	return newSession, nil
}

func (c *Core) FindActiveSession(ctx context.Context, sid string) (bool, error) {
	c.mutex.RLock()
	login, err := c.sessions.CheckActiveSession(ctx, sid, c.log)
	c.mutex.RUnlock()

	if err != nil {
		c.log.Errorf("find active session error: %s", err.Error())
		return false, fmt.Errorf("find active session error: %s", err.Error())
	}

	return login, nil
}

func (c *Core) KillSession(ctx context.Context, sid string) error {
	c.mutex.Lock()
	_, err := c.sessions.DeleteSession(ctx, sid, c.log)
	c.mutex.Unlock()

	if err != nil {
		c.log.Errorf("delete session error: %s", err.Error())
		return fmt.Errorf("delete sessionerror: %s", err.Error())
	}

	return nil
}

func (c *Core) CreateUserAccount(login string, password string) error {
	hashPassword := utils.HashPassword(password)
	err := c.profiles.CreateUser(login, hashPassword)
	if err != nil {
		c.log.Errorf("create user account error: %s", err.Error())
		return fmt.Errorf("create user account error: %s", err.Error())
	}

	return nil
}

func (c *Core) FindUserAccount(login string, password string) (*models.UserItem, bool, error) {
	hashPassword := utils.HashPassword(password)
	user, found, err := c.profiles.GetUser(login, hashPassword)
	if err != nil {
		c.log.Errorf("find user error: %s", err.Error())
		return nil, false, fmt.Errorf("find user account error: %s", err.Error())
	}
	return user, found, nil
}

func (c *Core) FindUserByLogin(login string) (bool, error) {
	found, err := c.profiles.FindUser(login)
	if err != nil {
		c.log.Errorf("find user by login error: %s", err.Error())
		return false, fmt.Errorf("find user by login error: %s", err.Error())
	}

	return found, nil
}
