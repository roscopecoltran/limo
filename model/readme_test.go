package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindReadmesShouldBeEmptyWhenNoReadmes(t *testing.T) {
	clearDB()

	readmes, err := FindReadmes(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(readmes))
}

func TestFindReadmesShouldFindAReadme(t *testing.T) {
	clearDB()

	readme, _, err := FindOrCreateReadmeByName(db, "solo")
	assert.Nil(t, err)
	assert.NotNil(t, readme)

	readmes, err := FindReadmes(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(readmes))
	assert.Equal(t, "solo", readmes[0].Name)
}

func TestFindReadmesShouldSortReadmesByName(t *testing.T) {
	clearDB()

	_, _, err := FindOrCreateReadmeByName(db, "delta")
	assert.Nil(t, err)
	_, _, err = FindOrCreateReadmeByName(db, "baker")
	assert.Nil(t, err)
	_, _, err = FindOrCreateReadmeByName(db, "apple")
	assert.Nil(t, err)
	_, _, err = FindOrCreateReadmeByName(db, "charlie")
	assert.Nil(t, err)

	readmes, err := FindReadmes(db)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(readmes))
	assert.Equal(t, "apple", readmes[0].Name)
	assert.Equal(t, "baker", readmes[1].Name)
	assert.Equal(t, "charlie", readmes[2].Name)
	assert.Equal(t, "delta", readmes[3].Name)
}

func TestFindOrCreateReadmeByNameShouldCreateReadme(t *testing.T) {
	clearDB()

	readme, created, err := FindOrCreateReadmeByName(db, "my-readme")
	assert.Nil(t, err)
	assert.NotNil(t, readme)
	assert.True(t, created)
	assert.Equal(t, "my-readme", readme.Name)

	var check Readme
	db.Where("name = ?", "my-readme").First(&check)
	assert.Equal(t, "my-readme", check.Name)
}

func TestFindOrCreateReadmeShouldNotCreateDuplicateNames(t *testing.T) {
	clearDB()

	readme, created, err := FindOrCreateReadmeByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, readme)
	assert.True(t, created)
	assert.Equal(t, "foo", readme.Name)

	readme, created, err = FindOrCreateReadmeByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, readme)
	assert.False(t, created)
	assert.Equal(t, "foo", readme.Name)

	var readmes []Readme
	db.Where("name = ?", "foo").Find(&readmes)
	assert.Equal(t, 1, len(readmes))
}

func TestFindReadmeByNameShouldReturnNilIfNotExists(t *testing.T) {
	clearDB()

	readme, err := FindReadmeByName(db, "this does not exist")
	assert.Nil(t, err)
	assert.Nil(t, readme)
}

func TestFindReadmeByNameShouldFindReadme(t *testing.T) {
	clearDB()

	readme, created, err := FindOrCreateReadmeByName(db, "creating a new readme")
	assert.Nil(t, err)
	assert.NotNil(t, readme)
	assert.True(t, created)
	assert.Equal(t, "creating a new readme", readme.Name)

	newReadme, err := FindReadmeByName(db, "creating a new readme")
	assert.Nil(t, err)
	assert.Equal(t, "creating a new readme", newReadme.Name)
}

func TestRenameReadmeShouldRenameReadme(t *testing.T) {
	clearDB()

	readme, created, err := FindOrCreateReadmeByName(db, "old name")
	assert.Nil(t, err)
	assert.NotNil(t, readme)
	assert.True(t, created)
	assert.Equal(t, "old name", readme.Name)

	err = readme.Rename(db, "new name")
	assert.Nil(t, err)
	assert.NotNil(t, readme)
	assert.Equal(t, "new name", readme.Name)
}

