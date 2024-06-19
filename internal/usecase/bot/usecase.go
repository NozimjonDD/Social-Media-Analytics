package bot

import (
	"bytes"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"tgstat/internal/entity"
)

var postsText = `Топ 15 публикации в {{.channelName}}:
{{ range $post := .posts }}
[{{$post.Title}}]({{$post.Link}}): {{$post.Views}}
{{end}}
`

type Post struct {
	Title string
	Link  string
	Views int
}

type ChannelRepository interface {
	GetAll(ctx context.Context) ([]*entity.Channel, error)
}

type PostRepository interface {
	GetAllForWeek(ctx context.Context, filter *entity.PostFilter) ([]*entity.Post, error)
}

type UseCase struct {
	bot      *tgbotapi.BotAPI
	chanRepo ChannelRepository
	postRepo PostRepository
}

func New(bot *tgbotapi.BotAPI, chanRepo ChannelRepository, postRepo PostRepository) *UseCase {
	return &UseCase{
		bot:      bot,
		chanRepo: chanRepo,
		postRepo: postRepo,
	}
}

func (uc *UseCase) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := uc.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.Text == "/start" {
				ctx := context.Background()
				channels, err := uc.chanRepo.GetAll(ctx)
				if err != nil {
					continue
				}
				for i := range channels {
					posts, err := uc.postRepo.GetAllForWeek(ctx, &entity.PostFilter{
						Limit:    15,
						ChanelId: strconv.Itoa(channels[i].Id),
						Order:    "views DESC",
					})
					if err != nil {
						log.Println(err)
						continue
					}

					var tgPosts = make([]*Post, 0, len(posts))
					for i := range posts {
						tgPosts = append(tgPosts, &Post{
							Title: posts[i].GetTitle(),
							Link:  posts[i].GetLink(),
							Views: posts[i].Views,
						})
					}

					t, err := template.New("text").Parse(postsText)
					if err != nil {
						log.Println(err)
						continue
					}

					var out bytes.Buffer
					err = t.Execute(&out, map[string]any{
						"channelName": channels[i].Username,
						"posts":       tgPosts,
					})
					if err != nil {
						log.Println(err)
						continue
					}

					s := removeHTMLTags(out.String())

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, s)
					msg.ParseMode = tgbotapi.ModeMarkdown
					msg.DisableWebPagePreview = true
					_, err = uc.bot.Send(msg)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}
	}
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
