package main

import (
	"fmt"
	connect4 "rotating-gravity-connect-4"
)

// import "connect4"

func main() {

	connect4.Game.Insert(1, 1)
	connect4.Game.Insert(1, 1)
	x := connect4.Game.Insert(2, 1)
	connect4.Game.Fall()
	if !x {
		fmt.Println("x")
	}
	fmt.Println(connect4.Game)
	connect4.Main()
}
