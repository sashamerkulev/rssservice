package mysql

import (
	"database/sql"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"time"
)

type SetupRegisterRepositoryImpl struct {
	DB *sql.DB
}

func (db SetupRegisterRepositoryImpl) FindUserIdByDeviceId(deviceId string, logger logger.Logger) (userId int64, err error) {
	rows, err := db.DB.Query("select userId from userDeviceToken WHERE deviceId = ?", deviceId)
	if err != nil {
		logger.Log("ERROR", "FINDUSERIDBYDEVICEID", err.Error())
		return -1, errors.UserNotFoundError
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return -1, errors.UserNotFoundError
		}
		return userId, nil
	}
	return -1, errors.UserNotFoundError
}

func (db SetupRegisterRepositoryImpl) AddTokenForUserIdAndDeviceId(userId int64, deviceId string, token string, logger logger.Logger) error {
	_, err := db.DB.Exec("insert into userDeviceToken(userId, deviceId, timestamp, token) values(?,?,?,?)", userId, deviceId, time.Now(), token)
	if err != nil {
		logger.Log("ERROR", "ADDTOKENFORUSERIDANDDEVICEID", err.Error())
		return errors.UserRegistrationError
	}
	return nil
}

func (db SetupRegisterRepositoryImpl) RegisterUser(deviceId string, firebaseId string, token string, logger logger.Logger) (userId int64, err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.TransactionOpenError
	}
	defer tx.Commit()
	result, err := db.DB.Exec("insert into userInfo(UserName, UserPhone) values(?,?)", "", "")
	if err != nil {
		tx.Rollback()
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.UserRegistrationError
	}
	newUserId, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.UserRegistrationError
	}
	result, err = db.DB.Exec("insert into userDevices(userId, deviceId, firebaseId) values(?,?,?)", newUserId, deviceId, firebaseId)
	if err != nil {
		tx.Rollback()
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.UserRegistrationError
	}
	err = addTokenForUserIdAndDeviceId(db.DB, newUserId, deviceId, token, logger)
	return newUserId, err
}

func (db SetupRegisterRepositoryImpl) UpdateFirebaseId(userId int64, setupId string, firebaseId string, logger logger.Logger) error {
	_, err := db.DB.Exec("update userDevices set firebaseId = ? where userId = ?", firebaseId, userId)
	if err != nil {
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return errors.UserRegistrationError
	}
	return nil
}

func addTokenForUserIdAndDeviceId(DB *sql.DB, userId int64, deviceId string, userToken string, logger logger.Logger) error {
	_, err := DB.Exec("insert into userDeviceToken(userId, deviceId, timestamp, token) values(?,?,?,?)", userId, deviceId, time.Now(), userToken)
	if err != nil {
		logger.Log("ERROR", "ADDTOKENFORUSERIDANDDEVICEID", err.Error())
		return errors.UserRegistrationError
	}
	return nil
}
