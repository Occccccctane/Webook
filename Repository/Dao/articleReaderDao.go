package Dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleReaderDao interface {
	Upsert(ctx context.Context, art Article) error
	UpsertV2(ctx context.Context, art PublishedArticle) error
}
type ArticleGORMReaderDao struct {
	db *gorm.DB
}

func NewArticleGORMReaderDao(db *gorm.DB) ArticleReaderDao {
	return &ArticleGORMReaderDao{db: db}
}

func (a *ArticleGORMReaderDao) Upsert(ctx context.Context, art Article) error {
	panic("implement me")
}
func (a *ArticleGORMReaderDao) UpsertV2(ctx context.Context, art PublishedArticle) error {
	panic("implement me")
}
