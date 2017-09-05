package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAcademicsShouldBeEmptyWhenNoAcademics(t *testing.T) {
	clearDB()

	academics, err := FindAcademics(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(academics))
}

func TestFindAcademicsShouldFindAAcademic(t *testing.T) {
	clearDB()

	academic, _, err := FindOrCreateAcademicByName(db, "solo")
	assert.Nil(t, err)
	assert.NotNil(t, academic)

	academics, err := FindAcademics(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(academics))
	assert.Equal(t, "solo", academics[0].Name)
}

func TestFindAcademicsShouldSortAcademicsByName(t *testing.T) {
	clearDB()

	_, _, err := FindOrCreateAcademicByName(db, "delta")
	assert.Nil(t, err)
	_, _, err = FindOrCreateAcademicByName(db, "baker")
	assert.Nil(t, err)
	_, _, err = FindOrCreateAcademicByName(db, "apple")
	assert.Nil(t, err)
	_, _, err = FindOrCreateAcademicByName(db, "charlie")
	assert.Nil(t, err)

	academics, err := FindAcademics(db)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(academics))
	assert.Equal(t, "apple", academics[0].Name)
	assert.Equal(t, "baker", academics[1].Name)
	assert.Equal(t, "charlie", academics[2].Name)
	assert.Equal(t, "delta", academics[3].Name)
}

func TestFindOrCreateAcademicByNameShouldCreateAcademic(t *testing.T) {
	clearDB()

	academic, created, err := FindOrCreateAcademicByName(db, "my-academic")
	assert.Nil(t, err)
	assert.NotNil(t, academic)
	assert.True(t, created)
	assert.Equal(t, "my-academic", academic.Name)

	var check Academic
	db.Where("name = ?", "my-academic").First(&check)
	assert.Equal(t, "my-academic", check.Name)
}

func TestFindOrCreateAcademicShouldNotCreateDuplicateNames(t *testing.T) {
	clearDB()

	academic, created, err := FindOrCreateAcademicByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, academic)
	assert.True(t, created)
	assert.Equal(t, "foo", academic.Name)

	academic, created, err = FindOrCreateAcademicByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, academic)
	assert.False(t, created)
	assert.Equal(t, "foo", academic.Name)

	var academics []Academic
	db.Where("name = ?", "foo").Find(&academics)
	assert.Equal(t, 1, len(academics))
}

func TestFindAcademicByNameShouldReturnNilIfNotExists(t *testing.T) {
	clearDB()

	academic, err := FindAcademicByName(db, "this does not exist")
	assert.Nil(t, err)
	assert.Nil(t, academic)
}

func TestFindAcademicByNameShouldFindAcademic(t *testing.T) {
	clearDB()

	academic, created, err := FindOrCreateAcademicByName(db, "creating a new academic")
	assert.Nil(t, err)
	assert.NotNil(t, academic)
	assert.True(t, created)
	assert.Equal(t, "creating a new academic", academic.Name)

	newAcademic, err := FindAcademicByName(db, "creating a new academic")
	assert.Nil(t, err)
	assert.Equal(t, "creating a new academic", newAcademic.Name)
}

func TestRenameAcademicShouldRenameAcademic(t *testing.T) {
	clearDB()

	academic, created, err := FindOrCreateAcademicByName(db, "old name")
	assert.Nil(t, err)
	assert.NotNil(t, academic)
	assert.True(t, created)
	assert.Equal(t, "old name", academic.Name)

	err = academic.Rename(db, "new name")
	assert.Nil(t, err)
	assert.NotNil(t, academic)
	assert.Equal(t, "new name", academic.Name)
}

