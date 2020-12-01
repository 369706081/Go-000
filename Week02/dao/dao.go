package dao

import (
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

type MysqlConnectiPool struct {
}

var instance *MysqlConnectiPool
var once sync.Once

var Db *gorm.DB
var err_db error

/**
 * @desc 协程池
 */
var GoPool *ants.Pool

func GetInstance() *MysqlConnectiPool {
	once.Do(func() {
		instance = &MysqlConnectiPool{}
	})
	return instance
}
func (m *MysqlConnectiPool) InitDataPool() (issucc bool) {
	mysqluser := "root"
	mysqlpass := "root"
	mysqlurls := "127.0.0.1:3306"
	mysqldb := "test"
	prefix := "test_"

	dsn := mysqluser + ":" + mysqlpass + "@tcp(" + mysqlurls + ")/" + mysqldb + "?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	Db, err_db = gorm.Open(mysql.Open(dsn), &gorm.Config{

		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   prefix,
		},

		SkipDefaultTransaction: true,
	})
	if err_db != nil {

		return false
	}
	sqlDB, err := Db.DB()
	if err != nil {

		return false
	}
	//设置连接池
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)

	sqlDB.SetConnMaxLifetime(time.Hour)

	return true

}

func (m *MysqlConnectiPool) GetMysqlDB() (db_con *gorm.DB) {
	return Db
}

func Init() {

	GetInstance().InitDataPool()

}

type Test struct {
	Id   int
	Name string
}

func (r *Test) Get() error {

	//sql.ErrNoRows
	tx := Db.First(&r)


	if tx.Error != nil {
		return errors.Wrap(tx.Error, "dao.Get err")
	}
	return nil
}
