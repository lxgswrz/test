package db

import (
	"errors"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Content string `gorm:"not null" json:"content"`
}

func GetComments(page, size int) ([]Comment, int64, error) {
	var total int64
	result := DB.Model(&Comment{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	offset := (page - 1) * size

	var comments []Comment
	result = DB.Order("id DESC").Offset(offset).Limit(size).Find(&comments)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return comments, total, nil
}

func AddComment(name, content string) (Comment, error) {
	comment := Comment{
		Name:    name,
		Content: content,
	}

	result := DB.Create(&comment)
	if result.Error != nil {
		return Comment{}, result.Error
	}

	return comment, nil
}

func DeleteComment(id int) error {
	result := DB.Delete(&Comment{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}
