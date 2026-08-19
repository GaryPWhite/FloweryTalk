package main

import (
	"flag"
	"log"
	"os"

	. "github.com/garypwhite/flowerytalk/internal/g2p"
	. "github.com/garypwhite/flowerytalk/internal/synthesizer"
)

func main() {
	prompt := flag.String("prompt", "", "what you want flowery to say \n(hint: give as ('here I come') not (here I come))")
	output := flag.String("output", "~/flowery.wav", "path to save stitched voice clip")
	flag.Parse()
	if len(*prompt) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if len(*output) == 0 {
		*output = "~/flowery.wav"
	}
	phonemizer := NewPhonemizer()
	phonemes := phonemizer.Parse(*prompt)
	synthesizer := NewSynthesizer()
	files, err := synthesizer.ResolvePhonemeFiles(phonemes)
	if err != nil {
		log.Fatalf("could not resolve files, err %e", err)
	}
	err = SynthFiles(files, *output)
	if err != nil {
		log.Fatalf("failed to write synth %s ", err)
	}
}
