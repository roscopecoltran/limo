package model

import (
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func TestNewRepoFromGithubShouldCopyFields(t *testing.T) {
	clearDB()

	id := 33
	name := "larry-bird"
	fullName := "celtics/larry-bird"
	description := "larry legend"
	homepage := "http://www.nba.com/celtics/"
	url := "http://www.nba.com/pacers/"
	language := "hoosier"
	repogazersCount := 10000

	timestamp := github.Timestamp{
		Time: time.Now(),
	}

	github := github.Repository{
		ID:              &id,
		Name:            &name,
		FullName:        &fullName,
		Description:     &description,
		Homepage:        &homepage,
		CloneURL:        &url,
		Language:        &language,
		RepogazersCount: &repogazersCount,
	}

	repo, err := NewRepoFromGithub(&timestamp, github)
	assert.Nil(t, err)
	assert.Equal(t, "33", repo.RemoteID)
	assert.Equal(t, name, *repo.Name)
	assert.Equal(t, fullName, *repo.FullName)
	assert.Equal(t, description, *repo.Description)
	assert.Equal(t, homepage, *repo.Homepage)
	assert.Equal(t, url, *repo.URL)
	assert.Equal(t, language, *repo.Language)
	assert.Equal(t, repogazersCount, repo.Repogazers)
}

func TestNewRepoFromGithubShouldHandleEmpty(t *testing.T) {
	clearDB()

	repo, err := NewRepoFromGithub(&github.Timestamp{}, github.Repository{})
	assert.NotNil(t, err)
	assert.Equal(t, "ID from GitHub is required", err.Error())
	assert.Nil(t, repo)
}

func TestNewRepoFromGithubShouldHandleOnlyID(t *testing.T) {
	clearDB()

	id := 33
	repo, err := NewRepoFromGithub(&github.Timestamp{}, github.Repository{
		ID: &id,
	})
	assert.Nil(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, "33", repo.RemoteID)
}

func TestNewRepoFromGitlabShouldCopyFields(t *testing.T) {
	clearDB()

	id := 33
	name := "larry-bird"
	fullName := "celtics/larry-bird"
	description := "larry legend"
	homepage := "http://www.nba.com/celtics/"
	url := "http://www.nba.com/pacers/"
	repogazersCount := 10000

	gitlab := gitlab.Project{
		ID:                id,
		Name:              name,
		NameWithNamespace: fullName,
		Description:       description,
		WebURL:            homepage,
		HTTPURLToRepo:     url,
		RepoCount:         repogazersCount,
	}

	repo, err := NewRepoFromGitlab(gitlab)
	assert.Nil(t, err)
	assert.Equal(t, "33", repo.RemoteID)
	assert.Equal(t, name, *repo.Name)
	assert.Equal(t, fullName, *repo.FullName)
	assert.Equal(t, description, *repo.Description)
	assert.Equal(t, homepage, *repo.Homepage)
	assert.Equal(t, url, *repo.URL)
	assert.Equal(t, (*string)(nil), repo.Language)
	assert.Equal(t, repogazersCount, repo.Repogazers)
}

func TestNewRepoFromGitlabShouldHandleOnlyID(t *testing.T) {
	clearDB()

	id := 33
	repo, err := NewRepoFromGitlab(gitlab.Project{
		ID: id,
	})
	assert.Nil(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, "33", repo.RemoteID)
}

func TestFuzzyFindReposByNameShouldFuzzyFind(t *testing.T) {
	clearDB()

	fullName := "Apple/Baker"
	name := "Charlie"

	repo := Repo{
		FullName: &fullName,
		Name:     &name,
	}
	assert.Nil(t, db.Create(&repo).Error)

	repos, err := FuzzyFindReposByName(db, "Apple/Baker")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repos))
	assert.Equal(t, fullName, *repos[0].FullName)
	assert.Equal(t, name, *repos[0].Name)

	repos, err = FuzzyFindReposByName(db, "Charlie")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repos))
	assert.Equal(t, fullName, *repos[0].FullName)
	assert.Equal(t, name, *repos[0].Name)

	repos, err = FuzzyFindReposByName(db, "apple/baker")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repos))
	assert.Equal(t, fullName, *repos[0].FullName)
	assert.Equal(t, name, *repos[0].Name)

	repos, err = FuzzyFindReposByName(db, "charlie")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repos))
	assert.Equal(t, fullName, *repos[0].FullName)
	assert.Equal(t, name, *repos[0].Name)

	repos, err = FuzzyFindReposByName(db, "apple")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repos))
	assert.Equal(t, fullName, *repos[0].FullName)
	assert.Equal(t, name, *repos[0].Name)

	repos, err = FuzzyFindReposByName(db, "harl")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repos))
	assert.Equal(t, fullName, *repos[0].FullName)
	assert.Equal(t, name, *repos[0].Name)

	repos, err = FuzzyFindReposByName(db, "boogers")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(repos))
}

func TestAddTagShouldAddTag(t *testing.T) {
	clearDB()

	tag, _, err := FindOrCreateTagByName(db, "celtics")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "celtics", tag.Name)

	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Isaiah Thomas" // Not a typo
	repo := &Repo{
		RemoteID: "remoteID",
		Name:     &name,
	}
	_, err = CreateOrUpdateRepo(db, repo, service)
	assert.Nil(t, err)

	err = repo.AddTag(db, tag)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repo.Tags))
	assert.Equal(t, "celtics", repo.Tags[0].Name)

	repos, err := FuzzyFindReposByName(db, "Isaiah Thomas")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repos))
	assert.Equal(t, "Isaiah Thomas", *repos[0].Name)

	err = repos[0].LoadTags(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repos[0].Tags))
	assert.Equal(t, "celtics", repos[0].Tags[0].Name)
}

