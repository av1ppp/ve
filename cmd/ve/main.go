package main

import (
	"errors"
	"os"
)

const (
	appName = "ver"
	appDesc = "Versioning tool"

	actionInit = "init"
	actionIncr = "incr"
	actionGet  = "get"
	actionHelp = "help"

	versionFile = "VERSION"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		printHelp()
		return
	}

	if err := handleAction(args[0]); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func handleAction(action string) error {
	switch action {
	case actionInit:
		return handleActionInit()
	case actionIncr:
		return handleActionIncr()
	case actionGet:
		return handleActionGet()
	case actionHelp:
		printHelp()
		return nil
	}
	panic("Unknown action '" + action + "'")
}

func handleActionInit() error {
	_, err := os.Stat(versionFile)
	if err == nil {
		return errors.New("Version file '" + versionFile + "' already exists")
	}

	if !os.IsNotExist(err) {
		return err
	}

	v := NewVersion()
	return v.WriteFile(versionFile)
}

func handleActionIncr() error {
	oldv, err := ParseVersionFile(versionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Version file '" + versionFile + "' does not exist. Run '" + appName + " " + actionInit + "' first")
		}
		return err
	}

	newv := oldv.Increment()
	return newv.WriteFile(versionFile)
}

func handleActionGet() error {
	version, err := ParseVersionFile(versionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Version file '" + versionFile + "' does not exist. Run '" + appName + " " + actionInit + "' first")
		}
		return err
	}
	os.Stdout.WriteString(version.String())
	return nil
}

func printHelp() {
	s := appName + " - " + appDesc + "\n\n"
	s += "Usage:\n"
	s += "  " + appName + " <action>\n\n"
	s += "Actions:\n"
	s += "  init - Initialize ver in current directory\n"
	s += "  incr - Increment version\n"
	s += "  get  - Get current version\n"
	s += "  help - Print this message\n"
	os.Stdout.WriteString(s + "\n")
}
