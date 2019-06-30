package data

import (
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

func GetUserIdByToken(userToken string) (userId int64, err error) {
	if len(userToken) == 0 {
		return -1, nil
	}
	rows, err := DB.Query("select userId from userDeviceToken WHERE token = ?", userToken)
	if err != nil {
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

func FindUserIdByDeviceId(deviceId string, logger logger.Logger) (userId int64, err error) {
	rows, err := DB.Query("select userId from userDeviceToken WHERE deviceId = ?", deviceId)
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

func AddTokenForUserIdAndDeviceId(userId int64, deviceId string, userToken string, logger logger.Logger) error {
	_, err := DB.Exec("insert into userDeviceToken(userId, deviceId, timestamp, token) values(?,?,?,?)", userId, deviceId, time.Now(), userToken)
	if err != nil {
		logger.Log("ERROR", "ADDTOKENFORUSERIDANDDEVICEID", err.Error())
		return errors.UserRegistrationError
	}
	return nil
}

func RegisterUser(deviceId string, firebaseId string, userToken string, logger logger.Logger) (userId int64, err error) {
	tx, err := DB.Begin()
	if err != nil {
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.TransactionOpenError
	}
	defer tx.Commit()
	result, err := DB.Exec("insert into userInfo(UserName, UserPhone) values(?,?)", "", "")
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
	result, err = DB.Exec("insert into userDevices(userId, deviceId, firebaseId) values(?,?,?)", newUserId, deviceId, firebaseId)
	if err != nil {
		tx.Rollback()
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.UserRegistrationError
	}
	err = AddTokenForUserIdAndDeviceId(newUserId, deviceId, userToken, logger)
	return newUserId, err
}

func UpdateFirebaseId(userId int64, firebaseId string, logger logger.Logger) error {
	_, err := DB.Exec("update userDevices set firebaseId = ? where userId = ?", firebaseId, userId)
	if err != nil {
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return errors.UserRegistrationError
	}
	return nil
}

func UpdateUser(userId int64, name string, phone string, logger logger.Logger) (err error) {
	_, err = DB.Exec("update userInfo set UserName=?, UserPhone=? where userId = ?", name, phone, userId)
	if err != nil {
		logger.Log("ERROR", "UPDATEUSER", err.Error())
		return errors.UserUpdateError
	}
	return nil
}

func UploadUserPhoto(userId int64, photo []byte, logger logger.Logger) (err error) {
	_, err = DB.Exec("update userInfo set UserPhoto=? where userId = ?", photo, userId)
	if err != nil {
		logger.Log("ERROR", "UPLOADUSERPHOTO", err.Error())
		return errors.UserPhotoUploadError
	}
	return nil
}

func GetUserInfo(userId int64, logger logger.Logger) (user model.User, err error) {
	rows, err := DB.Query("select UserId, UserName, UserPhone from userInfo where userId = ?", userId)
	if err != nil {
		logger.Log("ERROR", "GETUSERINFO", err.Error())
		return model.User{Name: ""}, nil
	}
	if rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.UserId, &user.Name, &user.Phone)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		return user, nil
	}
	return model.User{Name: ""}, errors.UserNotFoundError
}
