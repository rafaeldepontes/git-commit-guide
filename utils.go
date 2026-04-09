package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
)

const (
	FilePath        = "./data.json"
	UpToDateMessage = "Already up to date"

	// Constraints
	DateLayout = "2006-01-02"
	Separator  = "-"

	// Usages
	CommitAmountExplanation = `
	If the amount is specified it will ignore completly the 'min' and 'max' limitations, committing the exact amount.
	`

	YearsExplanation = `
	If positive the commit will happen in the future. The date can be calculated by the following formula:
	"today's date + <years_flag>"
	(e.g "binary -years=10" -> 02/01/2010 + 10 = 02/01/2020)

	If negative the commit will happen in the past. The date can be calculated by the following formula:
	"today's date + <years_flag>"
	(e.g "binary -years=-10" -> 02/01/2010 + -10 = 02/01/2000)
	`

	MonthsExplanation = `
	If positive the commit will happen in the future date. The date can be calculated by the following formula:
	"today's date + <months_flag>"
	(e.g "binary -months=10" -> 02/01/2010 + 10 = 02/11/2010)

	If negative the commit will happen in the past. The date can be calculated by the following formula:
	"today's date + <months_flag>"
	(e.g "binary -months=-10" -> 02/01/2010 + -10 = 02/03/2009)
	`

	DaysExplanation = `
	If positive the commit will happen in the future date. The date can be calculated by the following formula:
	"today's date + <days_flag>"
	(e.g "binary -days=10" -> 12/01/2010 + 10 = 22/01/2010)

	If negative the commit will happen in the past. The date can be calculated by the following formula:
	"today's date + <days_flag>"
	(e.g "binary -days=-10" -> 12/01/2010 + -10 = 02/01/2010)
	`

	MaxExplanation = `
	If the maximum is specified it CAN BE committed up to the maximum amount specified.

	If the value is lower than the minimum and the minimum isn't specified it will keep the maximum value and the minimum will be adjusted
	(e.g "binary -max=-2 ["min not specified, then min equals to 8 (default value)"] = max becomes its absolute value + min).

	If the value is lower than the minimum and the minimum is specified it will adjust the maximum relative to the minimum to keep a valid range
	(e.g "binary -min=20 -max=5 ["max is specified"] = min keeps its value (20) and max becomes min+max or 20+5...... 25).
	`

	MinExplanation = `
	If the minimum is specified it CAN BE committed the minimum amount or higher.

	If the value is higher than the default maximum and the maximum isn't specified it will add the minimum value to the max value keeping some range between them
	(e.g "binary -min=35 ["max not specified, then max equals to 30 (default value)"] = min keeps its value (35) and max becomes min+30 or 35+30...... 65).

	If the value is higher than the default maximum and the maximum is specified it will add the minimum value to the max value keeping some range between them
	(e.g "binary -min=35 -max=21["max is specified"] = min keeps its value (35) and max becomes min+max or 35+21...... 56).
	`
	DateExplanation = `
	When the -date flag is provided, the program ignores today's date and uses the date supplied by the user instead.

    The expected format is "YYYY-MM-DD". If any other format is provided, the program will fail.
	`

	DebugExplanation = `
	Activates debug logs and metrics.
	`

	SkipExplanation = `
	When -y is provided, the program will simply push all the changes to the remote repository, else it will ask if you really want to push those changes.
	`
)

// validateFlags holds business logic and some validations to asure behaviour.
func validateFlags(cfg *Config) {
	if cfg.Max <= 0 {
		cfg.Max = abs(cfg.Max)
	}

	if cfg.Min <= 0 {
		cfg.Min = abs(cfg.Min)
	}

	if cfg.Min >= cfg.Max {
		cfg.Max = cfg.Min + cfg.Max
	}

	if cfg.Amount <= 0 {
		cfg.Amount = rand.IntN(cfg.Max-cfg.Min+1) + cfg.Min
	}

	if validateDate(cfg.RawDate) {
		cfg.ParseDate()
	}
}

func abs(n int) int {
	if n > 0 {
		return n
	}
	return n * -1
}

func validateDate(date []byte) bool {
	if len(date) == 0 {
		return false
	}

	if len(date) != len(DateLayout) {
		panic(
			fmt.Sprintf(
				"[ERROR] invalid date lenght, your lenght was %d should be %d", len(date), len(DateLayout),
			),
		)
	}

	if !isNumber(date[:4]) ||
		string(date[4:5]) != Separator ||
		!isNumber(date[5:7]) ||
		string(date[7:8]) != Separator ||
		!isNumber(date[8:]) {
		panic(
			fmt.Sprintf(
				"[ERROR] invalid date format, your date was %s expected YYYY-MM-DD", string(date),
			),
		)
	}
	return true
}

func isNumber(src []byte) bool {
	_, err := strconv.Atoi(string(src))
	return err == nil
}

// scanLine uses channels and a context to verify if the user chose to
// cancel the build, it also uses a Scanner from bufio to work.
func scanLine(ctx context.Context) (string, error) {
	ch := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			ch <- scanner.Text()
			return
		}
		if err := scanner.Err(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return "", errors.New("Reverting changes...")
	case err := <-errCh:
		return "", err
	case line := <-ch:
		return line, nil
	}
}
