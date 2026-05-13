package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	if prompt != "" {
		fmt.Print(prompt)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}
