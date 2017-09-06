package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindKeywordsShouldBeEmptyWhenNoKeywords(t *testing.T) {
	clearDB()

	keywords, err := FindKeywords(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(keywords))
}

func TestFindKeywordsShouldFindAKeyword(t *testing.T) {
	clearDB()

	keyword, _, err := FindOrCreateKeywordByName(db, "solo")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)

	keywords, err := FindKeywords(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(keywords))
	assert.Equal(t, "solo", keywords[0].Name)
}

func TestFindKeywordsShouldSortKeywordsByName(t *testing.T) {
	clearDB()

	_, _, err := FindOrCreateKeywordByName(db, "delta")
	assert.Nil(t, err)
	_, _, err = FindOrCreateKeywordByName(db, "baker")
	assert.Nil(t, err)
	_, _, err = FindOrCreateKeywordByName(db, "apple")
	assert.Nil(t, err)
	_, _, err = FindOrCreateKeywordByName(db, "charlie")
	assert.Nil(t, err)

	keywords, err := FindKeywords(db)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(keywords))
	assert.Equal(t, "apple", keywords[0].Name)
	assert.Equal(t, "baker", keywords[1].Name)
	assert.Equal(t, "charlie", keywords[2].Name)
	assert.Equal(t, "delta", keywords[3].Name)
}

func TestFindOrCreateKeywordByNameShouldCreateKeyword(t *testing.T) {
	clearDB()

	keyword, created, err := FindOrCreateKeywordByName(db, "my-keyword")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)
	assert.True(t, created)
	assert.Equal(t, "my-keyword", keyword.Name)

	var check Keyword
	db.Where("name = ?", "my-keyword").First(&check)
	assert.Equal(t, "my-keyword", check.Name)
}

func TestFindOrCreateKeywordShouldNotCreateDuplicateNames(t *testing.T) {
	clearDB()

	keyword, created, err := FindOrCreateKeywordByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)
	assert.True(t, created)
	assert.Equal(t, "foo", keyword.Name)

	keyword, created, err = FindOrCreateKeywordByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)
	assert.False(t, created)
	assert.Equal(t, "foo", keyword.Name)

	var keywords []Keyword
	db.Where("name = ?", "foo").Find(&keywords)
	assert.Equal(t, 1, len(keywords))
}

func TestFindKeywordByNameShouldReturnNilIfNotExists(t *testing.T) {
	clearDB()

	keyword, err := FindKeywordByName(db, "this does not exist")
	assert.Nil(t, err)
	assert.Nil(t, keyword)
}

func TestFindKeywordByNameShouldFindKeyword(t *testing.T) {
	clearDB()

	keyword, created, err := FindOrCreateKeywordByName(db, "creating a new keyword")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)
	assert.True(t, created)
	assert.Equal(t, "creating a new keyword", keyword.Name)

	newKeyword, err := FindKeywordByName(db, "creating a new keyword")
	assert.Nil(t, err)
	assert.Equal(t, "creating a new keyword", newKeyword.Name)
}

func TestRenameKeywordShouldRenameKeyword(t *testing.T) {
	clearDB()

	keyword, created, err := FindOrCreateKeywordByName(db, "old name")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)
	assert.True(t, created)
	assert.Equal(t, "old name", keyword.Name)

	err = keyword.Rename(db, "new name")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)
	assert.Equal(t, "new name", keyword.Name)
}

func TestRenameKeywordToExistingNameShouldReturnError(t *testing.T) {
	clearDB()

	first, created, err := FindOrCreateKeywordByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.True(t, created)
	assert.Equal(t, "first", first.Name)

	second, created, err := FindOrCreateKeywordByName(db, "second")
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

func TestRenameKeywordByChangingCaseShouldRenameKeyword(t *testing.T) {
	clearDB()

	first, _, err := FindOrCreateKeywordByName(db, "first")
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

func TestRenameKeywordByChangingToSameNameShouldReturnError(t *testing.T) {
	clearDB()

	same, _, err := FindOrCreateKeywordByName(db, "same")
	assert.Nil(t, err)
	assert.NotNil(t, same)
	assert.Equal(t, "same", same.Name)

	err = same.Rename(db, "same")
	assert.NotNil(t, err)
}

func TestDeleteKeywordShouldDeleteKeyword(t *testing.T) {
	clearDB()

	keyword, created, err := FindOrCreateKeywordByName(db, "to delete")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)
	assert.True(t, created)
	assert.Equal(t, "to delete", keyword.Name)

	err = keyword.Delete(db)
	assert.Nil(t, err)

	deleted, err := FindKeywordByName(db, "to delete")
	assert.Nil(t, err)
	assert.Nil(t, deleted)
}

func TestDeleteKeywordShouldDeleteAssociationsToStars(t *testing.T) {
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

	keyword, _, err := FindOrCreateKeywordByName(db, "jaguars")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)
	assert.Equal(t, "jaguars", keyword.Name)

	err = star1.AddKeyword(db, keyword)
	assert.Nil(t, err)

	err = star2.AddKeyword(db, keyword)
	assert.Nil(t, err)

	err = star1.LoadKeywords(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star1.Keywords))
	assert.Equal(t, "jaguars", star1.Keywords[0].Name)

	err = star2.LoadKeywords(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star2.Keywords))
	assert.Equal(t, "jaguars", star2.Keywords[0].Name)

	err = keyword.Delete(db)
	assert.Nil(t, err)

	err = star1.LoadKeywords(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star1.Keywords))

	err = star2.LoadKeywords(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star2.Keywords))
}

func TestLoadStarsShouldReturnErrorWhenKeywordNotInDatabase(t *testing.T) {
	clearDB()

	keyword := &Keyword{
		Name: "not in db",
	}

	err := keyword.LoadStars(db, "")
	assert.NotNil(t, err)
	assert.Equal(t, "Keyword '0' not found", err.Error())
}

func TestLoadStarsShouldLoadNoStarsWhenKeywordHasNoStars(t *testing.T) {
	clearDB()

	keyword, _, err := FindOrCreateKeywordByName(db, "keyword")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)

	err = keyword.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(keyword.Stars))
}

func TestLoadStarsShouldFillStars(t *testing.T) {
	clearDB()

	keyword, _, err := FindOrCreateKeywordByName(db, "keyword")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddKeyword(db, keyword)
	assert.Nil(t, err)

	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddKeyword(db, keyword)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(keyword.Stars))
	err = keyword.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(keyword.Stars))
}

func TestLoadStarsShouldFillStarsWithMatch(t *testing.T) {
	clearDB()

	keyword, _, err := FindOrCreateKeywordByName(db, "keyword")
	assert.Nil(t, err)
	assert.NotNil(t, keyword)

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
	err = star1.AddKeyword(db, keyword)
	assert.Nil(t, err)

	name2 := "Jacksonville Suns"
	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
		FullName:  &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddKeyword(db, keyword)
	assert.Nil(t, err)

	name3 := "Florida Gators"
	star3 := &Star{
		RemoteID:  "3",
		ServiceID: service.ID,
		FullName:  &name3,
	}
	_, err = CreateOrUpdateStar(db, star3, service)
	assert.Nil(t, err)
	err = star3.AddKeyword(db, keyword)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(keyword.Stars))
	err = keyword.LoadStars(db, "jacksonville")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(keyword.Stars))
	assert.Equal(t, "Jacksonville Jaguars", *keyword.Stars[0].FullName)
	assert.Equal(t, "Jacksonville Suns", *keyword.Stars[1].FullName)
}
