package main

import (
	"github.com/gorhill/cronexpr"
	"time"
	"fmt"
)

func main(){
	var (
		exp *cronexpr.Expression
		err error
		cur time.Time
	)
	if exp, err = cronexpr.Parse("0 * * * * * *"); err != nil {
		fmt.Println(err)
	} else {
		cur = time.Now()
		fmt.Println(cur)
		fmt.Println(exp.Next(cur))
		//fmt.Println(*exp)
	}
	
}
