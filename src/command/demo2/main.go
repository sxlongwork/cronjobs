package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main(){
	var (
		cmd *exec.Cmd
		result []byte
		err error
	)
	cmd = exec.Command("/bin/bash", "-c", "pwd")
	if result, err = cmd.CombinedOutput(); err != nil {
		fmt.Println("exec command failed: ", err)
	} else {
		fmt.Println("result:", strings.TrimRight(string(result),"\n"))
		fmt.Println("exec completed.")
	}
}
