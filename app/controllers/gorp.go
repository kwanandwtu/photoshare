package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/db/app"
	"github.com/robfig/revel/samples/booking/app/models"
)

var (
	dbm *gorp.DbMap
)

type GorpPlugin struct {
	rev.EmptyPlugin
}

func (p GorpPlugin) OnAppStart() {
	db.DbPlugin{}.OnAppStart()
	dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}
	dbm.TraceOn("[gorp]", rev.INFO)

	result, _ := dbm.Exec("show tables")
	if rows, _ := result.RowsAFfected(); rows == 0 {
		rev.INFO.Println("Creating database..")
		t := dbm.AddTable(models.Photo{}).SetKeys(true, "PhotoId")
		t.ColMap("Taken").Transient = true
		t.ColMap("Uploaded").Transient = true
		dbm.CreateTables()
	}
}

type GorpController struct {
	*rev.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() rev.Result {
	txn, err := dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() rev.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() rev.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}