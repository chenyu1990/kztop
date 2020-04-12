package main

import (
	"fmt"
	"github.com/rumblefrog/go-a2s"
	"time"
)

func main()  {
	client, err := a2s.NewClient(
		"a.27015.club:27016",
		a2s.TimeoutOption(time.Second * 2), // Setting timeout option. Default is 3 seconds
		// ... Other options
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	info, err := client.QueryInfo()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bots: %v\n", info.Bots)
	fmt.Printf("Game: %v\n", info.Game)
	fmt.Printf("Map: %v\n", info.Map)
	fmt.Printf("MaxPlayers: %v\n", info.MaxPlayers)
	fmt.Printf("Players: %v\n", info.Players)
	fmt.Printf("Protocol: %v\n", info.Protocol)
	fmt.Printf("VAC: %v\n", info.VAC)
	fmt.Printf("Version: %v\n", info.Version)
	fmt.Printf("Visibility: %v\n", info.Visibility)
	fmt.Printf("Port: %v\n", info.ExtendedServerInfo.Port)
}
