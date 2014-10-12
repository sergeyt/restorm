package restorm

import (
	"github.com/jinzhu/gorm"
	. "github.com/sergeyt/hypster"
	. "strconv"
)

type (
	// IModel defines interface for models
	IModel interface {
		// NewArray creates array of this model
		NewArray() interface{}
	}

	// ModelFactory creates new model instance with given primary key
	ModelFactory func(id int64) IModel
)

// RegisterHandlers adds HTTP handlers for given collection of models
func RegisterHandlers(app *AppBuilder, collection string, factory ModelFactory) *AppBuilder {
	r := app.Route(collection)
	r.Post(post_model(factory))
	r.Get(get_models(factory))

	r = app.Route(collection + "/{id}")
	r.Get(get_model(factory))
	r.Update(update_model(factory))
	r.Delete(delete_model(factory))

	return app
}

// Collection handlers
// -------------------

// POST /{models}
func post_model(factory ModelFactory) Handler {
	return func(ctx *Context) (res interface{}, err error) {
		model := factory(0)
		ctx.Read(model)

		db := get_db(ctx)
		db.Save(model)

		if err = db.Error; err != nil {
			return
		}

		// TODO consider to return only issue id
		res = model
		return
	}
}

// TODO support basic query
// GET /{models}
func get_models(factory ModelFactory) Handler {
	return func(ctx *Context) (res interface{}, err error) {
		model := factory(0)
		models := model.NewArray()
		db := get_db(ctx)
		db.Find(models)

		if err = db.Error; err != nil {
			return
		}

		res = models
		return
	}
}

// Document handlers
// -----------------

// GET /{models}/{id}
func get_model(factory ModelFactory) Handler {
	return func(ctx *Context) (res interface{}, err error) {
		id, err := ParseInt(ctx.Vars["id"], 10, 64)
		if err != nil {
			return
		}

		model := factory(0)
		db := get_db(ctx)
		db.First(model, id)

		if err = db.Error; err != nil {
			return
		}

		res = model
		return
	}
}

// UPDATE /{models}/{id}
func update_model(factory ModelFactory) Handler {
	return func(ctx *Context) (res interface{}, err error) {
		id, err := ParseInt(ctx.Vars["id"], 10, 64)
		if err != nil {
			return
		}

		model := factory(id)
		ctx.Read(model)

		// TODO update only fields that are come in JSON input
		db := get_db(ctx)
		db.Model(model).Updates(model)

		if err = db.Error; err != nil {
			return
		}

		res = true
		return
	}
}

// DELETE /{models}/{id}
func delete_model(factory ModelFactory) Handler {
	return func(ctx *Context) (res interface{}, err error) {
		id, err := ParseInt(ctx.Vars["id"], 10, 64)
		if err != nil {
			return
		}

		model := factory(id)
		db := get_db(ctx)
		db.Delete(model)

		if err = db.Error; err != nil {
			return
		}

		res = true
		return
	}
}

// Helpers

func get_db(ctx *Context) *gorm.DB {
	return ctx.GetService("db").(*gorm.DB)
}
