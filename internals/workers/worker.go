package worker

import (
	"log"

	"github.com/geoffowuor/url-shortener/internals/models"
	"gorm.io/gorm"
)

type ClickEventQueue struct {
	Events chan models.ClickEvent
}

func NewClickEventQueue(size int) *ClickEventQueue {
	return &ClickEventQueue{
		Events: make(chan models.ClickEvent, size),
	}
}

func StartClickWorker(db *gorm.DB, queue *ClickEventQueue) {
	go func() {
		for event := range queue.Events {
			log.Printf("WORKER RECEIVED: %+v", event)

			if err := db.Create(&event).Error; err != nil {
				log.Printf("WORKER DB ERROR: %v", err)
				continue
			}

			log.Println("WORKER SAVED EVENT")
		}
	}()
}
