package post

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"tgstat/internal/entity"
	"time"
)

const (
	Limit                            = "50"
	LimitDays          int           = 60
	PostUpdateInterval time.Duration = time.Hour * 24
)

type Repository interface {
	GetAll(ctx context.Context, filter *entity.PostFilter) ([]*entity.Post, error)
	GetAllForWeek(ctx context.Context, filter *entity.PostFilter) ([]*entity.Post, error)
	Create(ctx context.Context, posts []*entity.Post) error
	GetFrequencyWords(ctx context.Context, channelID int) ([]*entity.FrequencyWords, error)
}

type ChannelRepository interface {
	GetAll(ctx context.Context) ([]*entity.Channel, error)
	//GetLastUpdatedDate(ctx context.Context, channelID int) (int64, error)
}

type Service interface {
	GetAll(ctx context.Context, v url.Values) (*entity.Posts, error)
}

type UseCase struct {
	channelRepo ChannelRepository
	postRepo    Repository
	postService Service
}

func NewUseCase(channelRepo ChannelRepository, postRepo Repository, postService Service) *UseCase {
	return &UseCase{
		channelRepo: channelRepo,
		postRepo:    postRepo,
		postService: postService,
	}
}

func (uc *UseCase) SyncAllPosts(ctx context.Context, channelID string) error {
	v := url.Values{}
	v.Set("channelId", channelID)
	v.Set("limit", Limit)
	v.Set("offset", "0")

	var days int

	for days != LimitDays {
		startTime, _ := time.Parse("2006-01-02", time.Now().AddDate(0, 0, -days).Format("2006-01-02"))
		v.Set("startTime", strconv.FormatInt(startTime.Unix(), 10))
		endTime := startTime.Add(time.Second * 86_399)
		v.Set("endTime", strconv.FormatInt(endTime.Unix(), 10))

		days++

		result, err := uc.postService.GetAll(ctx, v)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := uc.clearAndCreatePosts(ctx, result.Items); err != nil {
			fmt.Println(err)
			continue
		}

		//if err := uc.postRepo.Create(ctx, result.Items); err != nil {
		//	fmt.Println(err)
		//	continue
		//}
	}
	return nil
}

func (uc *UseCase) RunSyncWorker(ctx context.Context) {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Day().At("23:59").Do(func() {
		channels, err := uc.channelRepo.GetAll(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, channel := range channels {
			t := time.Unix(channel.LastUpdated, 0)
			now := time.Now()

			v := url.Values{}
			v.Set("channelId", strconv.Itoa(channel.Id))
			v.Set("limit", Limit)

			var (
				count, counter, offset int
			)

			counter = -1

			for count != counter {
				v.Set("startTime", strconv.FormatInt(t.Unix(), 10))
				v.Set("endTime", strconv.FormatInt(now.Unix(), 10))
				v.Set("offset", strconv.Itoa(offset))

				result, err := uc.postService.GetAll(ctx, v)
				if err == nil {
					fmt.Println(err)
					continue
				}

				if err := uc.clearAndCreatePosts(ctx, result.Items); err != nil {
					fmt.Println(err)
					continue
				}

				counter += len(result.Items)
				count = result.Count
				offset += len(result.Items)
			}
		}
	})
	fmt.Println(err)
}

func (uc *UseCase) GetAllPosts(ctx context.Context, filter *entity.PostFilter) ([]*entity.Post, error) {
	return uc.postRepo.GetAll(ctx, filter)
}

func (uc *UseCase) GetAllForWeek(ctx context.Context, filter *entity.PostFilter) ([]*entity.Post, error) {
	return uc.postRepo.GetAllForWeek(ctx, filter)
}

func (uc *UseCase) GetFrequencyWords(ctx context.Context, channelID int) ([]*entity.FrequencyWords, error) {
	return uc.postRepo.GetFrequencyWords(ctx, channelID)
}

func (uc *UseCase) clearAndCreatePosts(ctx context.Context, posts []*entity.Post) error {
	for i := range posts {
		//indexStart := strings.Index(posts[i].Text, "<a")
		//indexEnd := strings.Index(posts[i].Text, "</a>")
		//if indexStart == 0 || indexEnd == 0 {
		//	continue
		//}
		//text := posts[i].Text[indexStart:]
		//text += posts[i].Text[:indexEnd+4]
		posts[i].Text = removeHTMLTags(posts[i].Text)
		//fmt.Println("-------------------------------------")
		//fmt.Println(posts[i].Text)
		//fmt.Println("-------------------------------------")
		//time.Sleep(time.Second * 3)
	}

	return uc.postRepo.Create(ctx, posts)
}

func removeHTMLTags(in string) string {
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}
