package mysql

import (
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
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

func GetFirebaseIdByCommentId(commentId int64, logger logger.Logger) (string, error) {
	rows, err := DB.Query("select ud.FirebaseId from usercommentlikes ucl	join userdevices ud on ucl.userId = ud.userId " +
		" where ucl.commentId = ?", commentId)
	if err != nil {
		logger.Log("ERROR", "GETUSERINFO", err.Error())
		return "", err
	}
	if rows.Next() {
		var firebaseId string
		err := rows.Scan(&firebaseId)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
			return "", errors.CommentNotFoundError
		}
		return firebaseId, nil
	}
	return "", errors.CommentNotFoundError
}