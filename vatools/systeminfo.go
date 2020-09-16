package vatools

import (
	"fmt"
	"os/exec"
)

func GetCupInfo() string {
	obExec := exec.Command("wmic")
	out, err := obExec.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	fmt.Println(string(out))
	return string(out)
}
