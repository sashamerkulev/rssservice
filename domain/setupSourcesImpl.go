package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/models"
)

type SetupSourcesRepository interface {
	GetSources() ([]models.Link, error)
}

type SetupSources struct {
	Logger logger.Logger
}

func (ss SetupSources) GetSources() ([]models.Link, error) {
	return models.Urls, nil
}
