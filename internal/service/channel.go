package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"tgstat/internal/entity"
	"tgstat/internal/pkg/client"
)

type ChannelService struct {
	client *client.Client
}

func NewChannelService(client *client.Client) *ChannelService {
	return &ChannelService{
		client: client,
	}
}

func (s *ChannelService) GetInfo(ctx context.Context, channelID string) (*entity.Channel, error) {
	result, err := s.client.CallAPI(ctx, &client.Request{
		Method:   http.MethodGet,
		Endpoint: "/channels/get",
		Query: url.Values{
			"channelId": []string{channelID},
		},
	})
	if err != nil {
		return nil, err
	}

	channel := new(entity.Channel)
	if err := json.Unmarshal(result.Response, channel); err != nil {
		return nil, err
	}

	return channel, nil
}
