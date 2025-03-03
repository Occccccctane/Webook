package Service

import (
	"GinStart/Domain"
	"GinStart/Repository"
	"GinStart/pkg/logger"
	"context"
	"errors"
)

type ArticleService interface {
	Save(ctx context.Context, art Domain.Article) (int64, error)
	Publish(ctx context.Context, art Domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid int64, aid int64) error
}

type ArticleServiceImpl struct {
	repo Repository.ArticleRepository
	//V1 Service进行分发
	readerRepo Repository.ArticleReadRepo
	authorRepo Repository.ArticleAuthorRepo
	l          logger.Logger
}

func NewArticleServiceImplV1(
	authorRepo Repository.ArticleAuthorRepo,
	readerRepo Repository.ArticleReadRepo,
	l logger.Logger) *ArticleServiceImpl {
	return &ArticleServiceImpl{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		l:          l,
	}

}

func NewArticleServiceImpl(repo Repository.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{
		repo: repo,
	}

}
func (a *ArticleServiceImpl) Withdraw(ctx context.Context, uid int64, aid int64) error {
	return a.repo.SyncStatus(ctx, uid, aid, Domain.ArticleStatusPrivate)
}

// PublishV1 在Service层将存储分发
func (a *ArticleServiceImpl) PublishV1(ctx context.Context, art Domain.Article) (int64, error) {
	art.Status = Domain.ArticleStatusPublished
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id

	for i := 0; i < 3; i++ {
		err = a.readerRepo.Save(ctx, art)
		if err != nil {
			a.l.Error("保存制作库成功同步线上库失败",
				logger.Int64("aid", art.Id),
				logger.Error(err))
		} else {
			return id, nil
		}

	}
	a.l.Error("保存线上库失败重试次数耗尽",
		logger.Int64("aid", art.Id),
		logger.Error(err))
	return id, errors.New("保存线上库失败重试次数耗尽")
}

// Publish 正常业务调用
func (a *ArticleServiceImpl) Publish(ctx context.Context, art Domain.Article) (int64, error) {
	art.Status = Domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}

func (a *ArticleServiceImpl) Save(ctx context.Context, art Domain.Article) (int64, error) {
	art.Status = Domain.ArticleStatusUnpublished
	if art.Id > 0 {
		return a.repo.Update(ctx, art)
	} else {
		return a.repo.Create(ctx, art)
	}

}
