package synthesizer

import (
	"os"
	"regexp"
	"strings"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/generators"
	"github.com/gopxl/beep/wav"
)

const (
	samplesPerSecond = 16000
	punct_regexp     = "^([!?.,]+)$" // detect punctuation to add silence/pauses
)

// accepts a list of phonemes, and generates text-to-speech for those phonemes.
// will generate a `.wav` file with appropriate timestamp / name provided
func SynthFiles(files [][]string, outputPath string) error {
	synthb := beep.NewBuffer(beep.Format{SampleRate: beep.SampleRate(samplesPerSecond), NumChannels: 2, Precision: 2})
	punct := regexp.MustCompile(punct_regexp)
	// fill buffer with loaded files
	for _, word := range files {
		for _, file := range word {
			if punct.Match([]byte(file)) {
				synthb.Append(punctSilence(file))
				continue
			}
			// select a random file from the list and add it to the buffer
			streamer, err := LoadAudioClip(file)
			if err != nil {
				return err
			}
			synthb.Append(streamer)
		}
		synthb.Append(generators.Silence(samplesPerSecond / 25)) // word boundaries
	}
	// save the buffer to a file
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()
	err = wav.Encode(out, synthb.Streamer(0, synthb.Len()), synthb.Format())
	if err != nil {
		return err
	}
	return nil
}

// TODO: add punctuation silence based on length/type.
func punctSilence(punct string) beep.Streamer {
	fullSilence := beep.NewBuffer(beep.Format{SampleRate: beep.SampleRate(samplesPerSecond), NumChannels: 2, Precision: 2})
	for _, rn := range strings.Split(punct, "") {
		divider := 4
		if rn == "," {
			divider = 10 // smaller stops for commas
		}
		fullSilence.Append(generators.Silence((samplesPerSecond / divider)))
	}
	return fullSilence.Streamer(0, fullSilence.Len())
}
