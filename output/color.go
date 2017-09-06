package output

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
)

const defaultInterval = 300
const minInterval = 250
const defaultColor = "yellow"

var spin *spinner.Spinner
var cfg *config.OutputConfig

// Color is a color text output
type Color struct {
}

// Configure configures the output
func (c *Color) Configure(oc *config.OutputConfig) {
	cfg = oc
}

// Inline displays text in line
func (c *Color) Inline(s string) {
	fmt.Print(color.GreenString(s))
}

// Info displays information
func (c *Color) Info(s string) {
	color.Green(s)
}

// Error displays an error
func (c *Color) Error(s string) {
	color.Red(s)
}

// Fatal displays an error and ends the program
func (c *Color) Fatal(s string) {
	c.Error(s)
	os.Exit(1)
}

// Event displays an event {
func (c *Color) Event(event *model.Event) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.YellowString(event.Who))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.GreenString(fmt.Sprintf(" %s", event.What)))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.BlueString(fmt.Sprintf(" %s", event.Which)))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.RedString(fmt.Sprintf(" %s", event.URL)))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.MagentaString(fmt.Sprintf(" %s", humanize.Time(event.When))))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// StarLine displays a star in one line
func (c *Color) StarLine(star *model.Star) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(*star.FullName))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", star.Stargazers)))
	if err != nil {
		c.Error(err.Error())
	}

	if star.Language != nil {
		_, err = buffer.WriteString(color.GreenString(fmt.Sprintf(" %s", *star.Language)))
		if err != nil {
			c.Error(err.Error())
		}
	}

	if star.URL != nil {
		_, err = buffer.WriteString(color.RedString(fmt.Sprintf(" %s", *star.URL)))
		if err != nil {
			c.Error(err.Error())
		}
	}

	fmt.Println(buffer.String())
}

// Star displays a star
func (c *Color) Star(star *model.Star) {
	c.StarLine(star)

	if len(star.Tags) > 0 {
		var buffer bytes.Buffer
		leader := ""
		for _, tag := range star.Tags {
			_, err := buffer.WriteString(color.MagentaString(fmt.Sprintf("%s%s", leader, tag.Name)))
			if err != nil {
				c.Error(err.Error())
			}
			leader = ", "
		}
		fmt.Println(buffer.String())
	}

	if star.Description != nil && *star.Description != "" {
		color.White(*star.Description)
	}

	if star.Homepage != nil && *star.Homepage != "" {
		color.Red(fmt.Sprintf("Home page: %s", *star.Homepage))
	}

	color.Green(fmt.Sprintf("Starred on %s", star.StarredAt.Format(time.UnixDate)))
}

// Repo displays a repo
func (c *Color) Repo(repo *model.Repo) {
	c.StarLine(repo)

	if len(repo.Tags) > 0 {
		var buffer bytes.Buffer
		leader := ""
		for _, tag := range repo.Tags {
			_, err := buffer.WriteString(color.MagentaString(fmt.Sprintf("%s%s", leader, tag.Name)))
			if err != nil {
				c.Error(err.Error())
			}
			leader = ", "
		}
		fmt.Println(buffer.String())
	}

	if repo.Description != nil && *repo.Description != "" {
		color.White(*repo.Description)
	}

	if repo.Homepage != nil && *repo.Homepage != "" {
		color.Red(fmt.Sprintf("Home page: %s", *repo.Homepage))
	}

	color.Green(fmt.Sprintf("Created on %s", repo.CreatedAt.Format(time.UnixDate)))
}

// Tag displays a tag
func (c *Color) Tag(tag *model.Tag) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(tag.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", tag.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// Topic displays a topic
func (c *Color) Topic(topic *model.Topic) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(topic.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", topic.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// Academic displays a academic
func (c *Color) Academic(academic *model.Academic) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(academic.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", academic.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// Software displays a software
func (c *Color) Software(software *model.Software) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(software.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", software.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// Tree displays a tree
func (c *Color) Tree(file *model.Tree) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(file.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", file.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// LanguageDetected displays a languageDetected
func (c *Color) LanguageDetected(languageDetected *model.LanguageDetected) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(languageDetected.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", languageDetected.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// LanguageDetected displays a languageDetected
func (c *Color) Pkg(pkg *model.Pkg) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(pkg.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", pkg.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// LanguageDetected displays a languageDetected
func (c *Color) Readme(readme *model.Readme) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(readme.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", readme.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// Keyword displays a keyword
func (c *Color) Keyword(keyword *model.Keyword) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(keyword.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", keyword.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// Pattern displays a pattern
func (c *Color) Pattern(pattern *model.Pattern) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(pattern.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", pattern.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// Tick displays evidence that the program is working
func (c *Color) Tick() {
	if spin == nil {
		index := 0
		interval := defaultInterval
		clr := defaultColor
		if cfg != nil {
			index = cfg.SpinnerIndex
			if index < 0 || index > len(spinner.CharSets) {
				index = 0
			}
			interval = cfg.SpinnerInterval
			if interval < minInterval {
				interval = minInterval
			}
			clr = cfg.SpinnerColor
			if clr == "" {
				clr = defaultColor
			}
		}
		spin = spinner.New(spinner.CharSets[index], time.Duration(interval)*time.Millisecond)
		spin.Suffix = color.CyanString(" Updating")
		if err := spin.Color(clr); err != nil {
			c.Error(err.Error())
		}
	}
}

func init() {
	registerOutput(&Color{})
}
