package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMilestoneEquals(t *testing.T) {
	left := &ExtraGithub_Milestone{}
	right := &ExtraGithub_Milestone{}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{Title: "foo"}
	right = &ExtraGithub_Milestone{Title: "foo"}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{State: "open"}
	right = &ExtraGithub_Milestone{State: "open"}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{Description: "bar"}
	right = &ExtraGithub_Milestone{Description: "bar"}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{Title: "foo"}
	right = &ExtraGithub_Milestone{Title: "bar"}
	assert.False(t, left.Equals(right))

	left = &ExtraGithub_Milestone{Description: "bar"}
	right = &ExtraGithub_Milestone{Description: "baz"}
	assert.False(t, left.Equals(right))

	left = &ExtraGithub_Milestone{State: "open"}
	right = &ExtraGithub_Milestone{State: "closed"}
	assert.False(t, left.Equals(right))

	left = &Milestone{PreviousTitles: nil}
	right = &Milestone{PreviousTitles: nil}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{PreviousTitles: []string{}}
	right = &ExtraGithub_Milestone{PreviousTitles: []string{}}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{PreviousTitles: []string{"foo"}}
	right = &ExtraGithub_Milestone{PreviousTitles: []string{}}
	assert.False(t, left.Equals(right))

	left = &ExtraGithub_Milestone{PreviousTitles: []string{"foo", "bar", "baz"}}
	right = &ExtraGithub_Milestone{PreviousTitles: []string{"bar", "foo", "baz"}}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{PreviousTitles: []string{"foo", "bar", "baz"}}
	right = &ExtraGithub_Milestone{PreviousTitles: []string{"bar", "foo", "qux"}}
	assert.False(t, left.Equals(right))

	left = &ExtraGithub_Milestone{Number: 0}
	right = &ExtraGithub_Milestone{Number: 0}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{Number: 42}
	right = &ExtraGithub_Milestone{Number: 43}
	assert.False(t, left.Equals(right))

	left = &ExtraGithub_Milestone{DueOn: time.Time{}}
	right = &ExtraGithub_Milestone{DueOn: time.Time{}}
	assert.True(t, left.Equals(right))

	today, _ := time.Parse("2006-01-02T15:04:05Z0700", "2017-01-10T23:59:59Z")
	morning, _ := time.Parse("2006-01-02T15:04:05Z0700", "2017-01-10T06:00:00Z")
	lastmonth, _ := time.Parse("2006-01-02T15:04:05Z0700", "2016-12-10T23:59:59Z")
	left = &ExtraGithub_Milestone{DueOn: today}
	right = &ExtraGithub_Milestone{DueOn: today}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{}
	right = &ExtraGithub_Milestone{DueOn: today}
	assert.False(t, left.Equals(right))

	assert.False(t, left.Equals(nil))

	left = &ExtraGithub_Milestone{DueOn: today}
	right = &ExtraGithub_Milestone{DueOn: morning}
	assert.True(t, left.Equals(right))

	left = &ExtraGithub_Milestone{DueOn: today}
	right = &ExtraGithub_Milestone{DueOn: lastmonth}
	assert.False(t, left.Equals(right))
}