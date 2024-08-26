package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func twoDigitString(n int) string {
	s := strconv.Itoa(n)
	if len(s) < 2 {
		s = fmt.Sprintf("0%s", s)
	}
	return s
}

func newEntry(path string) error {
	now := time.Now()
	ep := fmt.Sprintf(
		"%s/%s/%s/%s%s",
		path,
		twoDigitString(now.Year()),
		twoDigitString(int(now.Month())),
		twoDigitString(now.Day()),
		".md",
	)
	ts := fmt.Sprintf(
		"# %s:%s\n\n\n",
		twoDigitString(now.Hour()),
		twoDigitString(now.Minute()),
	)
	// This overwrites the file, and that's bad...
	os.WriteFile(ep, []byte(ts), os.ModeAppend)
	editor, ok := os.LookupEnv("EDITOR")
	if !ok {
		fmt.Println("$EDITOR not set. Using vim.")
		editor = "vim"
	}
	bin, err := exec.LookPath(editor)
	if err != nil {
		panic(err)
	}
	args := []string{editor, "+", ep}
	err = syscall.Exec(bin, args, os.Environ())
	if err != nil {
		panic(err)
	}
	return nil
}

func handleTilde(jp string) string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(
		h,
		filepath.Base(jp),
	)
}

func main() {
	app := viper.New()
	app.SetDefault("journal_path", "~/journal/")
	app.SetConfigType("json")
	app.SetConfigName("settings.json")
	app.AddConfigPath(".")
	err := app.ReadInConfig()
	if err != nil {
		panic(err)
	}
	jpath := app.GetString("journal_path")
	if strings.HasPrefix(jpath, "~") {
		jpath = handleTilde(jpath)
	}
	newEntry(jpath)
}