func TestHasTagShouldReturnFalseWhenNoTags(t *testing.T) {
	clearDB()

	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	repo := &Repo{
		RemoteID: "brown",
		Name:     &name,
	}
	tag, _, err := FindOrCreateTagByName(db, "bucks")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "bucks", tag.Name)

	_, err = CreateOrUpdateRepo(db, repo, service)
	assert.Nil(t, err)

	err = repo.LoadTags(db)
	assert.Nil(t, err)

	assert.False(t, repo.HasTag(tag))
}

func TestHasTagShouldReturnFalseWhenTagIsNil(t *testing.T) {
	clearDB()

	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	repo := &Repo{
		RemoteID: "brown",
		Name:     &name,
	}
	_, err = CreateOrUpdateRepo(db, repo, service)
	assert.Nil(t, err)

	err = repo.LoadTags(db)
	assert.Nil(t, err)

	assert.False(t, repo.HasTag(nil))
}

func TestHasTagShouldReturnFalseWhenDoesNotHaveTag(t *testing.T) {
	clearDB()

	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	repo := &Repo{
		RemoteID: "brown",
		Name:     &name,
	}
	_, err = CreateOrUpdateRepo(db, repo, service)
	assert.Nil(t, err)

	bucks, _, err := FindOrCreateTagByName(db, "bucks")
	assert.Nil(t, err)
	assert.NotNil(t, bucks)
	assert.Equal(t, "bucks", bucks.Name)

	celtics, _, err := FindOrCreateTagByName(db, "celtics")
	assert.Nil(t, err)
	assert.NotNil(t, celtics)
	assert.Equal(t, "celtics", celtics.Name)

	err = repo.AddTag(db, celtics)
	assert.Nil(t, err)

	err = repo.LoadTags(db)
	assert.Nil(t, err)

	assert.False(t, repo.HasTag(bucks))
}

func TestHasTagShouldReturnTrueWhenHasOnlyTag(t *testing.T) {
	clearDB()

	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	repo := &Repo{
		RemoteID: "brown",
		Name:     &name,
	}
	_, err = CreateOrUpdateRepo(db, repo, service)
	assert.Nil(t, err)

	celtics, _, err := FindOrCreateTagByName(db, "celtics")
	assert.Nil(t, err)
	assert.NotNil(t, celtics)
	assert.Equal(t, "celtics", celtics.Name)

	err = repo.AddTag(db, celtics)
	assert.Nil(t, err)

	err = repo.LoadTags(db)
	assert.Nil(t, err)

	assert.True(t, repo.HasTag(celtics))
}

func TestHasTagShouldReturnTrueWhenHasTag(t *testing.T) {
	clearDB()

	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	repo := &Repo{
		RemoteID: "brown",
		Name:     &name,
	}
	_, err = CreateOrUpdateRepo(db, repo, service)
	assert.Nil(t, err)

	draft, _, err := FindOrCreateTagByName(db, "2016-draft")
	assert.Nil(t, err)
	assert.NotNil(t, draft)
	assert.Equal(t, "2016-draft", draft.Name)

	celtics, _, err := FindOrCreateTagByName(db, "celtics")
	assert.Nil(t, err)
	assert.NotNil(t, celtics)
	assert.Equal(t, "celtics", celtics.Name)

	err = repo.AddTag(db, celtics)
	assert.Nil(t, err)

	err = repo.LoadTags(db)
	assert.Nil(t, err)

	assert.True(t, repo.HasTag(celtics))
}

func TestLoadTagsShouldReturnErrorWhenRepoNotInDatabase(t *testing.T) {
	clearDB()

	name := "not in db"
	repo := &Repo{
		RemoteID: "not in db",
		Name:     &name,
	}

	err := repo.LoadTags(db)
	assert.NotNil(t, err)
	assert.Equal(t, "Repo '0' not found", err.Error())
}

func TestFindRepoByIDShouldReturnErrorWhenDoesNotExist(t *testing.T) {
	clearDB()

	repo, err := FindRepoByID(db, 1)
	assert.NotNil(t, err)
	assert.Equal(t, "Repo '1' not found", err.Error())
	assert.Nil(t, repo)
}

func TestFindRepoByIDShouldReturnRepo(t *testing.T) {
	clearDB()

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	repo := &Repo{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateRepo(db, repo, service)
	assert.Nil(t, err)

	existing, err := FindRepoByID(db, repo.ID)
	assert.Nil(t, err)
	assert.NotNil(t, existing)
}

func TestCreateOrUpdateRepoShouldUpdateRepo(t *testing.T) {
	clearDB()

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	repo := &Repo{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateRepo(db, repo, service)
	assert.Nil(t, err)

	name := "Updated"
	repo.Name = &name
	created, err := CreateOrUpdateRepo(db, repo, service)
	assert.Nil(t, err)
	assert.False(t, created)

	updated, err := FindRepoByID(db, repo.ID)
	assert.Nil(t, err)
	assert.Equal(t, "Updated", *updated.Name)
}
