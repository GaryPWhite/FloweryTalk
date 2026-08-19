package synthesizer

import (
	"fmt"
	"regexp"
)

// Returns a list of files (embedded) for a given list of phonemes.
// These phoneme files can be given to LoadAudioClips to get sound file contents.
// Gives an error when no phoneme can be found, even in fallbacks.
// When a phoneme is not found, a list of fallbacks will be given. Fallback list is
// about 3x larger, and is sorted by the provided manifest.json (best-match-first)
func (s Synthesizer) ResolvePhonemeFiles(phonemes []string) ([][]string, error) {
	allFiles := make([][]string, len(phonemes))
	punct := regexp.MustCompile(punct_regexp) // detect punctuation to add silence/pauses
	for i, p := range phonemes {
		if punct.Match([]byte(p)) {
			allFiles[i] = []string{p}
			continue
		} else if p == "\\w" {
			allFiles[i] = []string{"\\w"}
			continue
		}

		files, ok := s.clips[p]
		if ok {
			allFiles[i] = files
			continue
		}
		// could not find any clips for phoneme, try fallbacks
		fbFiles := []string{}
		fallbacks, ok := s.fallbacks[p]
		for _, fb := range fallbacks {
			files, ok := s.clips[fb]
			if ok { // TODO: could add some randomness here, shuffle fallbacks before `range`ing it
				fbFiles = append(fbFiles, files...)
				break // let's use the first found suggestion...
			}
		}
		if !ok || len(fbFiles) == 0 { // check both failure conditions, no fallback and/or no fallback files
			return allFiles, fmt.Errorf("Could not produce phoneme for <%s>, no phoneme or fallback found.", p)
		}
		allFiles[i] = fbFiles
	}
	return allFiles, nil
}
