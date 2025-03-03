package Dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleAuthorDao interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}
type ArticleGORMAuthorDao struct {
	db *gorm.DB
}

func NewArticleGORMAuthorDao(db *gorm.DB) ArticleAuthorDao {
	return &ArticleGORMAuthorDao{db: db}
}

func (a *ArticleGORMAuthorDao) Create(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *ArticleGORMAuthorDao) Update(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}
