package g2p

import (
	"bufio"
	"embed"
	"log"
	"regexp"
	"strings"
)

type Phonemizer struct {
	phoneticWords map[string][]string
}

//go:embed cmudict-0.7b
var cmudict embed.FS

func NewPhonemizer() Phonemizer {
	cmu, err := cmudict.Open("cmudict-0.7b")
	if err != nil {
		log.Fatalf("failed to create Phoenemizer, %s", err)
	}
	defer cmu.Close()
	allWords := make(map[string][]string)
	b := bufio.NewScanner(cmu)
	for b.Scan() {
		line := b.Text()
		split := strings.Split(line, " ")
		if len(split) == 0 || strings.HasPrefix(line, ";;;") { // ignore empty or comment lines
			continue
		}
		// trim empty entries by starting at 2
		allWords[split[0]] = split[2:]
	}
	err = b.Err()
	if err != nil {
		log.Fatalf("failed to scan cmudict file %s", err)
	}
	return Phonemizer{phoneticWords: allWords}
}

// Parse will take a string of text, and give back a whole bunch of phonemes.
// the phonemes will be returned in the provided order, with punctuation preserved.
// Note that we do not parse apostrophes, etc -- punctuation is given back to allow
// for space between words (commas, periods, etc.)
func (p Phonemizer) Parse(sentence string) (phonemes []string) {
	re := regexp.MustCompile(`[^A-Za-z\s!?.,']`)
	punct := regexp.MustCompile("([!?.,])")
	sentence = re.ReplaceAllString(sentence, "")        // remove unneeded punct
	sentence = punct.ReplaceAllString(sentence, " $1 ") // add padding to parse punct
	for _, word := range strings.Split(sentence, " ") {
		if len(word) == 0 { // double spaces, etc. do not add anything
			continue
		}
		if punct.Match([]byte(word)) { // add punctuation by itself
			phonemes = append(phonemes, word)
			continue
		}
		word = strings.ToUpper(word)
		phones, ok := p.phoneticWords[word]
		if !ok {
			phones, ok = p.nearMatch(word)
			if !ok {
				phones = p.byLetter(word)
			}
		}
		if len(phones) == 0 {
			log.Printf("could not find phones for %s", word)
		}
		phonemes = append(phonemes, phones...)
	}
	return
}

// attempt to remove suffixes, if none exist, or no word exists without the suffix, false.
// if suffixes can be removed, this returns a []string of the word and suffix phonemes.
func (p Phonemizer) nearMatch(word string) ([]string, bool) {
	// remove common suffixes (s, -ing, -ed)
	for _, suffix := range []string{"s", "ing", "ed"} {
		wordWithout := strings.TrimSuffix(word, suffix)
		phonemes, ok := p.phoneticWords[wordWithout]
		if ok {
			suffixPhonemes, _ := p.phoneticWords[suffix]
			phonemes = append(phonemes, suffixPhonemes...)
			return phonemes, true
		}
	}
	return []string{}, false
}

// last resort, sound out the word out by letter.
func (p Phonemizer) byLetter(word string) []string {
	var acc []string
	for _, letter := range strings.Split(word, "") {
		phonemes, ok := p.phoneticWords[letter]
		if ok {
			acc = append(acc, phonemes...)
		}
	}
	return acc
}
