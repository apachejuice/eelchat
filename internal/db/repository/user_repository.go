package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"math/rand"
	"sync"
	"time"

	"github.com/apachejuice/eelchat/internal/db/model"
	"github.com/gofrs/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	u        UserRepository
	userOnce sync.Once
)

func NewUserRepository() UserRepository {
	userOnce.Do(func() {
		u = UserRepository{db: connectDB()}
		repoLog.Debug("Connected user repository to database", "table", u.TableName())
	})

	return u
}

type UserRepository struct {
	db *sql.DB
}

func (u UserRepository) One(qms ...qm.QueryMod) (*model.User, error) {
	return basicOne(func(ctx context.Context) (*model.User, error) {
		return model.Users(qms...).One(ctx, u.db)
	})
}

func (u UserRepository) All(qms ...qm.QueryMod) (model.UserSlice, error) {
	return basicAll(func(ctx context.Context) (model.UserSlice, error) {
		return model.Users(qms...).All(ctx, u.db)
	})
}

func (u UserRepository) Update(closure func(entity *model.User), qms ...qm.QueryMod) error {
	return basicUpdate[*model.User, model.UserSlice](u, u.db, closure, qms...)
}

func (u UserRepository) UpdateAll(closure func(entity *model.User), qms ...qm.QueryMod) (int64, error) {
	return basicUpdateAll[*model.User, model.UserSlice](u, u.db, closure, qms...)
}

func (u UserRepository) Insert(entity *model.User) error {
	return doCtx(func(ctx context.Context) error {
		return entity.Insert(ctx, u.db, boil.Infer())
	})
}

func (u UserRepository) Delete(qms ...qm.QueryMod) (*model.User, error) {
	return basicDelete[*model.User, model.UserSlice](u, u.db, qms...)
}

func (u UserRepository) DeleteAll(qms ...qm.QueryMod) (model.UserSlice, error) {
	return basicDeleteAll[*model.User, model.UserSlice](u, u.db, qms...)
}

func (u UserRepository) Count(qms ...qm.QueryMod) (int64, error) {
	return basicCount(u.TableName(), func(ctx context.Context) (int64, error) {
		return model.Users(qms...).Count(ctx, u.db)
	})
}

func (u UserRepository) TableName() string {
	return model.TableNames.User
}

// Create() implementation

func makeUserId() string {
	id, err := uuid.NewV4()
	if err != nil {
		repoLog.Fatal("Error generating primary key", "error", err)
	}

	return base64.RawStdEncoding.EncodeToString(id.Bytes())
}

func makeDiscriminator() string {
	chars := "ABCDEFGHIJKLMOPQRSTUVXYZ1234567890"
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	discrim := ""

	for i := 0; i < 4; i++ {
		discrim += string(chars[r.Intn(len(chars))])
	}

	return discrim
}

func (u UserRepository) Create(username, pwHash, email string) (*model.User, error) {
	dbEmail := null.StringFromPtr(nil)
	if email != "" {
		dbEmail = null.StringFrom(email)
	}

	user := &model.User{
		ID:            makeUserId(),
		Username:      username,
		Discriminator: makeDiscriminator(),
		PasswordHash:  pwHash,
		RegisteredAt:  time.Now().UTC(),
		LastLogin:     null.TimeFromPtr(nil),
		Email:         dbEmail,
	}

	return user, u.Insert(user)
}
