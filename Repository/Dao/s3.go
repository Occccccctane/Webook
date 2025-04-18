package Dao

import (
	"GinStart/Domain"
	"GinStart/pkg/Converser"
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

type S3ArticleDao struct {
	GormArticleDao
	oss *s3.S3
}

func NewS3ArticleDao(db *gorm.DB, oss *s3.S3) *S3ArticleDao {
	return &S3ArticleDao{
		GormArticleDao: GormArticleDao{
			db: db,
		},
		oss: oss,
	}
}

func (d *S3ArticleDao) Sync(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		dao := NewGormArticleDao(tx)
		id = art.Id
		if art.Id > 0 {
			err = dao.Update(ctx, art)
		} else {
			id, err = dao.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		art.Id = id
		now := time.Now().UnixMilli()
		pubArt := S3PublishedArticle{
			Id:       art.Id,
			Title:    art.Title,
			AuthorId: art.AuthorId,
			Status:   art.Status,
			Ctime:    now,
			Utime:    now,
		}
		err = tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":  pubArt.Title,
				"status": pubArt.Status,
				"utime":  pubArt.Utime,
			}),
		}).Create(&pubArt).Error
		return err
	})
	if err != nil {
		return 0, err
	}
	_, err = d.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      Converser.ToPtr[string]("webook"),
		Key:         Converser.ToPtr[string](strconv.FormatInt(art.Id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: Converser.ToPtr("text/plain;charset=utf-8"),
	})
	return id, err
}

func (d *S3ArticleDao) SyncStatus(ctx context.Context, uid int64, aid int64, status uint8) error {
	now := time.Now().UnixMilli()
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? and author_id = ?", aid, uid).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("非法数据更新")
		}
		return tx.Model(&S3PublishedArticle{}).
			Where("id = ? ", aid).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			}).Error
	})
	if err != nil {
		return err
	}
	if status == Domain.ArticleStatusPrivate {
		_, err = d.oss.DeleteObjectsWithContext(ctx, &s3.DeleteObjectsInput{
			Bucket: Converser.ToPtr[string]("webook"),
			Delete: &s3.Delete{
				Objects: []*s3.ObjectIdentifier{
					{
						Key: Converser.ToPtr[string](strconv.FormatInt(aid, 10)),
					},
				},
			},
		})
	}
	return err
}

type S3PublishedArticle struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Title string `gorm:"type=varchar(4096)" `
	//作者ID
	AuthorId int64 `gorm:"index" `
	Status   uint8
	Ctime    int64
	Utime    int64
}
