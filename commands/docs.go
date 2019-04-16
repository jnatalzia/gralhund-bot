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

type GodDoc struct {
	name        string
	alignment   string
	domains     []string
	description string
}

var godDocs []GodDoc
var mappedGodDocs = make(map[string]*GodDoc)

type CharDoc struct {
	name        string
	isNPC       bool
	race        string
	class       string
	description string
}

var charDocs []CharDoc
var mappedCharDocs = make(map[string]*CharDoc)

var projectRoot = os.Getenv("PROJECTROOT")

func retrieveDocsFromFile(docsfile string) [][]string {
	if projectRoot == "" {
		fmt.Println("You must specify a $PROJECTROOT")
		os.Exit(1)
	}

	openFile, err := os.Open(path.Join(projectRoot, docsfile))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	r := csv.NewReader(openFile)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

func determineFirstNameFromString(name string) string {
	return strings.ToLower(strings.ReplaceAll(strings.Split(name, " ")[0], ",", ""))
}

func retrieveGodDocs() {
	records := retrieveDocsFromFile("docs/gods.csv")
	godDocs = []GodDoc{}
	relevantRecords := records[1:]
	for _, r := range relevantRecords {
		godFirstName := determineFirstNameFromString(r[0])

		result := GodDoc{
			name:        r[0],
			alignment:   r[1],
			domains:     strings.Split(r[2], ","),
			description: r[3],
		}
		godDocs = append(godDocs, result)
		mappedGodDocs[godFirstName] = &result
	}
}

func retrieveCharDocs() {
	records := retrieveDocsFromFile("docs/characters.csv")
	charDocs = []CharDoc{}
	relevantRecords := records[1:]
	for _, r := range relevantRecords {
		charFirstName := determineFirstNameFromString(r[0])
		result := CharDoc{
			name:        r[0],
			isNPC:       r[1] == "yes",
			race:        r[2],
			class:       r[3],
			description: r[4],
		}
		charDocs = append(charDocs, result)
		mappedCharDocs[charFirstName] = &result
	}
}

func RetrieveDocs() {
	retrieveGodDocs()
	retrieveCharDocs()
}

func writeDocForGod(doc GodDoc, extended bool) string {
	result := ""
	result = result + "**" + doc.name + "**" + " " + "_" + doc.alignment + "_\n"
	result = result + "Domains: " + strings.Join(doc.domains, ", ") + "\n"
	if extended {
		result = result + "Description: " + doc.description + "\n"
	}
	result = result + "\n"
	return result
}

func GetGodDocs() string {
	result := "**The Gods of Taldorei**\n--------------------------\n"

	for _, doc := range godDocs {
		result = result + writeDocForGod(doc, false)
	}

	return result
}

func GetGodDoc(godName string) (string, error) {
	god := mappedGodDocs[godName]
	if god == nil {
		return "", errors.New("God not found")
	}
	return writeDocForGod(*god, true), nil
}

func writeDocForChar(doc CharDoc, extended bool) string {
	result := ""
	result = result + "**" + doc.name + "**\n" + "  _The " + doc.race + " " + doc.class + "_\n"
	if extended {
		result = result + "\n" + doc.description + "\n"
	}
	result = result + "\n"
	return result
}

func GetCharDocs() string {
	result := "**The Characters of our Tale**\n--------------------------\n"

	for _, doc := range charDocs {
		result = result + writeDocForChar(doc, false)
	}

	return result
}

func GetCharDoc(charName string) (string, error) {
	char := mappedCharDocs[charName]
	if char == nil {
		return "", errors.New("Character not found")
	}
	return writeDocForChar(*char, true), nil
}
