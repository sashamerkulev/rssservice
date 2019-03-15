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
	UserId  int64
	Like    bool
	Dislike bool
}

type ArticleUser struct {
	Article
	ArticleActivity
	ArticleUserActivity
}

//type IRegisterUser interface {
//	RegisterUser() (userToken string, err error)
//}
//
//type IUpdateUser interface {
//	UpdateUserFunc() (err error)
//}
//
//type IArticleUser interface {
//	GetArticleUser(userToken string, lastTime time.Time, logger logger.Logger) ([]ArticleUser, error)
//}
