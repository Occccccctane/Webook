package Repository

import (
	"GinStart/Domain"
	"context"
)

type ArticleAuthorRepo interface {
	Create(ctx context.Context, art Domain.Article) (int64, error)
	Update(ctx context.Context, art Domain.Article) error
}
