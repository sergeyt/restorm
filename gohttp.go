package restorm

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gohttp/app"
	"github.com/jinzhu/gorm"
)

// TODO struct to pass options

// RegisterHandlers adds HTTP handlers for given model collection.
func RegisterHandlers(a *app.App, db *gorm.DB, path string, model interface{}) *app.App {
	if a == nil {
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

	a.Post(path, postModel(db, typ))
	a.Get(path, getModels(db, typ))

	var pat = path + "/:id"
	a.Get(pat, getModel(db, typ))
	a.Put(pat, updateModel(db, typ))
	a.Del(pat, deleteModel(db, typ))

	return a
}

func typeOf(value interface{}) reflect.Type {
	var t reflect.Type
	switch value.(type) {
	case reflect.Type:
		t = value.(reflect.Type)
	default:
		t = reflect.TypeOf(value)
	}
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

// Collection handlers
// -------------------

// POST /{models}
func postModel(db *gorm.DB, typ reflect.Type) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		model := create(typ, 0)
		readJSON(r, model)

		db.Save(model)

		if err := db.Error; err != nil {
			writeError(w, err)
			return
		}

		// TODO consider to return only model id
		writeJSON(w, model)
	})
}

// TODO support basic query
// GET /{models}
func getModels(db *gorm.DB, typ reflect.Type) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		models := reflect.New(reflect.SliceOf(typ))
		db.Find(models)

		if err := db.Error; err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, models)
	})
}

// Document handlers
// -----------------

// GET /{models}/{id}
func getModel(db *gorm.DB, typ reflect.Type) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			writeError(w, err)
			return
		}

		model := create(typ, 0)
		db.First(model, id)

		if err = db.Error; err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, model)
	})
}

// UPDATE /{models}/{id}
func updateModel(db *gorm.DB, typ reflect.Type) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			writeError(w, err)
			return
		}

		model := create(typ, id)
		readJSON(r, model)

		// TODO update only fields that are come in JSON input
		db.Model(model).Updates(model)

		if err = db.Error; err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, true)
	})
}

// DELETE /{models}/{id}
func deleteModel(db *gorm.DB, typ reflect.Type) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			writeError(w, err)
			return
		}

		model := create(typ, id)
		db.Delete(model)

		if err = db.Error; err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, true)
	})
}

// Helpers

func getID(r *http.Request) (int64, error) {
	var s = r.URL.Query().Get(":id")
	return strconv.ParseInt(s, 10, 64)
}

func create(typ reflect.Type, id int64) interface{} {
	val := reflect.New(typ)
	if id != 0 {
		// TODO fix set ID field
	}
	return val
}

func readJSON(r *http.Request, out interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(out)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	bytes, _ := json.Marshal(v)
	w.Write(bytes)
}

// TODO add error type
type errorPayload struct {
	error string
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, errorPayload{err.Error()})
}
