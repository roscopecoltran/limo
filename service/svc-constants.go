package service

// github
const GITHUB__RAW_URL = "https://raw.githubusercontent.com/"

// github - job queue(s)
const (
	NOTSTART = "NotStart"
	FETCHING = "Fetching"
	INDEXING = "Indexing"
	INDEXED  = "Indexed"
	ERROR    = "Error"
)