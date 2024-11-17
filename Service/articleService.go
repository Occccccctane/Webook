package Service

import (
	"GinStart/Domain"
	"context"
)

type ArticleService interface {
	Save(ctx context.Context, art Domain.Article) (int64, error)
}

type ArticleServiceImpl struct {
}

func (a *ArticleServiceImpl) Save(ctx context.Context, art Domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}
