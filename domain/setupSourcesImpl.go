package domain

import "github.com/sashamerkulev/rssservice/logger"

type SetupSources struct {
	Logger     logger.Logger
	Repository MainRepository
}
