package model

type User struct {
	UserId int64
	UserToken string
	Name   string
	Phone  string
}

type ArticleActivity struct {
	Likes    int64
	Dislikes int64
}

type ArticleUserActivity struct {
	Like    bool
	Dislike bool
}

type ArticleUser struct {
	Article
	ArticleActivity
	ArticleUserActivity
}
