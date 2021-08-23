package main

import (
	"github.com/gorhill/cronexpr"
	"fmt"
	"time"
)
type CronJob struct{
	exp *cronexpr.Expression
	next time.Time
}

func main(){
	var (
		jobs map[string]*CronJob = make(map[string]*CronJob)
		expr *cronexpr.Expression
		err error
		cur time.Time = time.Now()
	)
	if expr, err = cronexpr.Parse("*/10 * * * * * *"); err!= nil {
		fmt.Println(err)
	} else {
		job := &CronJob{expr, expr.Next(cur)}
		jobs["job1"] = job
	}
	if expr, err = cronexpr.Parse("*/20 * * * * * *"); err!= nil {
                fmt.Println(err)
        } else {
                job := &CronJob{expr, expr.Next(cur)}
                jobs["job2"] = job
        }
	go func(){
		for {
			cur = time.Now()
			for name, job := range jobs{
				if job.next.Before(cur) || job.next.Equal(cur) {
					go func(name string){
						fmt.Printf("begin to exec %s.\n", name)
					}(name)
					job.next = job.exp.Next(cur)
					fmt.Printf("%s next exec time: %v\n", name, job.next)
				}
			}
			//time.Sleep(100 * time.Millisecond)
			select {
				case <- time.NewTimer(100 * time.Millisecond).C:
			}
		}
	}()
	time.Sleep(1 * time.Minute)
}
