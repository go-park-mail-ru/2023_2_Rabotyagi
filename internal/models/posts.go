package models

type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}

type Post struct {
	ID              uint64 `json:"id"`
	AuthorID        uint64 `json:"author"`
	Title           string `json:"title"`
	Image           Image  `jsom:"image"`
	Description     string `json:"description"`
	Price           int    `json:"price"`
	SafeTransaction bool   `json:"safe"`
	Delivery        bool   `json:"delivery"`
	City            string `json:"city"`
}

type PrePost struct {
	AuthorID        uint64 `json:"author"`
	Title           string `json:"title"`
	Image           Image  `jsom:"image"`
	Description     string `json:"description"`
	Price           int    `json:"price"`
	SafeTransaction bool   `json:"safe"`
	Delivery        bool   `json:"delivery"`
	City            string `json:"city"`
}

type PostInFeed struct {
	ID              uint64 `json:"id"`
	Title           string `json:"title"`
	Image           Image  `json:"image"`
	Price           int    `json:"price"`
	SafeTransaction bool   `json:"safe"`
	Delivery        bool   `json:"delivery"`
	City            string `json:"city"`
}
