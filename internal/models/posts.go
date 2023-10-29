package models

type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}

type Product struct {
	ID              uint64 `json:"id"          valid:"required"`
	AuthorID        uint64 `json:"author"      valid:"required"`
	Title           string `json:"title"       valid:"required"`
	Image           Image  `json:"image"       valid:"required"`
	Description     string `json:"description" valid:"required"`
	Price           uint   `json:"price"       valid:"required"`
	SafeTransaction bool   `json:"safe"        valid:"required"`
	Delivery        bool   `json:"delivery"    valid:"required"`
	City            string `json:"city"        valid:"required"`
}

type PreProduct struct {
	AuthorID        uint64 `json:"author"      valid:"required"`
	Title           string `json:"title"       valid:"required"`
	Image           Image  `json:"image"       valid:"required"`
	Description     string `json:"description" valid:"required"`
	Price           uint   `json:"price"       valid:"required"`
	SafeTransaction bool   `json:"safe"        valid:"required"`
	Delivery        bool   `json:"delivery"    valid:"required"`
	City            string `json:"city"        valid:"required"`
}

type ProductInFeed struct {
	ID              uint64 `json:"id"       valid:"required"`
	Title           string `json:"title"    valid:"required"`
	Image           Image  `json:"image"    valid:"required"`
	Price           uint   `json:"price"    valid:"required"`
	SafeTransaction bool   `json:"safe"     valid:"required"`
	Delivery        bool   `json:"delivery" valid:"required"`
	City            string `json:"city"     valid:"required"`
}
