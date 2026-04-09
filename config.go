package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

type Config struct {
	Years       int
	Months      int
	Days        int
	NanoSeconds int
	Max         int
	Min         int
	Amount      int
	Debug       bool
	Skip        bool
	RawDate     []byte
}

func parseFlags() Config {
	var cfg Config

	flag.IntVar(&cfg.Years, "years", 0, YearsExplanation)
	flag.IntVar(&cfg.Months, "months", 0, MonthsExplanation)
	flag.IntVar(&cfg.Days, "days", 0, DaysExplanation)
	flag.IntVar(&cfg.Max, "max", 30, MaxExplanation)
	flag.IntVar(&cfg.Min, "min", 8, MinExplanation)
	flag.IntVar(&cfg.Amount, "amount", 0, CommitAmountExplanation)
	flag.BoolVar(&cfg.Debug, "debug", false, DebugExplanation)
	flag.BoolVar(&cfg.Skip, "y", false, SkipExplanation)
	date := flag.String("date", "", DateExplanation)

	flag.Parse()

	cfg.RawDate = []byte(*date)

	return cfg
}

func (c Config) isDebug() bool {
	return c.Debug
}

func (c *Config) ParseDate() {
	date, err := time.Parse(DateLayout, string(c.RawDate))
	if err != nil {
		panic(fmt.Errorf("[ERROR] couldn't parse date from raw date: %s, %w", c.RawDate, err))
	}

	c.Years = date.Year()
	c.Months = int(date.Month())
	c.Days = date.Day()
	c.NanoSeconds = rand.IntN(100)
}

func (c Config) GetDate(ith int) string {
	if len(c.RawDate) > 0 {
		var sb strings.Builder
		sb.WriteString(
			time.Date(
				c.Years,
				time.Month(c.Months),
				c.Days,
				23,
				0,
				0,
				c.NanoSeconds+ith,
				time.UTC,
			).Format(
				time.RFC3339Nano,
			),
		)
		return sb.String()
	}

	return time.Now().UTC().AddDate(
		c.Years,
		c.Months,
		c.Days,
	).Format(
		time.RFC3339Nano,
	)
}
