package Service

import (
	"GinStart/Domain"
	"GinStart/Repository"
	"context"
)

type ArticleService interface {
	Save(ctx context.Context, art Domain.Article) (int64, error)
	Publish(ctx context.Context, art Domain.Article) (int64, error)
}

type ArticleServiceImpl struct {
	repo Repository.ArticleRepository
}

func NewArticleServiceImpl(repo Repository.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{
		repo: repo,
	}

}

func (a *ArticleServiceImpl) Publish(ctx context.Context, art Domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *ArticleServiceImpl) Save(ctx context.Context, art Domain.Article) (int64, error) {
	if art.Id > 0 {
		return a.repo.Update(ctx, art)
	} else {
		return a.repo.Create(ctx, art)
	}
}
