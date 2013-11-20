// Package wordcull includes the logic to include and exclude
// sets of words for creation of pass phrases.
package wordcull

import (
	"bufio"
	"fmt"
	"log"
	//"sort"
	"strings"
	//"io"
	"os"
	//"strings"
)

// portmanteau determines if a word is made up of precisely two
// other words in a given list.
func portmanteau(s string, list map[string]bool) bool {
	pres, sufs := make([]string, 0), make([]string, 0)
	for test := range list {
		if strings.HasPrefix(s, test) {
			pres = append(pres, test)
		}
		if strings.HasSuffix(s, test) {
			sufs = append(sufs, test)
		}
	}

	for _, pre := range pres {
		for _, suf := range sufs {
			if pre+suf == s {
				//fmt.Println(pre, suf, "=", s)
				return true
			}
		}
	}

	return false

	// do this
}

// checkWord determines whether a given word is within two bounds.
//
// TODO(Luke): Change this from magic numbers to user-specified
// bounds.
func checkWord(s string) bool {

	if len(s) < 4 {
		return false
	}

	if len(s) > 13 {
		return false
	}

	return true
}

// Wordcull provides the main logic for the package.
//
// TODO(Luke): Move the list of excluded words to a text file.
//
// TODO(Luke): Allow user-specified culling.
func Wordcull() {
	f, _ := os.Open("google-10000-english.txt")
	defer f.Close()
	fo, _ := os.Create("newwords.txt")
	defer fo.Close()

	words := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := strings.Trim(scanner.Text(), "\t\n\f\r ")
		words[word] = checkWord(word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Remove words solely on arbitrage
	words["ment"] = false
	words["ping"] = false
	words["trans"] = false
	words["tion"] = false
	words["inter"] = false
	words["tions"] = false
	words["hart"] = false
	words["comm"] = false
	words["para"] = false
	words["deutsch"] = false
	words["deutschland"] = false
	words["milf"] = false
	words["milfhunter"] = false
	words["mary"] = false
	words["maryland"] = false
	words["conf"] = false
	words["anti"] = false
	words["cock"] = false
	words["cincinnati"] = false
	words["sexcam"] = false
	words["xnxx"] = false
	words["fuck"] = false
	words["hentai"] = false

	//sortList(words)
	for word := range words {
		if portmanteau(word, words) {
			words[word] = false
		}
	}

	printer := make([]string, 0)

	for word := range words {
		if words[word] {
			printer = append(printer, word)
		}
	}

	sortList(printer)

	for _, word := range printer {
		fmt.Fprintln(fo, word)
	}
}
