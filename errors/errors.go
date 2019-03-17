package errors

func UserNotFoundError() error {
	return Error("User not found")
}

func TransactionOpenError() error {
	return Error("Transaction is failed")
}

func UserUpdateError() error {
	return Error("User updating is failed")
}

func UserRegistrationError() error {
	return Error("User registration is failed")
}

func ArticleNotFoundError() error {
	return Error("Article not found")
}

func Error(text string) error {
	return &errorString{text}
}

type errorString struct {
	message string
}

func (e *errorString) Error() string {
	return e.message
}
