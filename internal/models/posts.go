package models

type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}

type Post struct {
	ID              uint64 `json:"id"          valid:"required"`
	AuthorID        uint64 `json:"author"      valid:"required"`
	Title           string `json:"title"       valid:"required"`
	Image           Image  `jsom:"image"       valid:"required"`
	Description     string `json:"description" valid:"required"`
	Price           uint   `json:"price"       valid:"required"`
	SafeTransaction bool   `json:"safe"        valid:"required"`
	Delivery        bool   `json:"delivery"    valid:"required"`
	City            string `json:"city"        valid:"required"`
}

type PrePost struct {
	AuthorID        uint64 `json:"author"      valid:"required"`
	Title           string `json:"title"       valid:"required"`
	Image           Image  `jsom:"image"       valid:"required"`
	Description     string `json:"description" valid:"required"`
	Price           uint   `json:"price"       valid:"required"`
	SafeTransaction bool   `json:"safe"        valid:"required"`
	Delivery        bool   `json:"delivery"    valid:"required"`
	City            string `json:"city"        valid:"required"`
}

type PostInFeed struct {
	ID              uint64 `json:"id"       valid:"required"`
	Title           string `json:"title"    valid:"required"`
	Image           Image  `jsom:"image"    valid:"required"`
	Price           uint   `json:"price"    valid:"required"`
	SafeTransaction bool   `json:"safe"     valid:"required"`
	Delivery        bool   `json:"delivery" valid:"required"`
	City            string `json:"city"     valid:"required"`
}
