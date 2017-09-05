package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSoftwaresShouldBeEmptyWhenNoSoftwares(t *testing.T) {
	clearDB()

	softwares, err := FindSoftwares(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(softwares))
}

func TestFindSoftwaresShouldFindASoftware(t *testing.T) {
	clearDB()

	software, _, err := FindOrCreateSoftwareByName(db, "solo")
	assert.Nil(t, err)
	assert.NotNil(t, software)

	softwares, err := FindSoftwares(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(softwares))
	assert.Equal(t, "solo", softwares[0].Name)
}

func TestFindSoftwaresShouldSortSoftwaresByName(t *testing.T) {
	clearDB()

	_, _, err := FindOrCreateSoftwareByName(db, "delta")
	assert.Nil(t, err)
	_, _, err = FindOrCreateSoftwareByName(db, "baker")
	assert.Nil(t, err)
	_, _, err = FindOrCreateSoftwareByName(db, "apple")
	assert.Nil(t, err)
	_, _, err = FindOrCreateSoftwareByName(db, "charlie")
	assert.Nil(t, err)

	softwares, err := FindSoftwares(db)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(softwares))
	assert.Equal(t, "apple", softwares[0].Name)
	assert.Equal(t, "baker", softwares[1].Name)
	assert.Equal(t, "charlie", softwares[2].Name)
	assert.Equal(t, "delta", softwares[3].Name)
}

func TestFindOrCreateSoftwareByNameShouldCreateSoftware(t *testing.T) {
	clearDB()

	software, created, err := FindOrCreateSoftwareByName(db, "my-software")
	assert.Nil(t, err)
	assert.NotNil(t, software)
	assert.True(t, created)
	assert.Equal(t, "my-software", software.Name)

	var check Software
	db.Where("name = ?", "my-software").First(&check)
	assert.Equal(t, "my-software", check.Name)
}

func TestFindOrCreateSoftwareShouldNotCreateDuplicateNames(t *testing.T) {
	clearDB()

	software, created, err := FindOrCreateSoftwareByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, software)
	assert.True(t, created)
	assert.Equal(t, "foo", software.Name)

	software, created, err = FindOrCreateSoftwareByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, software)
	assert.False(t, created)
	assert.Equal(t, "foo", software.Name)

	var softwares []Software
	db.Where("name = ?", "foo").Find(&softwares)
	assert.Equal(t, 1, len(softwares))
}

func TestFindSoftwareByNameShouldReturnNilIfNotExists(t *testing.T) {
	clearDB()

	software, err := FindSoftwareByName(db, "this does not exist")
	assert.Nil(t, err)
	assert.Nil(t, software)
}

func TestFindSoftwareByNameShouldFindSoftware(t *testing.T) {
	clearDB()

	software, created, err := FindOrCreateSoftwareByName(db, "creating a new software")
	assert.Nil(t, err)
	assert.NotNil(t, software)
	assert.True(t, created)
	assert.Equal(t, "creating a new software", software.Name)

	newSoftware, err := FindSoftwareByName(db, "creating a new software")
	assert.Nil(t, err)
	assert.Equal(t, "creating a new software", newSoftware.Name)
}

func TestRenameSoftwareShouldRenameSoftware(t *testing.T) {
	clearDB()

	software, created, err := FindOrCreateSoftwareByName(db, "old name")
	assert.Nil(t, err)
	assert.NotNil(t, software)
	assert.True(t, created)
	assert.Equal(t, "old name", software.Name)

	err = software.Rename(db, "new name")
	assert.Nil(t, err)
	assert.NotNil(t, software)
	assert.Equal(t, "new name", software.Name)
}

