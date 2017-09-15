package news

import "github.com/jinzhu/gorm"

/*
User model
*/
type User struct {
	ID       uint   `json:"id" yaml:"id"  gorm:"primary_key"`
	Name     string `json:"name" yaml:"name"  sql:"size:255"`
	Username string `json:"username" yaml:"username"  sql:"size:255;unique_index"`
	Email    string `json:"email" yaml:"email"  sql:"size:255"`
	Password string `json:"-" sql:"size:255"`
}

// ---------------------------------------------------------------------------

/*
Save user on database
*/
func (u *User) Save(db *gorm.DB) (*User, error) {
	if err := db.Create(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

/*
GetUser returns the user with that id
*/
func GetUser(db *gorm.DB, id int) (*User, error) {
	var user User
	if err := db.Where("ID = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

/*
GetUsers returns collection of users
*/
func GetUsers(db *gorm.DB) (*[]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

/*
FindUserByUsername find user
*/
func FindUserByUsername(db *gorm.DB, username string) (*User, error) {
	var users []User
	if err := db.Where("username = ?", username).Find(&users).Error; err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}
