package models

import (
	"strings"
	"time"
	"unicode"
)

type Comment struct {
	ID          uint64    `json:"id"           valid:"required"`
	SenderID    uint64    `json:"sender_id"    valid:"required"`
	RecipientID uint64    `json:"recipient_id" valid:"required"`
	Text        string    `json:"text"         valid:"required, length(1|4000)~Текст должен быть длинной от 1 до 4000 симвволов"` //nolint:nolintlint
	Rating      uint8     `json:"rating"       valid:"required,min=1,max=5"`
	CreatedAt   time.Time `json:"created_at"   valid:"required"`
}

//easyjson:json
type PreComment struct {
	SenderID    uint64 `json:"sender_id"    valid:"required"`
	RecipientID uint64 `json:"recipient_id" valid:"required"`
	Text        string `json:"text"         valid:"required, length(1|4000)~Текст должен быть длинной от 1 до 4000 симвволов"` //nolint:nolintlint
	Rating      uint8  `json:"rating"       valid:"required,range(1|5)"`
}

//easyjson:json
type CommentChanges struct {
	Text   string `json:"text"         valid:"required, length(1|4000)~Текст должен быть длинной от 1 до 4000 симвволов"` //nolint:nolintlint
	Rating uint8  `json:"rating"       valid:"required,range(1|5)"`
}

func (p *PreComment) Trim() {
	p.Text = strings.TrimFunc(p.Text, unicode.IsSpace)
}

func (p *CommentChanges) Trim() {
	p.Text = strings.TrimFunc(p.Text, unicode.IsSpace)
}
