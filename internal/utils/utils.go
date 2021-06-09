package utils

import (
	"bufio"
	"fmt"
	"os"
)

func BreakExecution() {
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Inputerror: %v", err)
		BreakExecution()
	}
}
