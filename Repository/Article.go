package Repository

import (
	"GinStart/Domain"
	"GinStart/Repository/Dao"
	"context"
)

type ArticleRepository interface {
	Create(ctx context.Context, art Domain.Article) (int64, error)
}

type ArticleRepositoryImpl struct {
	dao Dao.ArticleDao
}

func NewArticleRepositoryImpl(dao Dao.ArticleDao) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao: dao,
	}
}
func (a *ArticleRepositoryImpl) Create(ctx context.Context, art Domain.Article) (int64, error) {
	return a.dao.Insert(ctx, a.toEntity(art))
}

func (a *ArticleRepositoryImpl) toEntity(art Domain.Article) Dao.Article {
	return Dao.Article{
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
	}
}
