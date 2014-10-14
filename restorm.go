package restorm

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/sergeyt/hypster"
)

// RegisterHandlers adds HTTP handlers for given collection of models
func RegisterHandlers(app *hypster.AppBuilder, collection string, typ reflect.Type) *hypster.AppBuilder {
	app.
		Route("/" + collection).
		Post(postModel(typ)).
		Get(getModels(typ))

	app.
		Route("/" + collection + "/{id}").
		Get(getModel(typ)).
		Update(updateModel(typ)).
		Delete(deleteModel(typ))

	return app
}

// Collection handlers
// -------------------

// POST /{models}
func postModel(typ reflect.Type) hypster.Handler {
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
func getModels(typ reflect.Type) hypster.Handler {
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
func getModel(typ reflect.Type) hypster.Handler {
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
func updateModel(typ reflect.Type) hypster.Handler {
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
func deleteModel(typ reflect.Type) hypster.Handler {
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

func create(typ reflect.Type, id int64) interface{} {
	val := reflect.New(typ)
	// set id
	return val
}

func getdb(ctx *hypster.Context) (*gorm.DB, error) {
	db := ctx.GetService("db").(*gorm.DB)
	if db == nil {
		return nil, fmt.Errorf(`no "db" service of type gorm.DB`)
	}
	return db, nil
}
