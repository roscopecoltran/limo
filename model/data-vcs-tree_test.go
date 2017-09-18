package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindTreesShouldBeEmptyWhenNoTrees(t *testing.T) {
	clearDB()

	trees, err := FindTrees(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(trees))
}

func TestFindTreesShouldFindATree(t *testing.T) {
	clearDB()

	tree, _, err := FindOrCreateTreeByName(db, "solo")
	assert.Nil(t, err)
	assert.NotNil(t, tree)

	trees, err := FindTrees(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(trees))
	assert.Equal(t, "solo", trees[0].Name)
}

func TestFindTreesShouldSortTreesByName(t *testing.T) {
	clearDB()

	_, _, err := FindOrCreateTreeByName(db, "delta")
	assert.Nil(t, err)
	_, _, err = FindOrCreateTreeByName(db, "baker")
	assert.Nil(t, err)
	_, _, err = FindOrCreateTreeByName(db, "apple")
	assert.Nil(t, err)
	_, _, err = FindOrCreateTreeByName(db, "charlie")
	assert.Nil(t, err)

	trees, err := FindTrees(db)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(trees))
	assert.Equal(t, "apple", trees[0].Name)
	assert.Equal(t, "baker", trees[1].Name)
	assert.Equal(t, "charlie", trees[2].Name)
	assert.Equal(t, "delta", trees[3].Name)
}

func TestFindOrCreateTreeByNameShouldCreateTree(t *testing.T) {
	clearDB()

	tree, created, err := FindOrCreateTreeByName(db, "my-tree")
	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.True(t, created)
	assert.Equal(t, "my-tree", tree.Name)

	var check Tree
	db.Where("name = ?", "my-tree").First(&check)
	assert.Equal(t, "my-tree", check.Name)
}

func TestFindOrCreateTreeShouldNotCreateDuplicateNames(t *testing.T) {
	clearDB()

	tree, created, err := FindOrCreateTreeByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.True(t, created)
	assert.Equal(t, "foo", tree.Name)

	tree, created, err = FindOrCreateTreeByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.False(t, created)
	assert.Equal(t, "foo", tree.Name)

	var trees []Tree
	db.Where("name = ?", "foo").Find(&trees)
	assert.Equal(t, 1, len(trees))
}

func TestFindTreeByNameShouldReturnNilIfNotExists(t *testing.T) {
	clearDB()

	tree, err := FindTreeByName(db, "this does not exist")
	assert.Nil(t, err)
	assert.Nil(t, tree)
}

func TestFindTreeByNameShouldFindTree(t *testing.T) {
	clearDB()

	tree, created, err := FindOrCreateTreeByName(db, "creating a new tree")
	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.True(t, created)
	assert.Equal(t, "creating a new tree", tree.Name)

	newTree, err := FindTreeByName(db, "creating a new tree")
	assert.Nil(t, err)
	assert.Equal(t, "creating a new tree", newTree.Name)
}

func TestRenameTreeShouldRenameTree(t *testing.T) {
	clearDB()

	tree, created, err := FindOrCreateTreeByName(db, "old name")
	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.True(t, created)
	assert.Equal(t, "old name", tree.Name)

	err = tree.Rename(db, "new name")
	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, "new name", tree.Name)
}

func TestRenameTreeToExistingNameShouldReturnError(t *testing.T) {
	clearDB()

	first, created, err := FindOrCreateTreeByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.True(t, created)
	assert.Equal(t, "first", first.Name)

	second, created, err := FindOrCreateTreeByName(db, "second")
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

func TestRenameTreeByChangingCaseShouldRenameTree(t *testing.T) {
	clearDB()

	first, _, err := FindOrCreateTreeByName(db, "first")
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

func TestRenameTreeByChangingToSameNameShouldReturnError(t *testing.T) {
	clearDB()

	same, _, err := FindOrCreateTreeByName(db, "same")
	assert.Nil(t, err)
	assert.NotNil(t, same)
	assert.Equal(t, "same", same.Name)

	err = same.Rename(db, "same")
	assert.NotNil(t, err)
}

func TestDeleteTreeShouldDeleteTree(t *testing.T) {
	clearDB()

	tree, created, err := FindOrCreateTreeByName(db, "to delete")
	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.True(t, created)
	assert.Equal(t, "to delete", tree.Name)

	err = tree.Delete(db)
	assert.Nil(t, err)

	deleted, err := FindTreeByName(db, "to delete")
	assert.Nil(t, err)
	assert.Nil(t, deleted)
}

func TestDeleteTreeShouldDeleteAssociationsToStars(t *testing.T) {
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

	tree, _, err := FindOrCreateTreeByName(db, "jaguars")
	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, "jaguars", tree.Name)

	err = star1.AddTree(db, tree)
	assert.Nil(t, err)

	err = star2.AddTree(db, tree)
	assert.Nil(t, err)

	err = star1.LoadTrees(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star1.Trees))
	assert.Equal(t, "jaguars", star1.Trees[0].Name)

	err = star2.LoadTrees(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star2.Trees))
	assert.Equal(t, "jaguars", star2.Trees[0].Name)

	err = tree.Delete(db)
	assert.Nil(t, err)

	err = star1.LoadTrees(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star1.Trees))

	err = star2.LoadTrees(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star2.Trees))
}

func TestLoadStarsShouldReturnErrorWhenTreeNotInDatabase(t *testing.T) {
	clearDB()

	tree := &Tree{
		Name: "not in db",
	}

	err := tree.LoadStars(db, "")
	assert.NotNil(t, err)
	assert.Equal(t, "Tree '0' not found", err.Error())
}

func TestLoadStarsShouldLoadNoStarsWhenTreeHasNoStars(t *testing.T) {
	clearDB()

	tree, _, err := FindOrCreateTreeByName(db, "tree")
	assert.Nil(t, err)
	assert.NotNil(t, tree)

	err = tree.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tree.Stars))
}

func TestLoadStarsShouldFillStars(t *testing.T) {
	clearDB()

	tree, _, err := FindOrCreateTreeByName(db, "tree")
	assert.Nil(t, err)
	assert.NotNil(t, tree)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddTree(db, tree)
	assert.Nil(t, err)

	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddTree(db, tree)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(tree.Stars))
	err = tree.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tree.Stars))
}

func TestLoadStarsShouldFillStarsWithMatch(t *testing.T) {
	clearDB()

	tree, _, err := FindOrCreateTreeByName(db, "tree")
	assert.Nil(t, err)
	assert.NotNil(t, tree)

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
	err = star1.AddTree(db, tree)
	assert.Nil(t, err)

	name2 := "Jacksonville Suns"
	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
		FullName:  &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddTree(db, tree)
	assert.Nil(t, err)

	name3 := "Florida Gators"
	star3 := &Star{
		RemoteID:  "3",
		ServiceID: service.ID,
		FullName:  &name3,
	}
	_, err = CreateOrUpdateStar(db, star3, service)
	assert.Nil(t, err)
	err = star3.AddTree(db, tree)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tree.Stars))
	err = tree.LoadStars(db, "jacksonville")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tree.Stars))
	assert.Equal(t, "Jacksonville Jaguars", *tree.Stars[0].FullName)
	assert.Equal(t, "Jacksonville Suns", *tree.Stars[1].FullName)
}
