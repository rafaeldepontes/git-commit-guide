# Git-Commit-Guide

Git-Commit-Guide is a small Go tool that automates a Git workflow by updating a
file, creating commits with custom dates, and pushing the result to a
remote repository.

It is mainly useful for:

- generating artificial commit histories

---

## Features

- Pulls the latest changes before starting
- Writes a fresh timestamp to `data.json`
- Stages and commits changes automatically
- Supports custom commit dates
- Supports exact or random commit counts
- Includes a debug mode for inspecting runtime values
- Uses the Git CLI through Go's `os/exec`

---

## Requirements

- Go installed
- Git installed and available in your `PATH`
- A valid Git repository with a configured remote

---

## Installation

Clone the repository and build the binary:

```bash
git clone <your-repo>
cd git-commit-guide
go build
```

---

## Usage

Run the binary with optional flags:

```bash
./git-commit-guide [flags]
```

### Flags

Flag Description

---

`-years` Offset commit dates by years
`-months` Offset commit dates by months
`-days` Offset commit dates by days
`-date` Use an exact date in `YYYY-MM-DD` format
`-min` Minimum number of commits
`-max` Maximum number of commits
`-amount` Exact number of commits
`-debug` Enable debug output

---

## Date behavior

You can generate commits in two ways:

### Relative date

Use `-years`, `-months`, and/or `-days` to shift the commit date
relative to the current UTC date.

Example:

```bash
./git-commit-guide -days=-10 -amount=20
```

This creates 20 commits dated 10 days in the past.

### Exact date

Use `-date` to force a specific date.

Example:

```bash
./git-commit-guide -date=2024-01-15 -amount=10
```

When `-date` is provided, the program ignores the relative date flags and uses the supplied date instead.

The expected format is:

    YYYY-MM-DD

Invalid formats will cause the program to fail.

---

## Commit count behavior

- `-amount` creates exactly that many commits
- If `-amount` is not provided, the program picks a random value
  between `-min` and `-max`
- Default values:
  - `-min=8`
  - `-max=30`

Example:

```bash
./git-commit-guide -months=3 -min=5 -max=20
```

This generates a random number of commits between 5 and 20, with commit
dates shifted 3 months into the future.

---

## How it works

The tool follows this workflow:

1.  Pull the latest changes from the remote repository
2.  Update `data.json` with a new timestamp
3.  Stage the file
4.  Commit with a custom message and date
5.  Repeat until the configured amount is reached
6.  Push the commits to the remote repository

Internally, it uses Go's `os/exec` package to call Git commands
directly.

---

## Project structure

    .
    ├── main.go
    ├── config.go
    ├── utils.go
    ├── git_action.go
    └── data.json

### Files

- `main.go` --- application entry point and execution flow
- `config.go` --- CLI flag parsing and date handling
- `utils.go` --- validation helpers and usage strings
- `git_action.go` --- Git command wrappers and file writing logic
- `data.json` --- file updated before each commit

---

## Disclaimer

This project is intended for **artificial commit histories only**.

Use it carefully in real repositories, especially when rewriting commit
dates or generating large numbers of commits.

---

## Credits

Project built with [Gini](https://gini-webserver.up.railway.app/)