func TestRenameAcademicToExistingNameShouldReturnError(t *testing.T) {
	clearDB()

	first, created, err := FindOrCreateAcademicByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.True(t, created)
	assert.Equal(t, "first", first.Name)

	second, created, err := FindOrCreateAcademicByName(db, "second")
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

func TestRenameAcademicByChangingCaseShouldRenameAcademic(t *testing.T) {
	clearDB()

	first, _, err := FindOrCreateAcademicByName(db, "first")
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

func TestRenameAcademicByChangingToSameNameShouldReturnError(t *testing.T) {
	clearDB()

	same, _, err := FindOrCreateAcademicByName(db, "same")
	assert.Nil(t, err)
	assert.NotNil(t, same)
	assert.Equal(t, "same", same.Name)

	err = same.Rename(db, "same")
	assert.NotNil(t, err)
}

func TestDeleteAcademicShouldDeleteAcademic(t *testing.T) {
	clearDB()

	academic, created, err := FindOrCreateAcademicByName(db, "to delete")
	assert.Nil(t, err)
	assert.NotNil(t, academic)
	assert.True(t, created)
	assert.Equal(t, "to delete", academic.Name)

	err = academic.Delete(db)
	assert.Nil(t, err)

	deleted, err := FindAcademicByName(db, "to delete")
	assert.Nil(t, err)
	assert.Nil(t, deleted)
}

func TestDeleteAcademicShouldDeleteAssociationsToStars(t *testing.T) {
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

	academic, _, err := FindOrCreateAcademicByName(db, "jaguars")
	assert.Nil(t, err)
	assert.NotNil(t, academic)
	assert.Equal(t, "jaguars", academic.Name)

	err = star1.AddAcademic(db, academic)
	assert.Nil(t, err)

	err = star2.AddAcademic(db, academic)
	assert.Nil(t, err)

	err = star1.LoadAcademics(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star1.Academics))
	assert.Equal(t, "jaguars", star1.Academics[0].Name)

	err = star2.LoadAcademics(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star2.Academics))
	assert.Equal(t, "jaguars", star2.Academics[0].Name)

	err = academic.Delete(db)
	assert.Nil(t, err)

	err = star1.LoadAcademics(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star1.Academics))

	err = star2.LoadAcademics(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star2.Academics))
}

func TestLoadStarsShouldReturnErrorWhenAcademicNotInDatabase(t *testing.T) {
	clearDB()

	academic := &Academic{
		Name: "not in db",
	}

	err := academic.LoadStars(db, "")
	assert.NotNil(t, err)
	assert.Equal(t, "Academic '0' not found", err.Error())
}

func TestLoadStarsShouldLoadNoStarsWhenAcademicHasNoStars(t *testing.T) {
	clearDB()

	academic, _, err := FindOrCreateAcademicByName(db, "academic")
	assert.Nil(t, err)
	assert.NotNil(t, academic)

	err = academic.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(academic.Stars))
}

func TestLoadStarsShouldFillStars(t *testing.T) {
	clearDB()

	academic, _, err := FindOrCreateAcademicByName(db, "academic")
	assert.Nil(t, err)
	assert.NotNil(t, academic)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddAcademic(db, academic)
	assert.Nil(t, err)

	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddAcademic(db, academic)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(academic.Stars))
	err = academic.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(academic.Stars))
}

func TestLoadStarsShouldFillStarsWithMatch(t *testing.T) {
	clearDB()

	academic, _, err := FindOrCreateAcademicByName(db, "academic")
	assert.Nil(t, err)
	assert.NotNil(t, academic)

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
	err = star1.AddAcademic(db, academic)
	assert.Nil(t, err)

	name2 := "Jacksonville Suns"
	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
		FullName:  &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddAcademic(db, academic)
	assert.Nil(t, err)

	name3 := "Florida Gators"
	star3 := &Star{
		RemoteID:  "3",
		ServiceID: service.ID,
		FullName:  &name3,
	}
	_, err = CreateOrUpdateStar(db, star3, service)
	assert.Nil(t, err)
	err = star3.AddAcademic(db, academic)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(academic.Stars))
	err = academic.LoadStars(db, "jacksonville")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(academic.Stars))
	assert.Equal(t, "Jacksonville Jaguars", *academic.Stars[0].FullName)
	assert.Equal(t, "Jacksonville Suns", *academic.Stars[1].FullName)
}

