package main

import (
	"os/exec"
	"fmt"
)

func main(){
	var (
		cmd *exec.Cmd
		err error
	)
	// cmd = exec.Command("/bin/bash", "-c", "pwd")
	cmd = exec.Command("/bin/bash", "-c", "aaa")
	if err = cmd.Run(); err != nil{
		fmt.Println("exec command failed: ", err)
	}
}
