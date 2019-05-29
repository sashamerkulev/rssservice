package model

import "time"

type ArticleComment struct {
	CommentId int64
	ArticleId int64
	UserId    int64
	Name      string
	Comment   string
	PubDate   time.Time
	Status    int
}

type UserArticleComment struct {
	ArticleComment
	Likes    int64
	Dislikes int64
	Like     int
	Dislike  int
	Owner    int
}
