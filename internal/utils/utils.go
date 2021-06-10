package utils

import (
	"bufio"
	"fmt"
	"os"
)

func BreakExecution() string {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Inputerror: %v", err)
		return BreakExecution()
	}
	return str
}
