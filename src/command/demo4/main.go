package main

import (
	"os/exec"
	"fmt"
	"strings"
	"time"
	"flag"
	"context"
)

func main(){
	var (
		cmd *exec.Cmd
		err error
		command string
		result []byte
		ctx context.Context
		cancelFunc context.CancelFunc
		strchan chan string = make(chan string)
	)
	flag.StringVar(&command, "c", "", "the command user input")
	flag.Parse()

	ctx, cancelFunc = context.WithCancel(context.TODO())

	cmd = exec.CommandContext(ctx, "/bin/bash", "-c", command)
	go func(){
		if result, err = cmd.CombinedOutput(); err != nil {
			//fmt.Println("ERROR:", err)
			strchan <- err.Error()
		} else {
			//fmt.Println(strings.TrimRight(string(result), "\n"))
			strchan <- strings.TrimRight(string(result), "\n")
		}
	}()

	time.Sleep(1 * time.Second)
	//cancelFunc = cancelFunc
	cancelFunc()
	//time.Sleep(5 * time.Second)
	select {
		case v:= <-strchan:
			fmt.Println(v)
	}
}
