package Dao

import (
	"context"
	"gorm.io/gorm"
)

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type=varchar(4096)"`
	Content string `gorm:"type=BLOB"`
	//作者ID
	AuthorId int64 `gorm:"index"`
	Ctime    int64
	Utime    int64
}

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

type ArticleGormDao struct {
	db *gorm.DB
}

func NewArticleGormDao(db *gorm.DB) ArticleDao {
	return &ArticleGormDao{
		db: db,
	}
}
func (a *ArticleGormDao) Insert(ctx context.Context, art Article) (int64, error) {
	err := a.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}
