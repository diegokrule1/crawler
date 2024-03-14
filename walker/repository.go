package walker

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type UrlRepository struct {
	Pool   *sql.DB
	Logger *zap.Logger
}

func (u *UrlRepository) CreateUrl(scheme, host, path string, id string, parentId *string) (int64, error) {
	u.Logger.Info(fmt.Sprintf("Persisting url %s/%s/%s", scheme, host, path))
	cont := context.Background()

	tx, err := u.Pool.BeginTx(cont, nil)
	if err != nil {
		u.Logger.Error("Could not create tx", zap.Error(err))
		return 0, err
	}

	//if err != nil {
	//	u.Logger.Error("could not create id", zap.Error(err))
	//	return 0, err
	//}

	rows, err := tx.Exec("insert into url (id, scheme,host, path, insertion_date, state, parent) values ($1, $2, $3, $4,$5, $6, $7)", id, scheme, host, path, time.Now(), Created, parentId)

	if err != nil {
		u.Logger.Error("Could not execute insert into url table", zap.Error(err))
		return 0, nil
	}

	err = tx.Commit()
	if err != nil {
		u.Logger.Error("could not commit tx", zap.Error(err))
		return 0, err
	}

	total, err := rows.RowsAffected()

	if err != nil {
		u.Logger.Error("Could not get total rows affected", zap.Error(err))
		return 0, err
	}

	u.Logger.Debug(fmt.Sprintf("Total rows inserted %d", total))
	return total, nil
}

func (u *UrlRepository) getUrlToProcess(urlId string) (*Url, error) {
	u.Logger.Info("Getting url to analyze")

	query := "select id, scheme, host, path from url where id=$1 limit 1"
	rows, err := u.Pool.Query(query, urlId)
	defer func(logger *zap.Logger) {
		err := rows.Close()
		if err != nil {
			logger.Error("could not close row", zap.Error(err))
		}
	}(u.Logger)

	if err != nil {
		return nil, err
	}
	var id, scheme, host, path string
	for rows.Next() {
		err := rows.Scan(&id, &scheme, &host, &path)
		if err != nil {
			u.Logger.Error("Could not populate url struct ", zap.Error(err))
			return nil, err
		}
	}
	return &Url{Id: id, Scheme: scheme, Host: host, Path: path}, nil
}

func (u *UrlRepository) changeState(id string, state State) (int64, error) {
	cont := context.Background()

	tx, err := u.Pool.BeginTx(cont, nil)
	if err != nil {
		u.Logger.Error("Could not create tx", zap.Error(err))
		return 0, err
	}

	rows, err := tx.Exec("update url set state=$1 where id=$2", state, id)

	if err != nil {
		u.Logger.Error("Could not execute update into url table", zap.Error(err))
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		u.Logger.Error("could not commit tx", zap.Error(err))
		return 0, err
	}

	total, err := rows.RowsAffected()

	if err != nil {
		u.Logger.Error("Could not get total rows affected", zap.Error(err))
		return 0, err
	}

	u.Logger.Debug(fmt.Sprintf("Total rows updated %d", total))
	return total, nil
}

func (u *UrlRepository) countNonterminated() (int64, error) {
	u.Logger.Info("Getting url to analyze")

	query := fmt.Sprintf("select count(*) from url where state=$1 or state=$2")
	rows, err := u.Pool.Query(query, Created, Processing)
	defer func(logger *zap.Logger) {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				logger.Error("could not close row", zap.Error(err))
			}
		}

	}(u.Logger)

	if err != nil {
		return 0, err
	}
	var total int64
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			u.Logger.Error("Could not populate url struct ", zap.Error(err))
			return 0, err
		}
	}
	return total, nil
}

func (u *UrlRepository) countTerminated() (int64, error) {
	u.Logger.Info("Getting url to analyze")

	query := fmt.Sprintf("select count(*) from url where state=$1")
	rows, err := u.Pool.Query(query, Processed)
	defer func(logger *zap.Logger) {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				logger.Error("could not close row", zap.Error(err))
			}
		}

	}(u.Logger)

	if err != nil {
		return 0, err
	}
	var total int64
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			u.Logger.Error("Could not populate url struct ", zap.Error(err))
			return 0, err
		}
	}
	return total, nil
}
