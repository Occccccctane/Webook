package Dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type=varchar(4096)"`
	Content string `gorm:"type=BLOB"`
	//作者ID
	AuthorId int64 `gorm:"index"`
	//状态： 1为存在0为删除
	Status int64
	Ctime  int64
	Utime  int64
}

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, entity Article) (int64, error)
}

type ArticleGormDao struct {
	db *gorm.DB
}

func NewArticleGormDao(db *gorm.DB) ArticleDao {
	return &ArticleGormDao{
		db: db,
	}
}

func (dao *ArticleGormDao) Update(ctx context.Context, art Article) (int64, error) {
	res := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? and author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"utime":   time.Now().UnixMilli(),
		})
	if res.Error != nil {
		return art.Id, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, errors.New("更新数据失败")
		//	可能作者不对或者帖子id不对
	}
	return art.Id, nil
}

func (dao *ArticleGormDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}
