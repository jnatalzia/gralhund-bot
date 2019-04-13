package commands

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type GameDoc struct {
	name        string
	alignment   string
	domains     []string
	description string
}

var gameDocs []GameDoc
var mappedDocs = make(map[string]*GameDoc)
var projectRoot = os.Getenv("PROJECTROOT")

func RetrieveDocs() {
	if projectRoot == "" {
		fmt.Println("You must specify a $PROJECTROOT")
		os.Exit(1)
	}

	openFile, err := os.Open(path.Join(projectRoot, "docs/gods.csv"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	r := csv.NewReader(openFile)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	gameDocs = []GameDoc{}
	relevantRecords := records[1:]
	for _, r := range relevantRecords {
		godFirstName := strings.Split(r[0], " ")[0]
		godFirstName = strings.ToLower(strings.ReplaceAll(godFirstName, ",", ""))

		result := GameDoc{
			name:        r[0],
			alignment:   r[1],
			domains:     strings.Split(r[2], ","),
			description: r[3],
		}
		gameDocs = append(gameDocs, result)
		mappedDocs[godFirstName] = &result
	}
}

func writeDocForGod(doc GameDoc, extended bool) string {
	result := ""
	result = result + "**" + doc.name + "**" + " " + "_" + doc.alignment + "_\n"
	result = result + "Domains: " + strings.Join(doc.domains, ", ") + "\n"
	if extended {
		result = result + "Description: " + doc.description + "\n"
	}
	result = result + "\n"
	return result
}

func GetDocs() string {
	result := "**The Gods of Taldorei**\n--------------------------\n"

	for _, doc := range gameDocs {
		result = result + writeDocForGod(doc, false)
	}

	return result
}

func GetDoc(godName string) (string, error) {
	god := mappedDocs[godName]
	if god == nil {
		return "", errors.New("God not found")
	}
	return writeDocForGod(*god, true), nil
}
