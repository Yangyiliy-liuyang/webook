package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/domain"
	"webook/internal/repository"
	repomocks "webook/internal/repository/mocks"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleRepository)
		art     domain.Article
		wantId  int64
		wantErr error
	}{
		{
			name: "新建并发表成功",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleRepository) {
				authorRepo := repomocks.newMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				readerRepo := repomocks.newMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Title:   "新建发表",
				Content: "新建发表",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 1,
		},
		{
			name: "修改并发表成功",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleRepository) {
				authorRepo := repomocks.newMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				readerRepo := repomocks.newMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Id:      2,
				Title:   "新建发表",
				Content: "新建发表",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 2,
		},
		{
			name: "新建并发表失败",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleRepository) {
				authorRepo := repomocks.newMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				readerRepo := repomocks.newMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      3,
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("发表失败"))
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Title:   "新建发表",
				Content: "新建发表",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  3,
			wantErr: errors.New("发表失败"),
		},
		{
			name: "修改并发表失败",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleRepository) {
				authorRepo := repomocks.newMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      3,
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				readerRepo := repomocks.newMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      3,
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("发表失败"))
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Id:      3,
				Title:   "新建发表",
				Content: "新建发表",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  3,
			wantErr: errors.New("发表失败"),
		},
		{
			name: "修改保存至制作库失败",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleRepository) {
				authorRepo := repomocks.newMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      3,
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("保存至制作库失败"))
				readerRepo := repomocks.newMockArticleReaderRepository(ctrl)
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Id:      3,
				Title:   "新建发表",
				Content: "新建发表",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  3,
			wantErr: errors.New("保存至制作库失败"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authorRepo, readerRepo := tc.mock(ctrl)
			svc := NewArticleServiceV1(authorRepo, readerRepo)
			artId, err := svc.PublishV1(context.Background(), tc.art)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, artId)
		})
	}
}