func TestRenameReadmeToExistingNameShouldReturnError(t *testing.T) {
	clearDB()

	first, created, err := FindOrCreateReadmeByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.True(t, created)
	assert.Equal(t, "first", first.Name)

	second, created, err := FindOrCreateReadmeByName(db, "second")
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

func TestRenameReadmeByChangingCaseShouldRenameReadme(t *testing.T) {
	clearDB()

	first, _, err := FindOrCreateReadmeByName(db, "first")
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

func TestRenameReadmeByChangingToSameNameShouldReturnError(t *testing.T) {
	clearDB()

	same, _, err := FindOrCreateReadmeByName(db, "same")
	assert.Nil(t, err)
	assert.NotNil(t, same)
	assert.Equal(t, "same", same.Name)

	err = same.Rename(db, "same")
	assert.NotNil(t, err)
}

func TestDeleteReadmeShouldDeleteReadme(t *testing.T) {
	clearDB()

	readme, created, err := FindOrCreateReadmeByName(db, "to delete")
	assert.Nil(t, err)
	assert.NotNil(t, readme)
	assert.True(t, created)
	assert.Equal(t, "to delete", readme.Name)

	err = readme.Delete(db)
	assert.Nil(t, err)

	deleted, err := FindReadmeByName(db, "to delete")
	assert.Nil(t, err)
	assert.Nil(t, deleted)
}

func TestDeleteReadmeShouldDeleteAssociationsToStars(t *testing.T) {
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

	readme, _, err := FindOrCreateReadmeByName(db, "jaguars")
	assert.Nil(t, err)
	assert.NotNil(t, readme)
	assert.Equal(t, "jaguars", readme.Name)

	err = star1.AddReadme(db, readme)
	assert.Nil(t, err)

	err = star2.AddReadme(db, readme)
	assert.Nil(t, err)

	err = star1.LoadReadmes(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star1.Readmes))
	assert.Equal(t, "jaguars", star1.Readmes[0].Name)

	err = star2.LoadReadmes(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star2.Readmes))
	assert.Equal(t, "jaguars", star2.Readmes[0].Name)

	err = readme.Delete(db)
	assert.Nil(t, err)

	err = star1.LoadReadmes(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star1.Readmes))

	err = star2.LoadReadmes(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star2.Readmes))
}

func TestLoadStarsShouldReturnErrorWhenReadmeNotInDatabase(t *testing.T) {
	clearDB()

	readme := &Readme{
		Name: "not in db",
	}

	err := readme.LoadStars(db, "")
	assert.NotNil(t, err)
	assert.Equal(t, "Readme '0' not found", err.Error())
}

func TestLoadStarsShouldLoadNoStarsWhenReadmeHasNoStars(t *testing.T) {
	clearDB()

	readme, _, err := FindOrCreateReadmeByName(db, "readme")
	assert.Nil(t, err)
	assert.NotNil(t, readme)

	err = readme.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(readme.Stars))
}

func TestLoadStarsShouldFillStars(t *testing.T) {
	clearDB()

	readme, _, err := FindOrCreateReadmeByName(db, "readme")
	assert.Nil(t, err)
	assert.NotNil(t, readme)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddReadme(db, readme)
	assert.Nil(t, err)

	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddReadme(db, readme)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(readme.Stars))
	err = readme.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(readme.Stars))
}

func TestLoadStarsShouldFillStarsWithMatch(t *testing.T) {
	clearDB()

	readme, _, err := FindOrCreateReadmeByName(db, "readme")
	assert.Nil(t, err)
	assert.NotNil(t, readme)

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
	err = star1.AddReadme(db, readme)
	assert.Nil(t, err)

	name2 := "Jacksonville Suns"
	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
		FullName:  &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddReadme(db, readme)
	assert.Nil(t, err)

	name3 := "Florida Gators"
	star3 := &Star{
		RemoteID:  "3",
		ServiceID: service.ID,
		FullName:  &name3,
	}
	_, err = CreateOrUpdateStar(db, star3, service)
	assert.Nil(t, err)
	err = star3.AddReadme(db, readme)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(readme.Stars))
	err = readme.LoadStars(db, "jacksonville")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(readme.Stars))
	assert.Equal(t, "Jacksonville Jaguars", *readme.Stars[0].FullName)
	assert.Equal(t, "Jacksonville Suns", *readme.Stars[1].FullName)
}
