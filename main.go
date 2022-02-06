package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	target := ""
	flag.StringVar(&target, "target", target, "tinygo target")
	flag.Parse()

	if target == "" {
		fmt.Fprintf(os.Stderr, "usage: %s --target TARGET\n", os.Args[0])
		os.Exit(1)
	}

	err := run(target)
	if err != nil {
		log.Fatal(err)
	}
}

func run(target string) error {
	tags, err := getTags(target)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "list", fmt.Sprintf("-tags=%s", strings.Join(tags, ",")), "-f", `{{join .GoFiles "\n"}}`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getTags(target string) ([]string, error) {
	buf := bytes.Buffer{}
	cmd := exec.Command("tinygo", "info", "-json", target)
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	x := struct {
		GOROOT     string   `json:"goroot"`
		GOOS       string   `json:"goos"`
		GOARCH     string   `json:"goarch"`
		GOARM      string   `json:"goarm"`
		BuildTags  []string `json:"build_tags"`
		GC         string   `json:"garbage_collector"`
		Scheduler  string   `json:"scheduler"`
		LLVMTriple string   `json:"llvm_triple"`
	}{}

	err = json.Unmarshal(buf.Bytes(), &x)
	if err != nil {
		return nil, err
	}

	return x.BuildTags, nil
}
