package Service

import (
	"GinStart/Domain"
	"GinStart/Repository"
	repositorymock "GinStart/Repository/mocks"
	"GinStart/pkg/logger"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestArticleServiceImpl_Publish(t *testing.T) {
	testCases := []struct {
		name        string
		mock        func(ctrl *gomock.Controller) (Repository.ArticleAuthorRepo, Repository.ArticleReadRepo)
		art         Domain.Article
		ExpectedId  int64
		ExpectedErr error
	}{
		{
			name: "新建成功",
			mock: func(ctrl *gomock.Controller) (Repository.ArticleAuthorRepo, Repository.ArticleReadRepo) {
				authRepo := repositorymock.NewMockArticleAuthorRepo(ctrl)
				authRepo.EXPECT().Create(gomock.Any(), Domain.Article{
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(int64(1), nil)
				readerRepo := repositorymock.NewMockArticleReadRepo(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), Domain.Article{
					Id: 1,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(nil)
				return authRepo, readerRepo
			},
			art: Domain.Article{
				Author: Domain.Author{
					Id: 123,
				},
				Content: "测试内容",
				Title:   "测试标题",
			},
			ExpectedId:  1,
			ExpectedErr: nil,
		},
		{
			name: "修改帖子并新发表失败，重试成功",
			mock: func(ctrl *gomock.Controller) (Repository.ArticleAuthorRepo, Repository.ArticleReadRepo) {
				authRepo := repositorymock.NewMockArticleAuthorRepo(ctrl)
				authRepo.EXPECT().Update(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(nil)
				readerRepo := repositorymock.NewMockArticleReadRepo(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(errors.New("mock db err"))
				readerRepo.EXPECT().Save(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(nil)
				return authRepo, readerRepo
			},
			art: Domain.Article{
				Id: 11,
				Author: Domain.Author{
					Id: 123,
				},
				Content: "测试内容",
				Title:   "测试标题",
			},
			ExpectedId:  11,
			ExpectedErr: nil,
		},
		{
			name: "修改帖子并新发表失败，重试失败",
			mock: func(ctrl *gomock.Controller) (Repository.ArticleAuthorRepo, Repository.ArticleReadRepo) {
				authRepo := repositorymock.NewMockArticleAuthorRepo(ctrl)
				authRepo.EXPECT().Update(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(nil)
				readerRepo := repositorymock.NewMockArticleReadRepo(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(errors.New("mock db err"))
				readerRepo.EXPECT().Save(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(errors.New("mock db err"))
				readerRepo.EXPECT().Save(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(errors.New("mock db err"))
				return authRepo, readerRepo
			},
			art: Domain.Article{
				Id: 11,
				Author: Domain.Author{
					Id: 123,
				},
				Content: "测试内容",
				Title:   "测试标题",
			},
			ExpectedId:  11,
			ExpectedErr: errors.New("保存线上库失败重试次数耗尽"),
		},
		{
			name: "保存到制作库失败",
			mock: func(ctrl *gomock.Controller) (Repository.ArticleAuthorRepo, Repository.ArticleReadRepo) {
				authRepo := repositorymock.NewMockArticleAuthorRepo(ctrl)
				authRepo.EXPECT().Update(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(errors.New("mock db err"))
				readerRepo := repositorymock.NewMockArticleReadRepo(ctrl)
				return authRepo, readerRepo
			},
			art: Domain.Article{
				Id: 11,
				Author: Domain.Author{
					Id: 123,
				},
				Content: "测试内容",
				Title:   "测试标题",
			},
			ExpectedId:  0,
			ExpectedErr: errors.New("mock db err"),
		},
		{
			name: "修改帖子并重发表成功",
			mock: func(ctrl *gomock.Controller) (Repository.ArticleAuthorRepo, Repository.ArticleReadRepo) {
				authRepo := repositorymock.NewMockArticleAuthorRepo(ctrl)
				authRepo.EXPECT().Update(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(nil)
				readerRepo := repositorymock.NewMockArticleReadRepo(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), Domain.Article{
					Id: 11,
					Author: Domain.Author{
						Id: 123,
					},
					Content: "测试内容",
					Title:   "测试标题",
				}).Return(nil)
				return authRepo, readerRepo
			},
			art: Domain.Article{
				Id: 11,
				Author: Domain.Author{
					Id: 123,
				},
				Content: "测试内容",
				Title:   "测试标题",
			},
			ExpectedId:  11,
			ExpectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authorRepo, readerRepo := tc.mock(ctrl)
			l := logger.NewNopLogger()
			svc := NewArticleServiceImplV1(authorRepo, readerRepo, l)
			id, err := svc.PublishV1(context.Background(), tc.art)
			assert.Equal(t, tc.ExpectedErr, err)
			assert.Equal(t, tc.ExpectedId, id)
		})
	}
}
