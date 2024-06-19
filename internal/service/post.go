package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"tgstat/internal/entity"
	"tgstat/internal/pkg/client"
)

type PostService struct {
	client *client.Client
}

func NewPostService(client *client.Client) *PostService {
	return &PostService{
		client: client,
	}
}

func (s *PostService) GetAll(ctx context.Context, v url.Values) (*entity.Posts, error) {
	result, err := s.client.CallAPI(ctx, &client.Request{
		Method:   http.MethodGet,
		Endpoint: "/channels/posts",
		Query:    v,
	})
	if err != nil {
		return nil, err
	}

	posts := new(entity.Posts)
	if result.Status != "ok" {
		return posts, result
	}
	if err := json.Unmarshal(result.Response, posts); err != nil {
		return posts, err
	}
	return posts, nil
}
