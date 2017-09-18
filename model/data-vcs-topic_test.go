package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindTopicsShouldBeEmptyWhenNoTopics(t *testing.T) {
	clearDB()

	topics, err := FindTopics(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(topics))
}

func TestFindTopicsShouldFindATopic(t *testing.T) {
	clearDB()

	topic, _, err := FindOrCreateTopicByName(db, "solo")
	assert.Nil(t, err)
	assert.NotNil(t, topic)

	topics, err := FindTopics(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(topics))
	assert.Equal(t, "solo", topics[0].Name)
}

func TestFindTopicsShouldSortTopicsByName(t *testing.T) {
	clearDB()

	_, _, err := FindOrCreateTopicByName(db, "delta")
	assert.Nil(t, err)
	_, _, err = FindOrCreateTopicByName(db, "baker")
	assert.Nil(t, err)
	_, _, err = FindOrCreateTopicByName(db, "apple")
	assert.Nil(t, err)
	_, _, err = FindOrCreateTopicByName(db, "charlie")
	assert.Nil(t, err)

	topics, err := FindTopics(db)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(topics))
	assert.Equal(t, "apple", topics[0].Name)
	assert.Equal(t, "baker", topics[1].Name)
	assert.Equal(t, "charlie", topics[2].Name)
	assert.Equal(t, "delta", topics[3].Name)
}

func TestFindOrCreateTopicByNameShouldCreateTopic(t *testing.T) {
	clearDB()

	topic, created, err := FindOrCreateTopicByName(db, "my-topic")
	assert.Nil(t, err)
	assert.NotNil(t, topic)
	assert.True(t, created)
	assert.Equal(t, "my-topic", topic.Name)

	var check Topic
	db.Where("name = ?", "my-topic").First(&check)
	assert.Equal(t, "my-topic", check.Name)
}

func TestFindOrCreateTopicShouldNotCreateDuplicateNames(t *testing.T) {
	clearDB()

	topic, created, err := FindOrCreateTopicByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, topic)
	assert.True(t, created)
	assert.Equal(t, "foo", topic.Name)

	topic, created, err = FindOrCreateTopicByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, topic)
	assert.False(t, created)
	assert.Equal(t, "foo", topic.Name)

	var topics []Topic
	db.Where("name = ?", "foo").Find(&topics)
	assert.Equal(t, 1, len(topics))
}

func TestFindTopicByNameShouldReturnNilIfNotExists(t *testing.T) {
	clearDB()

	topic, err := FindTopicByName(db, "this does not exist")
	assert.Nil(t, err)
	assert.Nil(t, topic)
}

func TestFindTopicByNameShouldFindTopic(t *testing.T) {
	clearDB()

	topic, created, err := FindOrCreateTopicByName(db, "creating a new topic")
	assert.Nil(t, err)
	assert.NotNil(t, topic)
	assert.True(t, created)
	assert.Equal(t, "creating a new topic", topic.Name)

	newTopic, err := FindTopicByName(db, "creating a new topic")
	assert.Nil(t, err)
	assert.Equal(t, "creating a new topic", newTopic.Name)
}

func TestRenameTopicShouldRenameTopic(t *testing.T) {
	clearDB()

	topic, created, err := FindOrCreateTopicByName(db, "old name")
	assert.Nil(t, err)
	assert.NotNil(t, topic)
	assert.True(t, created)
	assert.Equal(t, "old name", topic.Name)

	err = topic.Rename(db, "new name")
	assert.Nil(t, err)
	assert.NotNil(t, topic)
	assert.Equal(t, "new name", topic.Name)
}

func TestRenameTopicToExistingNameShouldReturnError(t *testing.T) {
	clearDB()

	first, created, err := FindOrCreateTopicByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.True(t, created)
	assert.Equal(t, "first", first.Name)

	second, created, err := FindOrCreateTopicByName(db, "second")
	assert.Nil(t, err)
	assert.NotNil(t, second)
	assert.True(t, created)
	assert.Equal(t, "second", second.Name)

	err = second.Rename(db, "first")
	assert.NotNil(t, err)
	assert.Equal(t, "second", second.Name)

	err = second.Rename(db, "First")
	assert.NotNil(t, err)
	assert.Equal(t, "second", second.Name)

	err = second.Rename(db, "FIRST")
	assert.NotNil(t, err)
	assert.Equal(t, "second", second.Name)
}

