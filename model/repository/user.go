package repository

import (
	"context"
	"github.com/PhantomWolf/recreationroom-auth/model/model"
)

type User interface {
	Create(ctx context.Context, user *model.User)
	Update(ctx context.Context, user *model.User)
	Delete(ctx context.Context, id int64)
	Find(ctx context.Context, id int64)
}
