package models

import (
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

//easyjson:json
type commentInFeedJSON struct {
	ID         uint64    `json:"id"           valid:"required"`
	SenderID   uint64    `json:"sender_id"    valid:"required"`
	SenderName string    `json:"sender_name"`
	Avatar     *string   `json:"avatar"     swaggertype:"string"`
	Text       string    `json:"text"         valid:"required, length(1|4000)~Текст должен быть длинной от 1 до 4000 симвволов"` //nolint:nolintlint
	Rating     uint8     `json:"rating"       valid:"required,min=1,max=5"`
	CreatedAt  time.Time `json:"created_at"   valid:"required"`
}

func (c *CommentInFeed) MarshalJSON() ([]byte, error) {
	commentJs := commentInFeedJSON{
		ID:         c.ID,
		SenderID:   c.SenderID,
		SenderName: c.SenderName,
		Avatar:     utils.NullStringToUnsafe(c.Avatar),
		Text:       c.Text,
		Rating:     c.Rating,
		CreatedAt:  c.CreatedAt,
	}

	return commentJs.MarshalJSON()
}

func (c *CommentInFeed) UnmarshalJSON(bytes []byte) error {
	var commentJs commentInFeedJSON

	if err := commentJs.UnmarshalJSON(bytes); err != nil {
		return err
	}

	c.ID = commentJs.ID
	c.SenderID = commentJs.SenderID
	c.SenderName = commentJs.SenderName
	c.Avatar = utils.UnsafeStringToNull(commentJs.Avatar)
	c.Text = commentJs.Text
	c.Rating = commentJs.Rating
	c.CreatedAt = commentJs.CreatedAt

	return nil
}
