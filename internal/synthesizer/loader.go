package synthesizer

import (
	"embed"
	"encoding/json"
	"log"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/wav"
)

//go:embed */*.wav
var soundFiles embed.FS

//go:embed phonemes/manifest.json
var manifestBytes []byte

//go:embed phonemes/fallbacks.json
var fallbacksBytes []byte

//go:embed words/manifest.json
var wordsManifestBytes []byte

type Synthesizer struct {
	clips     map[string][]string
	fallbacks map[string][]string
	words     map[string][]string
}

// Provides a
func NewSynthesizer() Synthesizer {
	// read clips / fallback manifest and return as a synth
	var synth Synthesizer
	err := json.Unmarshal(manifestBytes, &synth.clips)
	if err != nil {
		log.Fatalf("Failed to load manifest %s", err)
	}
	err = json.Unmarshal(fallbacksBytes, &synth.fallbacks)
	if err != nil {
		log.Fatalf("Failed to load manifest %s", err)
	}
	err = json.Unmarshal(wordsManifestBytes, &synth.words)
	if err != nil {
		log.Fatalf("Failed to load words manifest %s", err)
	}
	return synth
}

// Fetches clips from embedded phoneme files.
// Returns file contents in sequential order, the same order as files given.
// Produces an err if file is not found / unreadable.
func LoadAudioClip(file string) (beep.Streamer, error) {
	raw, err := soundFiles.Open(file)
	if err != nil {
		log.Printf("failed to read phoneme file given %s", file)
		return nil, err
	}

	streamer, _, err := wav.Decode(raw)
	if err != nil {
		log.Printf("failed to decode file %s", file)
		return nil, err
	}
	return streamer, nil
}
