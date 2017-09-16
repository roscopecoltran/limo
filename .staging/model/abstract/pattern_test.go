package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPatternsShouldBeEmptyWhenNoPatterns(t *testing.T) {
	clearDB()

	patterns, err := FindPatterns(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(patterns))
}

func TestFindPatternsShouldFindAPattern(t *testing.T) {
	clearDB()

	pattern, _, err := FindOrCreatePatternByName(db, "solo")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)

	patterns, err := FindPatterns(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(patterns))
	assert.Equal(t, "solo", patterns[0].Name)
}

func TestFindPatternsShouldSortPatternsByName(t *testing.T) {
	clearDB()

	_, _, err := FindOrCreatePatternByName(db, "delta")
	assert.Nil(t, err)
	_, _, err = FindOrCreatePatternByName(db, "baker")
	assert.Nil(t, err)
	_, _, err = FindOrCreatePatternByName(db, "apple")
	assert.Nil(t, err)
	_, _, err = FindOrCreatePatternByName(db, "charlie")
	assert.Nil(t, err)

	patterns, err := FindPatterns(db)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(patterns))
	assert.Equal(t, "apple", patterns[0].Name)
	assert.Equal(t, "baker", patterns[1].Name)
	assert.Equal(t, "charlie", patterns[2].Name)
	assert.Equal(t, "delta", patterns[3].Name)
}

func TestFindOrCreatePatternByNameShouldCreatePattern(t *testing.T) {
	clearDB()

	pattern, created, err := FindOrCreatePatternByName(db, "my-pattern")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, created)
	assert.Equal(t, "my-pattern", pattern.Name)

	var check Pattern
	db.Where("name = ?", "my-pattern").First(&check)
	assert.Equal(t, "my-pattern", check.Name)
}

func TestFindOrCreatePatternShouldNotCreateDuplicateNames(t *testing.T) {
	clearDB()

	pattern, created, err := FindOrCreatePatternByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, created)
	assert.Equal(t, "foo", pattern.Name)

	pattern, created, err = FindOrCreatePatternByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.False(t, created)
	assert.Equal(t, "foo", pattern.Name)

	var patterns []Pattern
	db.Where("name = ?", "foo").Find(&patterns)
	assert.Equal(t, 1, len(patterns))
}

func TestFindPatternByNameShouldReturnNilIfNotExists(t *testing.T) {
	clearDB()

	pattern, err := FindPatternByName(db, "this does not exist")
	assert.Nil(t, err)
	assert.Nil(t, pattern)
}

func TestFindPatternByNameShouldFindPattern(t *testing.T) {
	clearDB()

	pattern, created, err := FindOrCreatePatternByName(db, "creating a new pattern")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, created)
	assert.Equal(t, "creating a new pattern", pattern.Name)

	newPattern, err := FindPatternByName(db, "creating a new pattern")
	assert.Nil(t, err)
	assert.Equal(t, "creating a new pattern", newPattern.Name)
}

func TestRenamePatternShouldRenamePattern(t *testing.T) {
	clearDB()

	pattern, created, err := FindOrCreatePatternByName(db, "old name")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, created)
	assert.Equal(t, "old name", pattern.Name)

	err = pattern.Rename(db, "new name")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.Equal(t, "new name", pattern.Name)
}

func TestRenamePatternToExistingNameShouldReturnError(t *testing.T) {
	clearDB()

	first, created, err := FindOrCreatePatternByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.True(t, created)
	assert.Equal(t, "first", first.Name)

	second, created, err := FindOrCreatePatternByName(db, "second")
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

func TestRenamePatternByChangingCaseShouldRenamePattern(t *testing.T) {
	clearDB()

	first, _, err := FindOrCreatePatternByName(db, "first")
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

func TestRenamePatternByChangingToSameNameShouldReturnError(t *testing.T) {
	clearDB()

	same, _, err := FindOrCreatePatternByName(db, "same")
	assert.Nil(t, err)
	assert.NotNil(t, same)
	assert.Equal(t, "same", same.Name)

	err = same.Rename(db, "same")
	assert.NotNil(t, err)
}

func TestDeletePatternShouldDeletePattern(t *testing.T) {
	clearDB()

	pattern, created, err := FindOrCreatePatternByName(db, "to delete")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, created)
	assert.Equal(t, "to delete", pattern.Name)

	err = pattern.Delete(db)
	assert.Nil(t, err)

	deleted, err := FindPatternByName(db, "to delete")
	assert.Nil(t, err)
	assert.Nil(t, deleted)
}

func TestDeletePatternShouldDeleteAssociationsToStars(t *testing.T) {
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

	pattern, _, err := FindOrCreatePatternByName(db, "jaguars")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.Equal(t, "jaguars", pattern.Name)

	err = star1.AddPattern(db, pattern)
	assert.Nil(t, err)

	err = star2.AddPattern(db, pattern)
	assert.Nil(t, err)

	err = star1.LoadPatterns(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star1.Patterns))
	assert.Equal(t, "jaguars", star1.Patterns[0].Name)

	err = star2.LoadPatterns(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star2.Patterns))
	assert.Equal(t, "jaguars", star2.Patterns[0].Name)

	err = pattern.Delete(db)
	assert.Nil(t, err)

	err = star1.LoadPatterns(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star1.Patterns))

	err = star2.LoadPatterns(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star2.Patterns))
}

func TestLoadStarsShouldReturnErrorWhenPatternNotInDatabase(t *testing.T) {
	clearDB()

	pattern := &Pattern{
		Name: "not in db",
	}

	err := pattern.LoadStars(db, "")
	assert.NotNil(t, err)
	assert.Equal(t, "Pattern '0' not found", err.Error())
}

func TestLoadStarsShouldLoadNoStarsWhenPatternHasNoStars(t *testing.T) {
	clearDB()

	pattern, _, err := FindOrCreatePatternByName(db, "pattern")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)

	err = pattern.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(pattern.Stars))
}

func TestLoadStarsShouldFillStars(t *testing.T) {
	clearDB()

	pattern, _, err := FindOrCreatePatternByName(db, "pattern")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddPattern(db, pattern)
	assert.Nil(t, err)

	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddPattern(db, pattern)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(pattern.Stars))
	err = pattern.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(pattern.Stars))
}

func TestLoadStarsShouldFillStarsWithMatch(t *testing.T) {
	clearDB()

	pattern, _, err := FindOrCreatePatternByName(db, "pattern")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)

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
	err = star1.AddPattern(db, pattern)
	assert.Nil(t, err)

	name2 := "Jacksonville Suns"
	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
		FullName:  &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddPattern(db, pattern)
	assert.Nil(t, err)

	name3 := "Florida Gators"
	star3 := &Star{
		RemoteID:  "3",
		ServiceID: service.ID,
		FullName:  &name3,
	}
	_, err = CreateOrUpdateStar(db, star3, service)
	assert.Nil(t, err)
	err = star3.AddPattern(db, pattern)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(pattern.Stars))
	err = pattern.LoadStars(db, "jacksonville")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(pattern.Stars))
	assert.Equal(t, "Jacksonville Jaguars", *pattern.Stars[0].FullName)
	assert.Equal(t, "Jacksonville Suns", *pattern.Stars[1].FullName)
}
