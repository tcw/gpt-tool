package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/pkoukk/tiktoken-go-loader"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s [file]\n", os.Args[0])
		if err != nil {
			return
		}
		flag.PrintDefaults()
	}
	flag.Parse()

	argsWithoutProg := os.Args[1:]

	var text = ""

	if len(argsWithoutProg) == 0 {
		text = getTextFromStdin()
	} else if len(argsWithoutProg) == 1 {
		sourcePath, err := filepath.Abs(argsWithoutProg[0])
		check(err)
		bytes, err := os.ReadFile(sourcePath)
		check(err)
		text = string(bytes)
	}

	encoding := "cl100k_base"

	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return
	}

	// encode
	token := tke.Encode(text, nil, nil)

	//tokens
	//fmt.Println((token))
	// num_tokens
	fmt.Println(len(token))
}

func getTextFromStdin() string {
	var sb strings.Builder

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_, err := sb.WriteString(scanner.Text())
		check(err)
	}
	return sb.String()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
