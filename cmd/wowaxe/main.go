package main

import (
	"os"

	"github.com/coline-carle/battleaxe/battle"
	"github.com/coline-carle/battleaxe/cli"
)

func main() {
	cli.Run(battle.WoW, os.Args)
}
