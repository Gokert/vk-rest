package psx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
	"time"
	"vk-rest/configs"
	utils "vk-rest/pkg"
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

	repo := &PsxRepo{db: db}

	errs := make(chan error)
	go func() {
		errs <- repo.pingDb(3, log)
	}()

	if err := <-errs; err != nil {
		log.Error(err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	log.Info("Successfully connected to database")

	return repo, nil
}

func (repo *PsxRepo) pingDb(timer uint32, log *logrus.Logger) error {
	var err error
	var retries int

	for retries < utils.MaxRetries {
		err = repo.db.Ping()
		if err == nil {
			return nil
		}

		retries++
		log.Errorf("sql ping error: %s", err.Error())
		time.Sleep(time.Duration(timer) * time.Second)
	}

	return fmt.Errorf("sql max pinging error: %s", err.Error())
}

func (repo *PsxRepo) GetUserBalance(ctx context.Context, userId uint64) (uint64, error) {
	var total uint64 = 0

	rows, err := repo.db.QueryContext(ctx, `
	select profile.id, SUM(quest.cost) FROM quest
	LEFT JOIN quest_on_profile ON quest.id = quest_on_profile.id_quest
	LEFT JOIN profile ON profile.id = quest_on_profile.id_profile
	WHERE profile.id = $1
	GROUP BY profile.id`, userId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("sql get user balance error: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id uint64
		var balance uint64

		err := rows.Scan(&id, &balance)
		if err != nil {
			return 0, fmt.Errorf("scan get user balance error: %s", err.Error())
		}
		total += balance
	}

	return total, nil
}

func (repo *PsxRepo) GetUserStat(ctx context.Context, userId uint64) (*models.UserStat, error) {
	post := models.UserStat{}

	rows, err := repo.db.QueryContext(ctx, `
	select quest.id, quest.name, quest.cost FROM quest
	LEFT JOIN quest_on_profile ON quest.id = quest_on_profile.id_quest
	LEFT JOIN profile ON profile.id = quest_on_profile.id_profile
	WHERE profile.id = $1`, userId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("sql get user stat error: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var quest models.Quest

		err := rows.Scan(&quest.Id, &quest.Name, &quest.Cost)
		if err != nil {
			return nil, fmt.Errorf("scan get user stat: %s", err.Error())
		}
		post.Quests = append(post.Quests, quest)
	}

	return &post, nil
}

func (repo *PsxRepo) QuestionEvent(ctx context.Context, event *models.EventItem) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO quest_on_profile(id_profile, id_quest) VALUES($1, $2)", event.UserId, event.QuestId)
	if err != nil {
		return fmt.Errorf("insert question event error: %s", err.Error())
	}

	return nil
}

func (repo *PsxRepo) QuestionAdd(ctx context.Context, quest *models.Quest) (uint64, error) {
	var questID uint64
	err := repo.db.QueryRowContext(ctx, "INSERT INTO quest(name, cost) VALUES($1, $2) RETURNING id", quest.Name, quest.Cost).Scan(&questID)
	if err != nil {
		return 0, fmt.Errorf("question add error: %s", err.Error())
	}

	return questID, nil
}

func (repo *PsxRepo) GetUser(ctx context.Context, login string, password []byte) (*models.UserItem, bool, error) {
	post := &models.UserItem{}

	err := repo.db.QueryRowContext(ctx, "SELECT profile.id, profile.login, profile.balance FROM profile "+
		"WHERE profile.login = $1 AND profile.password = $2 ", login, password).Scan(&post.Id, &post.Login, &post.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("get query user error: %s", err.Error())
	}

	return post, true, nil
}

func (repo *PsxRepo) FindUser(ctx context.Context, login string) (bool, error) {
	post := &models.UserItem{}

	err := repo.db.QueryRowContext(ctx,
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

func (repo *PsxRepo) CreateUser(ctx context.Context, login string, password []byte) error {
	var userID uint64
	err := repo.db.QueryRowContext(ctx, "INSERT INTO profile(login, balance, password) VALUES($1, $2, $3) RETURNING id", login, 0, password).Scan(&userID)
	if err != nil {
		return fmt.Errorf("create user error: %s", err.Error())
	}

	return nil
}

func (repo *PsxRepo) GetUserId(ctx context.Context, login string) (uint64, error) {
	var userID uint64

	err := repo.db.QueryRowContext(ctx,
		"SELECT profile.id FROM profile WHERE profile.login = $1", login).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found for login: %s", login)
		}
		return 0, fmt.Errorf("get userpro file id error: %s", err.Error())
	}

	return userID, nil
}
