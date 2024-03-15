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
		wantErr bool
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleRepository) {
				authorRepo := repomocks.newMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.newMockArticleReaderRepository(ctrl)
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
