package model

import (
	//"fmt"
	//"strings"
	//"time"
	// "github.com/jinzhu/gorm"
	"github.com/google/go-github/github"
	//"github.com/sirupsen/logrus"
)

/*
	VCS - github bucket
*/
type GatewayBucket_Github struct {
	// gorm.Model
	Disable 		bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Repos 			[]GatewayBucket_GithubRepository 		`yaml:"repos,omitempty" json:"repos,omitempty"` 		//
	Search 			[]GatewayBucket_GithubSearch 			`yaml:"repos,omitempty" json:"repos,omitempty"` 		//
	Related  		[]GatewayBucket_GithubRelated 			`yaml:"related,omitempty" json:"related,omitempty"` 	// https://github.com/google/go-github/blob/master/github/repos.go#L22-L117
	Starred 		[]github.StarredRepository 				`yaml:"starred,omitempty" json:"starred,omitempty"`		// https://github.com/google/go-github/blob/master/github/activity_star.go#L13-L17
	Options 		GatewayGlobal_SearchOptions 			`yaml:"options,omitempty" json:"options,omitempty"` 	// 
}

type GatewayBucket_GithubRepository struct {
	// gorm.Model
	Repo 			*github.Repository 		  				`yaml:"repos,omitempty" json:"repos,omitempty"`			// https://github.com/google/go-github/blob/master/github/repos.go#L22-L117
	Readme 			*github.RepositoryContent 				`yaml:"readme,omitempty" json:"readme,omitempty"`		// https://github.com/google/go-github/blob/master/github/repos_contents.go#L22-L38
	Trees 			[]GatewayBucket_GithubTree 				`yaml:"trees,omitempty" json:"trees,omitempty"` 		//
	Languages 		map[string]int 							`yaml:"language,omitempty" json:"language,omitempty"` 	//
	User 			*github.User  							`yaml:"user,omitempty" json:"user,omitempty"` 			//
	//LanguageDetection 	map[string]int 					`yaml:"language_detected,omitempty" json:"language_detected,omitempty"`
}

type GatewayBucket_GithubRepoExtraInfo struct {
	// gorm.Model
	Readme 			*github.RepositoryContent 				`yaml:"readme,omitempty" json:"readme,omitempty"`		// https://github.com/google/go-github/blob/master/github/repos_contents.go#L22-L38
	Trees 			[]GatewayBucket_GithubTree 				`yaml:"trees,omitempty" json:"trees,omitempty"` 		// 
	Languages 		map[string]int 							`yaml:"language,omitempty" json:"language,omitempty"` 	// 
	User 			*github.User  							`yaml:"user,omitempty" json:"user,omitempty"` 			// 
	//LanguageDetection 	map[string]int 					`yaml:"language_detected,omitempty" json:"language_detected,omitempty"`
}

// https://github.com/google/go-github/blob/master/github/git_trees.go#L13-L34
// /repos/:owner/:repo/git/trees/:sha
type GatewayBucket_GithubSearch struct {
	// gorm.Model
	Query 				string 								`yaml:"query,omitempty" json:"query,omitempty"` 		//
	Hash 				string 								`yaml:"sha,omitempty" json:"sha,omitempty"` 			//
	Date 				string 								`yaml:"date,omitempty" json:"date,omitempty"` 			//
	UnixTime 			string 								`yaml:"unixtime,omitempty" json:"unixtime,omitempty"` 	//
	Options  			[]github.SearchOptions 				`yaml:"options,omitempty" json:"options,omitempty"`		// https://github.com/google/go-github/blob/master/github/search.go#L29-L49
	Results struct {
		Repos 			[]github.RepositoriesSearchResult 	`yaml:"repos,omitempty" json:"repos,omitempty"`			// https://github.com/google/go-github/blob/master/github/search.go#L52-L56
		Codes 			[]github.CodeSearchResult 			`yaml:"codes,omitempty" json:"codes,omitempty"`			// https://github.com/google/go-github/blob/master/github/search.go#L149-L154
		Users 			[]github.UsersSearchResult 			`yaml:"users,omitempty" json:"users,omitempty"`			// https://github.com/google/go-github/blob/master/github/search.go#L114-L119
		Issues 			[]github.IssuesSearchResult 		`yaml:"issues,omitempty" json:"issues,omitempty"`		// https://github.com/google/go-github/blob/master/github/search.go#L99-L102
		Commits 		[]github.CommitsSearchResult 		`yaml:"commits,omitempty" json:"commits,omitempty"`		// https://github.com/google/go-github/blob/master/github/search.go#L67-L72	
	}
}

type GatewayBucket_GithubRelated struct {
	// gorm.Model
	Disable 		bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	AgentName 		string 									`defaullt:"sniperkit" yaml:"agent_name,omitempty" json:"agent_name,omitempty"`
	AgentType 		string 									`defaullt:"first_generation" yaml:"agent_type,omitempty" json:"agent_type,omitempty"`
	Related 		[]GatewayBucket_GithubSearch 			`yaml:"results,omitempty" json:"results,omitempty"`
	Featured 		[]GatewayLink 							`yaml:"highlights,omitempty" json:"highlights,omitempty"` 		// https://github.com/google/go-github/blob/master/github/git_trees.go#L13-L34
	Options 		GatewayGlobal_SearchOptions 			`yaml:"options,omitempty" json:"options,omitempty"` 			// 
}

type GatewayBucket_GithubTree struct {
	// gorm.Model
	Disable 		bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Branch 			string 									`yaml:"branch,omitempty" json:"branch,omitempty"`
	Recursive 		bool 									`yaml:"recursive,omitempty" json:"recursive,omitempty"`
	Files 			[]github.Tree 							`yaml:"files,omitempty" json:"files,omitempty"` 	// https://github.com/google/go-github/blob/master/github/git_trees.go#L13-L34
}

type GatewayBucket_GithubSearchOptions Global_SearchOptions

