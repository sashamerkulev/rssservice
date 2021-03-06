package errors

var (
	UserNotFoundError            = Error("User not found")
	TransactionOpenError         = Error("Transaction is failed")
	UserUpdateError              = Error("User updating is failed")
	UserPhotoUploadError         = Error("User photo updating is failed")
	UserRegistrationError        = Error("User registration is failed")
	ArticleNotFoundError         = Error("Article not found")
	CommentNotFoundError         = Error("Comment not found")
	EmptyParametersError         = Error("Empty  articleId or commentId parameters")
	WrongArticleIdParameterError = Error("Wrong articleId parameters")
	WrongCommentIdParameterError = Error("Wrong commentId parameters")
)

func Error(text string) error {
	return &errorString{text}
}

type errorString struct {
	message string
}

func (e *errorString) Error() string {
	return e.message
}
