package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/reader"
)

type SetupSourcesRepository interface {
	GetSources() ([]reader.Link, error)
}

type SetupSources struct {
	Logger     logger.Logger
}

func (ss SetupSources) GetSources() ([]reader.Link, error) {
	return reader.Urls, nil
}
