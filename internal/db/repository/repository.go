package repository

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// The base interface for objects that retrieve database things.
type Repository[T any, ST ~[]T] interface {
	// Returns one object, found with the given query mods.
	One(qms ...qm.QueryMod) (T, error)
	// Returns many objects, found with the given query mods.
	All(qms ...qm.QueryMod) (ST, error)
	// Returns the count of objects locatable with the given query mods.
	Count(qms ...qm.QueryMod) (int64, error)

	// Updates a row with the given query mods. Use the closure to mutate the object, do not perform database operations.
	Update(closure func(entity T), qms ...qm.QueryMod) error
	// Updates all rows with the given query mods. Use the closure to mutate the object, do not perform database operations.
	UpdateAll(closure func(entity T), qms ...qm.QueryMod) (int64, error)

	// Inserts the given object in the database.
	// Usually a repository also has a Create() method, which excludes certain auto-assigned fields.
	Insert(entity T) error

	// Deletes a row with the given query mods.
	Delete(qms ...qm.QueryMod) (T, error)
	// Deletes all rows with the given query mods
	DeleteAll(qms ...qm.QueryMod) (ST, error)

	// The table name in the database
	TableName() string
}

type entityLike interface {
	Update(ctx context.Context, ctxe boil.ContextExecutor, cols boil.Columns) (int64, error)
	Delete(ctx context.Context, ctxe boil.ContextExecutor) (int64, error)
}

func basicOne[T any](closure func(ctx context.Context) (T, error)) (T, error) {
	var entity T
	err := doCtx(func(ctx context.Context) error {
		e, err := closure(ctx)
		entity = e
		return err
	})

	return entity, err
}

func basicAll[ST any](closure func(ctx context.Context) (ST, error)) (ST, error) {
	var entities ST
	err := doCtx(func(ctx context.Context) error {
		e, err := closure(ctx)
		entities = e
		return err
	})

	return entities, err
}

func basicUpdate[T any, ST ~[]T](r Repository[T, ST], db *sql.DB, update func(entity T), qms ...qm.QueryMod) error {
	entity, err := r.One(qms...)
	if err != nil {
		return err
	}

	update(entity)
	return doCtx(func(ctx context.Context) error {
		if _, err := any(entity).(entityLike).Update(ctx, db, boil.Infer()); err != nil {
			return err
		}

		return nil
	})
}

func basicUpdateAll[T any, ST ~[]T](r Repository[T, ST], db *sql.DB, update func(entity T), qms ...qm.QueryMod) (int64, error) {
	entities, err := r.All(qms...)
	if err != nil {
		return 0, nil
	}

	return int64(len(entities)), doCtx(func(ctx context.Context) error {
		for _, entity := range entities {
			update(entity)
			if _, err := any(entity).(entityLike).Update(ctx, db, boil.Infer()); err != nil {
				return err
			}
		}

		return nil
	})
}

func basicDelete[T any, ST ~[]T](r Repository[T, ST], db *sql.DB, qms ...qm.QueryMod) (T, error) {
	entity, err := r.One(qms...)
	if err != nil {
		return any(nil).(T), err
	}

	return entity, doCtx(func(ctx context.Context) error {
		_, err := any(entity).(entityLike).Delete(ctx, db)
		return err
	})
}

func basicDeleteAll[T any, ST ~[]T](r Repository[T, ST], db *sql.DB, qms ...qm.QueryMod) (ST, error) {
	entities, err := r.All(qms...)
	if err != nil {
		return nil, err
	}

	return entities, doCtx(func(ctx context.Context) error {
		for _, entity := range entities {
			if _, err := any(entity).(entityLike).Delete(ctx, db); err != nil {
				return err
			}
		}

		return nil
	})
}

func basicCount(table string, closure func(ctx context.Context) (int64, error)) (int64, error) {
	var count int64
	err := doCtx(func(ctx context.Context) error {
		c, err := closure(ctx)
		count = c
		return err
	})

	repoLog.Debug("Retrieved entity count", "table", table, "count", count) // helps with debugging specific insertion errors
	return count, err
}
