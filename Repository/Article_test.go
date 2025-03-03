package Repository

import (
	"GinStart/Domain"
	"GinStart/Repository/Dao"
	daomock "GinStart/Repository/mocks/dao"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestArticleRepositoryImpl_SyncV1(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(m *gomock.Controller) (Dao.ArticleReaderDao, Dao.ArticleAuthorDao)
		art     Domain.Article
		wantId  int64
		wantErr error
	}{
		{
			name: "新建同步成功",
			mock: func(m *gomock.Controller) (Dao.ArticleReaderDao, Dao.ArticleAuthorDao) {
				adao := daomock.NewMockArticleAuthorDao(m)
				adao.EXPECT().Create(gomock.Any(), Dao.Article{
					Content:  "content",
					Title:    "title",
					AuthorId: 123,
				}).Return(int64(1), nil)
				rdao := daomock.NewMockArticleReaderDao(m)
				rdao.EXPECT().Upsert(gomock.Any(), Dao.Article{
					Id:       1,
					Content:  "content",
					Title:    "title",
					AuthorId: 123,
				}).Return(nil)
				return rdao, adao
			},
			art: Domain.Article{
				Author: Domain.Author{
					Id: 123,
				},
				Content: "content",
				Title:   "title",
			},
			wantId:  1,
			wantErr: nil,
		},
		{
			name: "修改同步成功",
			mock: func(m *gomock.Controller) (Dao.ArticleReaderDao, Dao.ArticleAuthorDao) {
				adao := daomock.NewMockArticleAuthorDao(m)
				adao.EXPECT().Update(gomock.Any(), Dao.Article{
					Id:       11,
					Content:  "content",
					Title:    "title",
					AuthorId: 123,
				}).Return(nil)
				rdao := daomock.NewMockArticleReaderDao(m)
				rdao.EXPECT().Upsert(gomock.Any(), Dao.Article{
					Id:       11,
					Content:  "content",
					Title:    "title",
					AuthorId: 123,
				}).Return(nil)
				return rdao, adao
			},
			art: Domain.Article{
				Id: 11,
				Author: Domain.Author{
					Id: 123,
				},
				Content: "content",
				Title:   "title",
			},
			wantId:  11,
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			readerDao, authorDao := tc.mock(ctrl)
			repo := NewArticleRepositoryImplV2(readerDao, authorDao)
			Id, err := repo.SyncV1(context.Background(), tc.art)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, Id)
		})
	}
}
