package model

import (
	"errors"
	"fmt"
	"strings"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/ckaznocha/taggraph"
)

// https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/bot/repos.go

// Topic represents a topic in the database
type Topic struct {
	gorm.Model
	Name      		string `gorm:"type:varchar(64);not null;unique"`
	Count 			int    `gorm:"-"`
	Stars     		[]Star `gorm:"many2many:star_topics;"`
}

type TopicResult struct {
	Topic  	*Topic
	Error 	error
}

// should provide a map[string]map[string]
func TestTopicsGraph(query string) (taggraph.Tagger, error) {
	tagg.AddChildToTag("shirts", "clothes")
	tagg.AddChildToTag("pants", "clothes")
	tagg.AddChildToTag("dress clothes", "clothes")
	tagg.AddChildToTag("shirts", "dress clothes")
	tagg.AddChildToTag("shirts", "tops")
	tagg.AddChildToTag("tops", "casual")
	tagg.AddChildToTag("casual", "clothes")
	entities, ok := tagg.GetTag(query)
	if !ok {
		log.WithFields(logrus.Fields{"action": "PrintTopicsGraph", "step": "GetTag", "model": "Topic", "query": query}).Warnf("Tag topic not found", query)
		return nil, fmt.Errorf("Tag topic not found", query)
	}
	// iterate
	paths := entities.PathsToAllAncestorsAsString("->")
	for _, path := range paths {
		log.WithFields(logrus.Fields{"action": "PrintTopicsGraph", "step": "PathsToAllAncestorsAsString", "model": "Topic", "query": query, "path": path}).Info("New path discovered.")
	}
	return entities, nil

}

// FindTopics finds all topics
func FindTopics(db *gorm.DB) ([]Topic, error) {
	var topics []Topic
	db.Order("name").Find(&topics)
	return topics, db.Error
}

// FindTopicsWithStarCount finds all topics and gets their count of stars
func FindTopicsWithStarCount(db *gorm.DB) ([]Topic, error) {
	var topics []Topic
	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.TOPIC_ID) AS STARCOUNT
		FROM TOPICS T
		LEFT JOIN STAR_TOPICS ST ON T.ID = ST.TOPIC_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()

	if err != nil {
		return topics, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		var topic Topic
		if err = rows.Scan(&topic.Name); err != nil {
			return topics, err
		}
		topics = append(topics, topic)
	}
	return topics, db.Error
}

// FindTopicByName finds a topic by name
func FindTopicByName(db *gorm.DB, name string) (*Topic, error) {
	var topic Topic
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&topic).RecordNotFound() {
		return nil, db.Error
	}
	return &topic, db.Error
}

// FindOrCreateTopicByName finds a topic by name, creating if it doesn't exist
func FindOrCreateTopicByName(db *gorm.DB, name string) (*Topic, bool, error) {
	var topic Topic
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&topic).RecordNotFound() {
		topic.Name = name
		err := db.Create(&topic).Error
		return &topic, true, err
	}
	return &topic, false, nil
}

// LoadStars loads the stars for a topic
func (topic *Topic) LoadStars(db *gorm.DB, match string) error {
	// Make sure topic exists in database, or we will panic
	var existing Topic
	if db.Where("id = ?", topic.ID).First(&existing).RecordNotFound() {
		err := fmt.Errorf("Topic '%d' not found", topic.ID)
		log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "LoadStars"}).Errorf("Topic '%d' not found", topic.ID)
		return err
	}
	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM STARS S
			INNER JOIN STAR_TOPICS ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.TOPIC_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			topic.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		topic.Stars = stars
		return db.Error
	}
	return db.Model(topic).Association("Stars").Find(&topic.Stars).Error
}

// Rename renames a topic -- new name must not already exist
func (topic *Topic) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == topic.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(topic.Name) {
		existing, err := FindTopicByName(db, name)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "Rename"}).Warnf("%#s", err)
			return err
		}
		if existing != nil {
			err := fmt.Errorf("Topic '%s' already exists", existing.Name)
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "Rename"}).Errorf("Topic '%s' already exists", existing.Name)
			return err
		}
	}

	topic.Name = name
	return db.Save(topic).Error
}

// Delete deletes a topic and disassociates it from any stars
func (topic *Topic) Delete(db *gorm.DB) error {
	if err := db.Model(topic).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(topic).Error
}
