package model

import "fmt"

type User struct {
	UserId int64
	Name   string
	Phone  string
	Token  string
}

type ArticleActivity struct {
	Likes    int64
	Dislikes int64
	Comments int64
}

type ArticleUserActivity struct {
	Like    int
	Dislike int
	Comment int
}

type ArticleUser struct {
	Article
	ArticleActivity
	ArticleUserActivity
}

func MakeSureNameExists(userInfo User) User {
	if len(userInfo.Name) == 0 {
		userInfo.Name = "Гость_" + fmt.Sprint(userInfo.UserId)
	}
	return userInfo
}

func MakeSureCommenterExists(userArticleComment UserArticleComment) UserArticleComment {
	if len(userArticleComment.Name) == 0 {
		userArticleComment.Name = "Гость_" + fmt.Sprint(userArticleComment.UserId)
	}
	return userArticleComment
}
