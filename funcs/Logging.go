package netcat

import (
	"fmt"
	"os"
)

func Logging(logs string) {
	//open file and log messages, it will create if not existing already
	file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	//append messages to the file
	_, err = file.WriteString(logs)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
