package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postres"
	dbname   = "postgres"
)

func main() {

	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s [text]\n", os.Args[0])
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

	embedding := fetchEmbedding(text)

	insertEmbedding(text, embedding)
}

func insertEmbedding(content string, embedding []float64) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("vector-db-cmd", psqlconn)
	check(err)

	defer db.Close()

	insertStmt := fmt.Sprintf(`insert into documents(content, embedding) values('%s', '%v')`,
		content, embedding)
	_, e := db.Exec(insertStmt)
	check(e)
}

func fetchEmbedding(text string) []float64 {
	url := "https://api.openai.com/v1/embeddings"

	var jsonStr = []byte(fmt.Sprintf(`'{
    "input": "%s",
    "model": "text-embedding-ada-002"
  }'`, text))

	openAiKey := os.Getenv("OPENAI_API_KEY")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)

	return toEmbedding(string(body))
}

func findClosestNeighbors(search string) []int64 {
	return nil //todo
}

func toEmbedding(emb string) []float64 {
	in := []byte(emb)

	var ebr EmbeddingResponse
	err := json.Unmarshal(in, &ebr)
	if err != nil {
		panic(err)
	}
	embedding := ebr.Data[0].Embedding

	return embedding
}

type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
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
