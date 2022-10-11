package models

type Comment struct {
	Model
	Text string `json:"text"`
	Likes []*User `json:"likes" gorm:"many2many:user_comments;"`
	RepliedID uint `json:"repliedId"`
	UserID uint `json:"userId"`
	Comments []Comment `json:"comments" gorm:"foreignKey:RepliedID"`
}
