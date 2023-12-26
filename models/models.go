package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	Username string `json:"username" gorm:"not null"`
	Posts    []Post `json:"posts" gorm:"foreignKey:UserID"`
}

type Post struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Title   string `json:"title" gorm:"not null"`
	Content string `json:"content" gorm:"not null"`
	UserID  uint   `json:"userID" gorm:"index"`
	// set default to null to allow deleting users and not their posts
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID"`
}

type Comment struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Text   string `json:"text" gorm:"not null"`
	UserID uint   `json:"userID"`
	PostID uint   `json:"postID" gorm:"index;not null"`
}
