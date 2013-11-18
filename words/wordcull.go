package words

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

func checkWord(s string) bool {

	if len(s) < 4 {
		return false
	}

	if len(s) > 13 {
		return false
	}

	return true
}

func wordcull() {
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
