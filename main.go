package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
)

func main() {
	cfg := parseFlags()
	ctx := context.Background()

	if cfg.isDebug() {
		fmt.Printf(
			"[DEBUG] flags - years: %d, months: %d, days: %d, max: %d, min: %d, amount: %d\n",
			cfg.Years,
			cfg.Months,
			cfg.Days,
			cfg.Max,
			cfg.Min,
			cfg.Amount,
		)
	}

	validateFlags(&cfg)

	out, err := Pull(cfg)
	if err != nil {
		panic(fmt.Errorf("[ERROR] couldn't pull from remote repository: %w", err))
	}

	if !bytes.Contains(out, []byte(UpToDateMessage)) {
		println("[WARN] local Git was not up to date")
	}

	f, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(fmt.Errorf("[ERROR] couldn't open the data.json file: %w", err))
	}
	defer f.Close()

	for i := range cfg.Amount {
		if err := GitFlow(cfg, f, i); err != nil {
			panic(err)
		}
	}

	if cfg.isDebug() {
		fmt.Printf("[DEBUG] total amount: %d\n", cfg.Amount)
		fmt.Printf(
			"[DEBUG] flags - years: %d, months: %d, days: %d, max: %d, min: %d, amount: %d\n",
			cfg.Years,
			cfg.Months,
			cfg.Days,
			cfg.Max,
			cfg.Min,
			cfg.Amount,
		)
	}

	switch {
	case !cfg.Skip:
		fmt.Printf("Do you want to push all %d changes? (y/n) ", cfg.Amount)
		ans, err := scanLine(ctx)
		if err != nil {
			return
		}

		ans = strings.ToLower(strings.TrimSpace(ans))
		if ans == "n" {
			_, err = Reset(cfg)
			if err != nil {
				panic(fmt.Errorf("[ERROR] couldn't revert the changes: %w", err))
			}
			return
		}
		fallthrough
	default:
		_, err = Push(cfg)
		if err != nil {
			panic(fmt.Errorf("[ERROR] couldn't push to remote repository: %w", err))
		}
	}
}
