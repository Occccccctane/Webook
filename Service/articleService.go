package Service

import (
	"GinStart/Domain"
	"GinStart/Repository"
	"context"
)

type ArticleService interface {
	Save(ctx context.Context, art Domain.Article) (int64, error)
}

type ArticleServiceImpl struct {
	repo Repository.ArticleRepository
}

func NewArticleServiceImpl(repo Repository.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{
		repo: repo,
	}

}

func (a *ArticleServiceImpl) Save(ctx context.Context, art Domain.Article) (int64, error) {
	return a.repo.Create(ctx, art)
}
