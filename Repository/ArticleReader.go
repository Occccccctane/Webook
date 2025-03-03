package Repository

import (
	"GinStart/Domain"
	"context"
)

type ArticleReadRepo interface {
	// Save 新则创建，旧则更新
	Save(ctx context.Context, art Domain.Article) error
}
