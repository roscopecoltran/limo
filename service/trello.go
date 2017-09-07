package service

import (
	"github.com/adlio/trello"
)

// Trello holds data that is used to authenticate to the trello API.
type Trello struct {
	*trello.Client
}