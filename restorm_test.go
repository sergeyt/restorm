package restorm

import (
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/franela/go-supertest"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/gomega"
	. "github.com/sergeyt/goblin"
	"github.com/sergeyt/hypster"
)

type User struct {
	ID        int64
	Name      string `sql:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Test(t *testing.T) {
	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	describe := g.Describe
	it := g.It

	var open = func() *gorm.DB {
		db, err := gorm.Open("sqlite3", "/tmp/restorm_test")
		if err != nil {
			g.Fail("cannot init database")
		}
		db.DB().Ping()
		// clean db
		db.DropTable(&User{})
		return &db
	}

	describe("with restorm", func() {
		it("I should easily POST, GET models", func(done Done) {
			db := open()
			services := make(map[string]interface{})
			services["db"] = db
			app := hypster.NewApp(services)
			RegisterHandlers(app, "users", User{})

			ts := httptest.NewServer(app)
			defer ts.Close()

			getUsers := func() {
				// TODO check body
				NewRequest(ts.URL).
					Get("/users").
					Expect(200, done)
			}

			NewRequest(ts.URL).
				Post("/users").
				Send(&User{Name: "test"}).
				Expect(200, getUsers)
		})
	})
}
