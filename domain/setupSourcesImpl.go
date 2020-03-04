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
	Sources []models.Link
}

func (ss SetupSources) GetSources() ([]models.Link, error) {
	return ss.Sources, nil
}
