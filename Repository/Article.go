package Repository

import (
	"GinStart/Domain"
	"GinStart/Repository/Dao"
	"context"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art Domain.Article) (int64, error)
	Update(ctx context.Context, art Domain.Article) (int64, error)
	Sync(ctx context.Context, art Domain.Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, aid int64, private Domain.ArticleStatus) error
}

type ArticleRepositoryImpl struct {
	dao       Dao.ArticleDao
	readerDao Dao.ArticleReaderDao
	authorDao Dao.ArticleAuthorDao
	db        *gorm.DB
}

func NewArticleRepositoryImpl(dao Dao.ArticleDao) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao: dao,
	}
}

func (a *ArticleRepositoryImpl) SyncStatus(ctx context.Context, uid int64, aid int64, private Domain.ArticleStatus) error {
	return a.dao.SyncStatus(ctx, uid, aid, private.ToUint8())
}

// NewArticleRepositoryImplV2 在repository层进行分发
func NewArticleRepositoryImplV2(readerDao Dao.ArticleReaderDao, authorDao Dao.ArticleAuthorDao) *ArticleRepositoryImpl {
	return &ArticleRepositoryImpl{
		readerDao: readerDao,
		authorDao: authorDao,
	}
}

// Sync 在Dao层进行分发
func (a *ArticleRepositoryImpl) Sync(ctx context.Context, art Domain.Article) (int64, error) {
	return a.dao.Sync(ctx, a.toEntity(art))
}

// SyncV2 同步同库事务实现
func (a *ArticleRepositoryImpl) SyncV2(ctx context.Context, art Domain.Article) (int64, error) {
	tx := a.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	//防止业务panic，将事务数据进行回滚
	defer tx.Rollback()
	authorDao := Dao.NewArticleGORMAuthorDao(tx)
	readerDao := Dao.NewArticleGORMReaderDao(tx)
	artEntity := a.toEntity(art)
	var (
		id  int64
		err error
	)
	id = art.Id
	if art.Id > 0 {
		err = authorDao.Update(ctx, artEntity)
	} else {
		id, err = authorDao.Create(ctx, artEntity)
	}
	if err != nil {
		return 0, err
	}
	artEntity.Id = id
	err = readerDao.UpsertV2(ctx, Dao.PublishedArticle(artEntity))
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}

// SyncV1 在repository层进行分发同步,使用不同数据库存储，不使用事务
func (a *ArticleRepositoryImpl) SyncV1(ctx context.Context, art Domain.Article) (int64, error) {
	artEntity := a.toEntity(art)
	var (
		id  int64
		err error
	)
	id = art.Id
	if art.Id > 0 {
		err = a.authorDao.Update(ctx, artEntity)
	} else {
		id, err = a.authorDao.Create(ctx, artEntity)
	}
	if err != nil {
		return 0, err
	}
	artEntity.Id = id
	err = a.readerDao.Upsert(ctx, artEntity)
	return id, err
}

func (a *ArticleRepositoryImpl) Update(ctx context.Context, art Domain.Article) (int64, error) {
	err := a.dao.Update(ctx, a.toEntity(art))

	return art.Id, err
}

func (a *ArticleRepositoryImpl) Create(ctx context.Context, art Domain.Article) (int64, error) {
	return a.dao.Insert(ctx, a.toEntity(art))
}

func (a *ArticleRepositoryImpl) toEntity(art Domain.Article) Dao.Article {
	return Dao.Article{
		Id:       art.Id,
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
		Status:   art.Status.ToUint8(),
	}
}
