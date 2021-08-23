package main

import (
	"os/exec"
	"fmt"
	"flag"
	"strings"
)

func main(){
	var (
		cmd *exec.Cmd
		err error
		command string
		result []byte
	)
	flag.StringVar(&command, "c","", "the command user input")
	flag.Parse()
	
	cmd = exec.Command("/bin/bash", "-c", command)
	if result, err = cmd.CombinedOutput(); err != nil {
		fmt.Printf("exec command '%s' filed, please chack your command.\n", command)
	} else {
		fmt.Printf("result: %s\n", strings.TrimRight(string(result),"\n"))
	}

}
