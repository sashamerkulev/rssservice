package fcm

import (
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

func SendNotificationLikeArticleComment(userName string, commentId int64, dislike bool, logger logger.Logger) {
	opt := option.WithCredentialsFile("F:/go/src/github.com/sashamerkulev/rssservice/fcm/serviceAccountKey.json")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Log("ERROR", "SENDNOTIFICATIONNEWARTICLESMESSAGE", err.Error())
		return
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		logger.Log("ERROR", "SENDNOTIFICATIONNEWARTICLESMESSAGE", err.Error())
		return
	}
	firebaseToken, err := mysql.GetFirebaseIdByArticleCommentId(commentId, logger)
	if err != nil {
		return
	}
	message := &messaging.Message{
		Data: map[string]string{
			"userName":      userName,
			"commentId":     fmt.Sprint(commentId),
			"likeOrDislike": fmt.Sprint(dislike),
		},
		Token: firebaseToken,
	}
	response, err := client.Send(ctx, message)
	if err != nil {
		logger.Log("ERROR", "SENDNOTIFICATIONNEWARTICLESMESSAGE", err.Error())
		return
	}
	logger.Log("DEBUG", "SENDNOTIFICATIONNEWARTICLESMESSAGE", response)
}
