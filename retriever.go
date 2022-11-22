package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	filclient "github.com/application-research/filclient-unstable"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"
	flatfs "github.com/ipfs/go-ds-flatfs"
	leveldb "github.com/ipfs/go-ds-leveldb"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/libp2p/go-libp2p"
)

func NewRetriever(ctx context.Context, lotusString string) *filclient.Client {
	// Instantiate libp2p node
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer node.Close()

	// print the node's listening addresses
	fmt.Println("Listen addresses:", node.Addrs())

	// Connect to lotus node - this can be a lite node
	api, closer, err := LotusConnection(lotusString)
	defer closer()
	if err != nil {
		log.Fatalf("could not connect to lotus")
	}

	// Get default wallet of node
	addr, err := api.WalletDefaultAddress(ctx)
	if err != nil {
		log.Fatalf("Could not get wallet address: %v", err)
	}

	// Initialize Blockstore
	parseShardFunc, err := flatfs.ParseShardFunc("/repo/flatfs/shard/v1/next-to-last/3")
	if err != nil {
		log.Fatalf("Blockstore parse shard func failed: %v", err)
	}
	ds, err := flatfs.CreateOrOpen(filepath.Join("/tmp", "blockstore"), parseShardFunc, false)
	if err != nil {
		log.Fatalf("Could not initialize blockstore: %v", err)
	}

	bs := blockstore.NewBlockstoreNoPrefix(ds)

	ldb, err := leveldb.NewDatastore("", nil)
	if err != nil {
		log.Fatalf("Could not initialize datastore: %v", err)
	}

	fc, err := filclient.New(ctx, node, api, addr, bs, ldb)
	if err != nil {
		log.Fatalf("Could not initialize filclient: %v", err)
	}

	fmt.Printf("fc: %v\n", fc)

	return fc
}

func DoRetrieval(ctx context.Context, fc *filclient.Client, c string, mid string) {
	cid, err := cid.Parse(c)
	if err != nil {
		log.Fatalf("parsing cid failed: %s", err)
	}

	minerAddr, err := address.NewFromString(mid)
	if err != nil {
		log.Fatalf("parsing miner addr failed: %s", err)
	}

	miner := fc.MinerByAddress(minerAddr)

	log.Println("starting retrieval...")

	transfer, err := miner.StartRetrievalTransfer(ctx, cid)
	if err != nil {
		log.Fatalf("initiating retrieval failed: %s", err)
	}

	<-transfer.Done()
}
