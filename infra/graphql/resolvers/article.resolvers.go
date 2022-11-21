package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/TudorHulban/rest-articles/app/service"
	"github.com/TudorHulban/rest-articles/infra/graphql/generated"
	"github.com/TudorHulban/rest-articles/infra/graphql/generated/models"
)

// UpdateArticleTitle is the resolver for the UpdateArticleTitle field.
func (r *mutationResolver) UpdateArticleTitle(ctx context.Context, id int, title string) (int, error) {
	paramsArticle := service.ParamsUpdateArticle{
		ID:    int64(id),
		Title: &title,
	}

	if errUpdate := r.serv.UpdateArticle(ctx, &paramsArticle); errUpdate != nil {
		return 0, errUpdate
	}

	return id, nil
}

// UpdateArticleURL is the resolver for the UpdateArticleURL field.
func (r *mutationResolver) UpdateArticleURL(ctx context.Context, id int, url string) (int, error) {
	paramsArticle := service.ParamsUpdateArticle{
		ID:  int64(id),
		URL: &url,
	}

	if errUpdate := r.serv.UpdateArticle(ctx, &paramsArticle); errUpdate != nil {
		return 0, errUpdate
	}

	return id, nil
}

// GetArticle is the resolver for the GetArticle field.
func (r *queryResolver) GetArticle(ctx context.Context, id int) (*models.Article, error) {
	reconstructedArticle, errGet := r.serv.GetArticle(ctx, int64(id))
	if errGet != nil {
		return nil, errGet // TODO: move to area error
	}

	return &models.Article{
		ID:    int(reconstructedArticle.ID),
		Title: reconstructedArticle.Title,
		URL:   reconstructedArticle.URL,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
