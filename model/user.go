package model

type User struct {
	UserId int64
	Name   string
	Phone  string
	Token  string
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
