package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/philippgille/chromem-go"
	"github.com/siuyin/dflt"
)

type Entry struct {
	ID      string
	Title   string
	Content string
}

const embeddingModel = "nomic-embed-text"

var (
	db *chromem.DB
)


func main() {
	coll := loadOrCreateVecDB()
	for {
		q := getQuery()
		md := showMatchingDocs(coll, q)
		rag(md, q)
	}
}

func loadOrCreateVecDB() *chromem.Collection {
	var err error
	const compress = false
	db, err = chromem.NewPersistentDB(dflt.EnvString("VECTORDB", "/tmp/vecdb"), compress)
	if err != nil {
		log.Fatal(err)
	}

	coll, err := db.GetOrCreateCollection("knowledge", nil,
		chromem.NewEmbeddingFuncOllama(embeddingModel, ""))
	if err != nil {
		log.Fatal(err)
	}
	if coll.Count() > 0 {
		return coll
	}

	coll = loadCSVData(coll)
	log.Println("collection created")
	return coll
}

const searchQueryPrefix = "search_query: "

func getQuery() string {
	fmt.Println("\n\n\nEnter your question:")
	r := bufio.NewReader(os.Stdin)
	q, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return searchQueryPrefix + q
}

const searchDocPrefix = "search_document: "

func showMatchingDocs(coll *chromem.Collection, q string) []chromem.Result {
	const numRes = 3
	res, err := coll.Query(context.Background(), q, numRes, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	for i, r := range res {
		fmt.Printf("\nDocument %d (similarity: %.3f): %s\n",
			i+1, r.Similarity, strings.TrimPrefix(r.Content, searchDocPrefix))
	}

	return res
}

func loadCSVData(coll *chromem.Collection) *chromem.Collection {
	f, err := os.Open("./dat.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(f)
	recs, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(recs) < 2 {
		return coll
	}

	for i := 1; i < len(recs); i++ {
		e := Entry{ID: recs[i][0], Title: recs[i][1], Content: recs[i][2]}
		d := chromem.Document{
			ID:      e.ID,
			Content: searchDocPrefix + "\ndocument id: " + e.ID + "\n" + e.Title + "\n" + e.Content + "\n"}
		if err := coll.AddDocument(context.Background(), d); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", d.ID, d.Content)
	}
	return coll

}
