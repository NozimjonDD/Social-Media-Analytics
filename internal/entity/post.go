package entity

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"strings"
)

type Posts struct {
	Count      int      `json:"count"`
	TotalCount int      `json:"total_count"`
	Channel    *Channel `json:"channel"`
	Items      []*Post  `json:"items"`
}

type Post struct {
	bun.BaseModel `bun:"table:posts"`
	PostID        uuid.UUID `json:"post_id" bun:"post_id,pk"`
	Id            int       `json:"id" bun:"id"`
	Title         string    `json:"title" bun:"-"`
	Date          int       `json:"date" bun:"date"`
	Views         int       `json:"views" bun:"views"`
	Link          string    `json:"link" bun:"link"`
	ChannelId     int       `json:"channel_id" bun:"channel_id"`
	IsDeleted     int       `json:"is_deleted" bun:"is_deleted"`
	Text          string    `json:"text" bun:"text"`
	Media         struct {
		MediaType string `json:"media_type" `
		MimeType  string `json:"mime_type" `
		Size      int    `json:"size" `
	} `json:"media" bun:"media"`
}

func (p *Post) MarshalJSON() ([]byte, error) {
	s := strings.Split(p.Text, "\n")
	if len(s) >= 2 {
		p.Title = s[0]
	}
	p.Title = strings.Replace(p.Title, "<b>", "", -1)
	return json.Marshal(*p)
}

func (p *Post) GetTitle() string {
	s := strings.Split(p.Text, "\n")
	if len(s) >= 1 {
		p.Title = s[0]
	}
	//fmt.Println(s)
	if len(p.Title) == 0 {
		//fmt.Println("xaxaxa")
		s := strings.Split(p.Text, "</b>")
		fmt.Println(p.Text)
		if len(s) >= 1 {
			p.Title = s[0]
			p.Title += "..."
		}
	}

	//var htmlEscaper = strings.NewReplacer(
	//	`<br>`, "",
	//	`</br>`, "",
	//	`<`, "\\<",
	//	`>`, "\\>",
	//	`.`, "\\.",
	//	`-`, "\\-",
	//)

	p.Title = strings.Replace(p.Title, "<b>", "", -1)
	p.Title = strings.Replace(p.Title, "</b>", "", -1)
	p.Title = strings.ReplaceAll(p.Title, ".", "\\.")
	p.Title = strings.ReplaceAll(p.Title, "-", "\\-")
	p.Title = strings.ReplaceAll(p.Title, "<", "\\<")
	p.Title = strings.ReplaceAll(p.Title, ">", "\\>")
	//p.Title = htmlEscaper.Replace(p.Title)
	return p.Title
}

func (p *Post) GetLink() string {
	return fmt.Sprintf("https://%s", p.Link)
}

type PostServiceGetAllOptions struct {
	ChannelID string
	Offset    int
	Limit     int
	StartTime int64
	EndTime   int64
}

type PostFilter struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	ChanelId string `form:"channel_id"`
	Order    string `form:"order"`
	ByWord   string `form:"by_word"`
}

func (p *Post) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		p.PostID = uuid.New()
	}
	return nil
}

type FrequencyWords struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}