func TestRenameTopicByChangingCaseShouldRenameTopic(t *testing.T) {
	clearDB()

	first, _, err := FindOrCreateTopicByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.Equal(t, "first", first.Name)

	err = first.Rename(db, "First")
	assert.Nil(t, err)
	assert.Equal(t, "First", first.Name)

	err = first.Rename(db, "FIRST")
	assert.Nil(t, err)
	assert.Equal(t, "FIRST", first.Name)
}

func TestRenameTopicByChangingToSameNameShouldReturnError(t *testing.T) {
	clearDB()

	same, _, err := FindOrCreateTopicByName(db, "same")
	assert.Nil(t, err)
	assert.NotNil(t, same)
	assert.Equal(t, "same", same.Name)

	err = same.Rename(db, "same")
	assert.NotNil(t, err)
}

func TestDeleteTopicShouldDeleteTopic(t *testing.T) {
	clearDB()

	topic, created, err := FindOrCreateTopicByName(db, "to delete")
	assert.Nil(t, err)
	assert.NotNil(t, topic)
	assert.True(t, created)
	assert.Equal(t, "to delete", topic.Name)

	err = topic.Delete(db)
	assert.Nil(t, err)

	deleted, err := FindTopicByName(db, "to delete")
	assert.Nil(t, err)
	assert.Nil(t, deleted)
}

func TestDeleteTopicShouldDeleteAssociationsToStars(t *testing.T) {
	clearDB()

	service, _, err := FindOrCreateServiceByName(db, "nfl")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nfl", service.Name)

	name1 := "Allen Hurns"
	star1 := &Star{
		RemoteID: "88",
		Name:     &name1,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)

	name2 := "Allen Robinson"
	star2 := &Star{
		RemoteID: "15",
		Name:     &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)

	topic, _, err := FindOrCreateTopicByName(db, "jaguars")
	assert.Nil(t, err)
	assert.NotNil(t, topic)
	assert.Equal(t, "jaguars", topic.Name)

	err = star1.AddTopic(db, topic)
	assert.Nil(t, err)

	err = star2.AddTopic(db, topic)
	assert.Nil(t, err)

	err = star1.LoadTopics(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star1.Topics))
	assert.Equal(t, "jaguars", star1.Topics[0].Name)

	err = star2.LoadTopics(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star2.Topics))
	assert.Equal(t, "jaguars", star2.Topics[0].Name)

	err = topic.Delete(db)
	assert.Nil(t, err)

	err = star1.LoadTopics(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star1.Topics))

	err = star2.LoadTopics(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star2.Topics))
}

func TestLoadStarsShouldReturnErrorWhenTopicNotInDatabase(t *testing.T) {
	clearDB()

	topic := &Topic{
		Name: "not in db",
	}

	err := topic.LoadStars(db, "")
	assert.NotNil(t, err)
	assert.Equal(t, "Topic '0' not found", err.Error())
}

func TestLoadStarsShouldLoadNoStarsWhenTopicHasNoStars(t *testing.T) {
	clearDB()

	topic, _, err := FindOrCreateTopicByName(db, "topic")
	assert.Nil(t, err)
	assert.NotNil(t, topic)

	err = topic.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(topic.Stars))
}

func TestLoadStarsShouldFillStars(t *testing.T) {
	clearDB()

	topic, _, err := FindOrCreateTopicByName(db, "topic")
	assert.Nil(t, err)
	assert.NotNil(t, topic)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddTopic(db, topic)
	assert.Nil(t, err)

	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddTopic(db, topic)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(topic.Stars))
	err = topic.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(topic.Stars))
}

func TestLoadStarsShouldFillStarsWithMatch(t *testing.T) {
	clearDB()

	topic, _, err := FindOrCreateTopicByName(db, "topic")
	assert.Nil(t, err)
	assert.NotNil(t, topic)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	name1 := "Jacksonville Jaguars"
	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
		FullName:  &name1,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddTopic(db, topic)
	assert.Nil(t, err)

	name2 := "Jacksonville Suns"
	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
		FullName:  &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddTopic(db, topic)
	assert.Nil(t, err)

	name3 := "Florida Gators"
	star3 := &Star{
		RemoteID:  "3",
		ServiceID: service.ID,
		FullName:  &name3,
	}
	_, err = CreateOrUpdateStar(db, star3, service)
	assert.Nil(t, err)
	err = star3.AddTopic(db, topic)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(topic.Stars))
	err = topic.LoadStars(db, "jacksonville")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(topic.Stars))
	assert.Equal(t, "Jacksonville Jaguars", *topic.Stars[0].FullName)
	assert.Equal(t, "Jacksonville Suns", *topic.Stars[1].FullName)
}
