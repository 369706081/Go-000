## 作业
我们在数据库操作的时候，比如 `dao` 层中当遇到一个 `sql.ErrNoRows` 的时候，是否应该 `Wrap` 这个 `error`，抛给上层。为什么？应该怎么做请写出代码
## 回答
在dao层遇到sql.ErrNoRows`时都应该进行wrap封装，在由业务层进行判断是否要中断业务，还是进行降级服务
## 代码实现：
```go
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
```
##业务层
```go 
package main

import (
	"Go-000/Week02/dao"
	"errors"
	"fmt"
	gerrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

func main()  {
	dao.Init()
	test :=dao.Test{
		Id:   2,

	}
	err :=test.Get()
	if err != nil && errors.Is(err,gorm.ErrRecordNotFound) {
			fmt.Println(gerrors.Cause(err))
	//to do	根据业务场景是否要吞掉err
	}else{
		fmt.Printf("stact error:\n%+v\n",err)
		return
	}
	fmt.Println("hello word")
}
