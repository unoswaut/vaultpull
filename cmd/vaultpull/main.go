package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/sync"
	"github.com/example/vaultpull/internal/vault"
)

func main() {
	cfgPath := flag.String("config", "vaultpull.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		log.Fatalf("vault client: %v", err)
	}

	s := sync.New(client, cfg)
	results := s.Run()

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "FAIL %s: %v\n", r.Path, r.Err)
		} else {
			fmt.Printf("OK   %s -> %s\n", r.Path, r.OutFile)
		}
	}

	if sync.HasErrors(results) {
		os.Exit(1)
	}
}
