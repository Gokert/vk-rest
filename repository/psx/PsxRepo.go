package psx

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
	"vk-rest/configs"
	"vk-rest/pkg/models"
)

type PsxRepo struct {
	db *sql.DB
}

func GetPsxRepo(config *configs.DbPsxConfig, log *logrus.Logger) (*PsxRepo, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password= %s host=%s port=%d sslmode=%s",
		config.User, config.Dbname, config.Password, config.Host, config.Port, config.Sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Errorf("sql open error: %s", err.Error())
		return nil, fmt.Errorf("get user repo err: %s", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Errorf("sql ping error: %s", err.Error())
		return nil, fmt.Errorf("get user repo error: %s", err.Error())
	}
	db.SetMaxOpenConns(config.MaxOpenConns)

	log.Info("Psx created successful")
	return &PsxRepo{db: db}, nil
}

func (repo *PsxRepo) QuestionEvent(event *models.EventItem) error {
	_, err := repo.db.Exec("INSERT INTO quest_on_profile(id_profile, id_quest) VALUES($1, $2)", event.UserId, event.QuestId)
	if err != nil {
		return fmt.Errorf("insert question event error: %s", err.Error())
	}

	return nil
}

func (repo *PsxRepo) QuestionAdd(quest *models.Quest) (uint64, error) {
	var questID uint64
	err := repo.db.QueryRow("INSERT INTO quest(name, cost) VALUES($1, $2) RETURNING id", quest.Name, quest.Cost).Scan(&questID)
	if err != nil {
		return 0, fmt.Errorf("create user error: %s", err.Error())
	}

	return questID, nil
}

func (repo *PsxRepo) GetUser(login string, password []byte) (*models.UserItem, bool, error) {
	post := &models.UserItem{}

	err := repo.db.QueryRow("SELECT profile.id, profile.login, profile.balance FROM profile "+
		"WHERE profile.login = $1 AND profile.password = $2 ", login, password).Scan(&post.Id, &post.Login, &post.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("get query user error: %s", err.Error())
	}

	return post, true, nil
}

func (repo *PsxRepo) FindUser(login string) (bool, error) {
	post := &models.UserItem{}

	err := repo.db.QueryRow(
		"SELECT login FROM profile "+
			"WHERE login = $1", login).Scan(&post.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("find user query error: %s", err.Error())
	}

	return true, nil
}

func (repo *PsxRepo) CreateUser(login string, password []byte) error {
	var userID uint64
	err := repo.db.QueryRow("INSERT INTO profile(login, balance, password) VALUES($1, $2, $3) RETURNING id", login, 0, password).Scan(&userID)
	if err != nil {
		return fmt.Errorf("create user error: %s", err.Error())
	}

	return nil
}

func (repo *PsxRepo) GetUserId(login string) (uint64, error) {
	var userID uint64

	err := repo.db.QueryRow(
		"SELECT profile.id FROM profile WHERE profile.login = $1", login).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found for login: %s", login)
		}
		return 0, fmt.Errorf("get userpro file id error: %s", err.Error())
	}

	return userID, nil
}
