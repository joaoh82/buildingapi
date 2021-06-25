package presenters

import (
	"github.com/joaoh82/buildingapi/interfaces"
	"github.com/rs/zerolog"
)

type presenters struct {
	logger zerolog.Logger
}

func NewPresenters(logger zerolog.Logger) interfaces.Presenters {
	return &presenters{logger: logger}
}
