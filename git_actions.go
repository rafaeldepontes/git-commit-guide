package main

import (
	"fmt"
	"os"
	"os/exec"
)

// GitFlow goes through all the git conventions up until committing.
//
// This function has no power over push or pull, it can only modify, stage and commit changes.
func GitFlow(cfg Config, f *os.File, ith int) error {
	fomartedDate := cfg.GetDate(ith)

	if err := WriteDate(cfg, f, fomartedDate); err != nil {
		return fmt.Errorf("[ERROR] couldn't modify the 'data.json' file: %w", err)
	}

	if err := Commit(cfg, FilePath, fomartedDate); err != nil {
		return fmt.Errorf("[ERROR] couldn't stage changes: %w", err)
	}
	return nil
}

// Pull pulls.
func Pull(cfg Config) ([]byte, error) {
	if cfg.isDebug() {
		fmt.Println("[DEBUG] pulling from git repository")
	}
	cmd := exec.Command("git", "pull")
	return cmd.CombinedOutput()
}

// Add stage the changes.
func Add(cfg Config, fPath string) error {
	if cfg.isDebug() {
		fmt.Println("[DEBUG] staging changes")
	}
	cmd := exec.Command("git", "add", fPath)
	_, err := cmd.CombinedOutput()
	return err
}

// Commit commit the staged changes.
func Commit(cfg Config, fPath string, date string) error {
	if cfg.isDebug() {
		fmt.Printf("[DEBUG] committing change: %s\n", date)
	}
	cmd := exec.Command(
		"git",
		"commit",
		fPath,
		"-m",
		fmt.Sprintf("\"feat: %s\"", date),
		fmt.Sprintf("--date=%s", date),
	)
	return cmd.Run()
}

// Push push the changes.
func Push(cfg Config) ([]byte, error) {
	if cfg.isDebug() {
		fmt.Println("[DEBUG] pushing to remote repository")
	}
	cmd := exec.Command("git", "push")
	return cmd.CombinedOutput()
}

// WriteDate only has access to 'data.json' and write access only.
func WriteDate(cfg Config, f *os.File, date string) error {
	if cfg.isDebug() {
		fmt.Println("[DEBUG] writing new date to 'data.json'")
	}
	_, err := f.WriteAt(fmt.Appendf([]byte(""), "{\n  \"date\": \"%s\"\n}", date), 0)
	return err
}
