package table

import (
	"elichika/utils"

	"fmt"
	"os/exec"
)

func readGitChange(table string) string {
	out, err := exec.Command("cmd", "/c", fmt.Sprintf("cd %s && git log -p --unified=0 %s.sql", git, table)).Output()
	utils.CheckErr(err)
	return string(out)
}
