package main

import (
	"elichika/exec/db_recovery/table"

	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: db_recovery <path/to/database/git/repository>/<database_name> <optional path/to/output/file>")
		return
	}
	output := table.Run(os.Args[1])
	if len(os.Args) >= 3 {
		os.WriteFile(os.Args[2], []byte(output), 0777)
	} else {
		fmt.Println(output)
	}

}
