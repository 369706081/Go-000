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
