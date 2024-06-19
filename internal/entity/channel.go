package entity

import "github.com/uptrace/bun"

type Channel struct {
	bun.BaseModel      `bun:"table:channels"`
	Id                 int           `json:"id" bun:"id,pk"`
	TgId               int           `json:"tg_id"`
	Link               string        `json:"link"`
	PeerType           string        `json:"peer_type"`
	Username           string        `json:"username"`
	ActiveUsernames    []string      `json:"active_usernames"`
	Title              string        `json:"title"`
	About              string        `json:"about"`
	Category           string        `json:"category"`
	Country            string        `json:"country"`
	Language           string        `json:"language"`
	Image100           string        `json:"image100"`
	Image640           string        `json:"image640"`
	ParticipantsCount  int           `json:"participants_count"`
	TgstatRestrictions []interface{} `json:"tgstat_restrictions"`
	LastUpdated        int64         `json:"-"`
}
