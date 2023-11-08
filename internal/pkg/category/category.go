package category

import "context"

type ICategoryStorage interface {
	GetFullCategories(ctx context.Context)
}