func TestRenameSoftwareToExistingNameShouldReturnError(t *testing.T) {
	clearDB()

	first, created, err := FindOrCreateSoftwareByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.True(t, created)
	assert.Equal(t, "first", first.Name)

	second, created, err := FindOrCreateSoftwareByName(db, "second")
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

func TestRenameSoftwareByChangingCaseShouldRenameSoftware(t *testing.T) {
	clearDB()

	first, _, err := FindOrCreateSoftwareByName(db, "first")
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

func TestRenameSoftwareByChangingToSameNameShouldReturnError(t *testing.T) {
	clearDB()

	same, _, err := FindOrCreateSoftwareByName(db, "same")
	assert.Nil(t, err)
	assert.NotNil(t, same)
	assert.Equal(t, "same", same.Name)

	err = same.Rename(db, "same")
	assert.NotNil(t, err)
}

func TestDeleteSoftwareShouldDeleteSoftware(t *testing.T) {
	clearDB()

	software, created, err := FindOrCreateSoftwareByName(db, "to delete")
	assert.Nil(t, err)
	assert.NotNil(t, software)
	assert.True(t, created)
	assert.Equal(t, "to delete", software.Name)

	err = software.Delete(db)
	assert.Nil(t, err)

	deleted, err := FindSoftwareByName(db, "to delete")
	assert.Nil(t, err)
	assert.Nil(t, deleted)
}

func TestDeleteSoftwareShouldDeleteAssociationsToStars(t *testing.T) {
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

	software, _, err := FindOrCreateSoftwareByName(db, "jaguars")
	assert.Nil(t, err)
	assert.NotNil(t, software)
	assert.Equal(t, "jaguars", software.Name)

	err = star1.AddSoftware(db, software)
	assert.Nil(t, err)

	err = star2.AddSoftware(db, software)
	assert.Nil(t, err)

	err = star1.LoadSoftwares(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star1.Softwares))
	assert.Equal(t, "jaguars", star1.Softwares[0].Name)

	err = star2.LoadSoftwares(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star2.Softwares))
	assert.Equal(t, "jaguars", star2.Softwares[0].Name)

	err = software.Delete(db)
	assert.Nil(t, err)

	err = star1.LoadSoftwares(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star1.Softwares))

	err = star2.LoadSoftwares(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star2.Softwares))
}

func TestLoadStarsShouldReturnErrorWhenSoftwareNotInDatabase(t *testing.T) {
	clearDB()

	software := &Software{
		Name: "not in db",
	}

	err := software.LoadStars(db, "")
	assert.NotNil(t, err)
	assert.Equal(t, "Software '0' not found", err.Error())
}

func TestLoadStarsShouldLoadNoStarsWhenSoftwareHasNoStars(t *testing.T) {
	clearDB()

	software, _, err := FindOrCreateSoftwareByName(db, "software")
	assert.Nil(t, err)
	assert.NotNil(t, software)

	err = software.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(software.Stars))
}

func TestLoadStarsShouldFillStars(t *testing.T) {
	clearDB()

	software, _, err := FindOrCreateSoftwareByName(db, "software")
	assert.Nil(t, err)
	assert.NotNil(t, software)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddSoftware(db, software)
	assert.Nil(t, err)

	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddSoftware(db, software)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(software.Stars))
	err = software.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(software.Stars))
}

func TestLoadStarsShouldFillStarsWithMatch(t *testing.T) {
	clearDB()

	software, _, err := FindOrCreateSoftwareByName(db, "software")
	assert.Nil(t, err)
	assert.NotNil(t, software)

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
	err = star1.AddSoftware(db, software)
	assert.Nil(t, err)

	name2 := "Jacksonville Suns"
	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
		FullName:  &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddSoftware(db, software)
	assert.Nil(t, err)

	name3 := "Florida Gators"
	star3 := &Star{
		RemoteID:  "3",
		ServiceID: service.ID,
		FullName:  &name3,
	}
	_, err = CreateOrUpdateStar(db, star3, service)
	assert.Nil(t, err)
	err = star3.AddSoftware(db, software)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(software.Stars))
	err = software.LoadStars(db, "jacksonville")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(software.Stars))
	assert.Equal(t, "Jacksonville Jaguars", *software.Stars[0].FullName)
	assert.Equal(t, "Jacksonville Suns", *software.Stars[1].FullName)
}
