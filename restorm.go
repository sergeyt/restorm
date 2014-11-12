package restorm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sergeyt/hypster"
)

// ForHypster adds HTTP handlers for given collection of models
func ForHypster(app *hypster.AppBuilder, path string, model interface{}) *hypster.AppBuilder {
	if app == nil {
		panic("app is nil")
	}
	if model == nil {
		panic("model is nil")
	}

	if len(path) == 0 {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	typ := typeOf(model)

	app.
		Route(path).
		Post(hypsterPostModel(typ)).
		Get(hypsterGetModels(typ))

	app.
		Route(path + "/{id}").
		Get(hypsterGetModel(typ)).
		Update(hypsterUpdateModel(typ)).
		Delete(hypsterDeleteModel(typ))

	return app
}

// Collection handlers
// -------------------

// POST /{models}
func hypsterPostModel(typ reflect.Type) hypster.Handler {
	return func(ctx *hypster.Context) (interface{}, error) {
		db, err := getdb(ctx)
		if err != nil {
			return nil, err
		}

		model := create(typ, 0)
		ctx.Read(model)

		db.Save(model)

		if err = db.Error; err != nil {
			return nil, err
		}

		// TODO consider to return only model id
		return model, nil
	}
}

// TODO support basic query
// GET /{models}
func hypsterGetModels(typ reflect.Type) hypster.Handler {
	return func(ctx *hypster.Context) (interface{}, error) {
		db, err := getdb(ctx)
		if err != nil {
			return nil, err
		}

		models := reflect.New(reflect.SliceOf(typ))
		db.Find(models)

		if err = db.Error; err != nil {
			return nil, err
		}

		return models, nil
	}
}

// Document handlers
// -----------------

// GET /{models}/{id}
func hypsterGetModel(typ reflect.Type) hypster.Handler {
	return func(ctx *hypster.Context) (interface{}, error) {
		id, err := strconv.ParseInt(ctx.Vars["id"], 10, 64)
		if err != nil {
			return nil, err
		}

		db, err := getdb(ctx)
		if err != nil {
			return nil, err
		}

		model := create(typ, 0)
		db.First(model, id)

		if err = db.Error; err != nil {
			return nil, err
		}

		return model, nil
	}
}

// UPDATE /{models}/{id}
func hypsterUpdateModel(typ reflect.Type) hypster.Handler {
	return func(ctx *hypster.Context) (interface{}, error) {
		id, err := strconv.ParseInt(ctx.Vars["id"], 10, 64)
		if err != nil {
			return nil, err
		}

		db, err := getdb(ctx)
		if err != nil {
			return nil, err
		}

		model := create(typ, id)
		ctx.Read(model)

		// TODO update only fields that are come in JSON input
		db.Model(model).Updates(model)

		if err = db.Error; err != nil {
			return nil, err
		}

		return true, nil
	}
}

// DELETE /{models}/{id}
func hypsterDeleteModel(typ reflect.Type) hypster.Handler {
	return func(ctx *hypster.Context) (interface{}, error) {
		id, err := strconv.ParseInt(ctx.Vars["id"], 10, 64)
		if err != nil {
			return nil, err
		}

		db, err := getdb(ctx)
		if err != nil {
			return nil, err
		}

		model := create(typ, id)
		db.Delete(model)

		if err = db.Error; err != nil {
			return nil, err
		}

		return true, nil
	}
}

// Helpers

func getdb(ctx *hypster.Context) (*gorm.DB, error) {
	db := ctx.GetService("db").(*gorm.DB)
	if db == nil {
		return nil, fmt.Errorf(`no "db" service of type gorm.DB`)
	}
	return db, nil
}
