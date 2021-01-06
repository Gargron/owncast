package chat

import (
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _db *gorm.DB

func setupPersistence() {
	_db = data.GetDatabase()
	_db.AutoMigrate(&models.ChatEvent{})
}

func addMessage(message models.ChatEvent) {
	result := _db.Create(&message)

	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func getChatHistory(filtered bool) []models.ChatEvent {
	history := make([]models.ChatEvent, 0)
	query := _db.Where("message_type != ?", "SYSTEM").Where("datetime(created_at) >= datetime('now', '-1 day')")

	if filtered {
		query = query.Where("deleted_at IS NULL")
	}

	result := query.Find(&history)

	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return history
}

func saveMessageVisibility(messageIDs []string, visible bool) error {
	var deletedAt gorm.DeletedAt

	if visible {
		deletedAt = gorm.DeletedAt{}
	} else {
		deletedAt = gorm.DeletedAt{Time: time.Now()}
	}

	result := _db.Model(models.ChatEvent{}).Where("id IN (?)", messageIDs).Updates(models.ChatEvent{DeletedAt: deletedAt})

	if result.Error != nil {
		log.Fatal(result.Error)
		return result.Error
	}

	return nil
}

func getMessageById(messageID string) (models.ChatEvent, error) {
	var chatEvent models.ChatEvent

	result := _db.First(&chatEvent, messageID)

	if result.Error != nil {
		log.Fatal(result.Error)
		return models.ChatEvent{}, result.Error
	}

	return chatEvent, nil
}
