package helpers

import "fmt"

func LogError(err error) {
	fmt.Println("[ \033[32mERROR\033[0m ]    ", err)
}
