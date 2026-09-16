package worker

import (
	"github.com/geoffowuor/url-shortener/internals/models"
)

type ClickEventQueue struct {
	Events chan models.ClickEvent
}

func NewClickEventQueue(size int) *ClickEventQueue {
	return &ClickEventQueue{
		Events: make(chan models.ClickEvent, size),
	}
}
