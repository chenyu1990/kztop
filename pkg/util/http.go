package util

import (
	"fmt"
	"net"
)

func HandleHttpError(err error) {
	switch err := err.(type) {
	case net.Error:
		if err.Timeout() {
			fmt.Println(err.Error())
			return
		}
	default:
		panic(err)
	}
}