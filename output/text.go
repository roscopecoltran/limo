package output

import (
	"bytes"
	"fmt"
	"os"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/roscopecoltran/sniperkit-limo/config"
	"github.com/roscopecoltran/sniperkit-limo/model"
)

// Text is a monochrome text output
type Text struct {
}

// Configure no-ops
func (t *Text) Configure(oc *config.OutputConfig) {
}

// Inline displays text in line
func (t *Text) Inline(s string) {
	fmt.Print(s)
}

// Info displays information
func (t *Text) Info(s string) {
	fmt.Println(s)
}

// Error displays an error
func (t *Text) Error(s string) {
	fmt.Fprintln(os.Stderr, s)
}

// Fatal displays an error and ends the program
func (t *Text) Fatal(s string) {
	t.Error(s)
	os.Exit(1)
}

// Event displays an event {
func (t *Text) Event(event *model.Event) {
	fmt.Printf("%s %s %s %s %s\n", event.Who, event.What, event.Which, event.URL, humanize.Time(event.When))
}

// StarLine displays a star in one line
func (t *Text) StarLine(star *model.Star) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(*star.FullName)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = buffer.WriteString(fmt.Sprintf(" *:%d", star.Stargazers))
	if err != nil {
		t.Error(err.Error())
	}

	if star.Language != nil {
		_, err := buffer.WriteString(fmt.Sprintf(" %s", *star.Language))
		if err != nil {
			t.Error(err.Error())
		}
	}

	if star.URL != nil {
		_, err := buffer.WriteString(fmt.Sprintf(" %s", *star.URL))
		if err != nil {
			t.Error(err.Error())
		}
	}
	fmt.Println(buffer.String())
}

// StarLine displays a star in one line
func (t *Text) RepoLine(repo *model.Repo) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(*repo.FullName)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = buffer.WriteString(fmt.Sprintf(" *:%d", repo.Stargazers))
	if err != nil {
		t.Error(err.Error())
	}

	if repo.Language != nil {
		_, err := buffer.WriteString(fmt.Sprintf(" %s", *repo.Language))
		if err != nil {
			t.Error(err.Error())
		}
	}

	if repo.URL != nil {
		_, err := buffer.WriteString(fmt.Sprintf(" %s", *repo.URL))
		if err != nil {
			t.Error(err.Error())
		}
	}
	fmt.Println(buffer.String())
}

// Star displays a star
func (t *Text) Star(star *model.Star) {
	t.StarLine(star)

	if len(star.Tags) > 0 {
		var buffer bytes.Buffer
		leader := ""
		for _, tag := range star.Tags {
			_, err := buffer.WriteString(fmt.Sprintf("%s%s", leader, tag.Name))
			if err != nil {
				t.Error(err.Error())
			}
			leader = ", "
		}
		fmt.Println(buffer.String())
	}

	if star.Description != nil && *star.Description != "" {
		fmt.Println(*star.Description)
	}

	if star.Homepage != nil && *star.Homepage != "" {
		fmt.Printf("Home page: %s\n", *star.Homepage)
	}

	fmt.Printf("Starred on %s\n", star.StarredAt.Format(time.UnixDate))
}

// Star displays a star
func (t *Text) Repo(repo *model.Repo) {
	t.RepoLine(repo)

	if len(repo.Tags) > 0 {
		var buffer bytes.Buffer
		leader := ""
		for _, tag := range repo.Tags {
			_, err := buffer.WriteString(fmt.Sprintf("%s%s", leader, tag.Name))
			if err != nil {
				t.Error(err.Error())
			}
			leader = ", "
		}
		fmt.Println(buffer.String())
	}

	if repo.Description != nil && *repo.Description != "" {
		fmt.Println(*repo.Description)
	}

	if repo.Homepage != nil && *repo.Homepage != "" {
		fmt.Printf("Home page: %s\n", *repo.Homepage)
	}

	fmt.Printf("Created on %s\n", repo.CreatedAt.Format(time.UnixDate))
}

// Tag displays a tag
func (t *Text) Tag(tag *model.Tag) {
	fmt.Printf("%s *:%d\n", tag.Name)
}

// Topic displays a topic
func (t *Text) Topic(topic *model.Topic) {
	fmt.Printf("%s *:%d\n", topic.Name)
}

// Tree displays a tree
func (t *Text) Tree(tree *model.Tree) {
	fmt.Printf("%s/tree/%s\n", tree.RemoteURI, tree.SHA)
}

// Software displays a software
func (t *Text) LanguageDetected(LanguageDetected *model.LanguageDetected) {
	fmt.Printf("%s \n", LanguageDetected.Name)
}

// Readme displays a readme
func (t *Text) Readme(readme *model.Readme) {
	fmt.Printf("filename: %s (size: %d)\n %s \n", readme.Name, readme.Size, readme.Decoded)
	//fmt.Printf("%s *:%d\n", readme.Name, readme.StarCount)
}

// Tick displays evidence that the program is working
func (t *Text) Tick() {
	fmt.Print(".")
}

func init() {
	registerOutput(&Text{})
}
