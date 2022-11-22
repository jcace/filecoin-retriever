package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "retrieve",
		Usage: "retrieve <cid> <sp id>",
		Action: func(cctx *cli.Context) error {
			cid := cctx.Args().Get(0)
			spid := cctx.Args().Get(1)
			apiInfo := os.Getenv("FULLNODE_API_INFO")
			if cid == "" || spid == "" {
				return cli.Exit("please provide both cid and sp id arguments", 1)
			}
			if apiInfo == "" {
				return cli.Exit("missing FULLNODE_API_INFO env var", 1)
			}

			fc, closer := NewRetriever(cctx.Context, apiInfo)
			defer closer()
			DoRetrieval(cctx.Context, fc, cid, spid)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
