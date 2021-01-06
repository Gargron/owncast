// This is a centralized place to connect to the database, and hold a reference to it.
// Other packages can share this reference.  This package would also be a place to add any kind of
// persistence-related convenience methods or migrations.

package data

import (
	"database/sql"
	"os"

	"github.com/bwmarrin/snowflake"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	schemaVersion = 0
)

var (
	_db   *gorm.DB
	_node *snowflake.Node
)

func GetDatabase() *gorm.DB {
	return _db
}

func GetNode() *snowflake.Node {
	return _node
}

func SetupPersistence() error {
	file := config.Config.DatabaseFilePath

	// Create empty DB file if it doesn't exist.
	if !utils.DoesFileExists(file) {
		log.Traceln("Creating new database at", file)

		_, err := os.Create(file)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})

	if err != nil {
		return err
	}

	_db = db

	node, err := snowflake.NewNode(1)

	if err != nil {
		return err
	}

	_node = node

	return nil
}

func migrateDatabase(db *sql.DB, from, to int) error {
	return nil
}
