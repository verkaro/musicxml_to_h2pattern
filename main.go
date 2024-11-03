package main

import (
    "flag"
    "fmt"
    "log"
    "mxml2patt/config"
    "mxml2patt/parser"
    "mxml2patt/hydrogen"
)

func main() {
    // Define command-line flags
    inputFile := flag.String("input", "", "Path to the MusicXML input file (overrides config.yaml)")
    outputFile := flag.String("output", "", "Path for the output .h2pattern file (overrides config.yaml)")
    flag.Parse()

    // Load configuration
    cfg, err := config.LoadConfig("config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    fmt.Printf("Config loaded: %+v\n", cfg)

    // Override config paths with flags if provided
    if *inputFile != "" {
        cfg.InputFile = *inputFile
    }
    if *outputFile != "" {
        cfg.OutputFile = *outputFile
    }

    // Display help message if input file is not specified
    if cfg.InputFile == "" {
        fmt.Println("Usage: go run main.go [--input <input file>] [--output <output file>]")
        fmt.Println("Flags:")
        flag.PrintDefaults()
        return
    }

    // Parse MusicXML file
    notes, err := parser.ParseMusicXML(cfg.InputFile, cfg)
    if err != nil {
        log.Fatalf("Failed to parse MusicXML: %v", err)
    }
    fmt.Printf("Parsed notes: %+v\n", notes)

    // Calculate maximum ticks for scaling positions
    totalTicks := 0
    for _, note := range notes {
        if note.Position > totalTicks {
            totalTicks = note.Position
        }
    }

    // Define pattern length in ticks
    patternLength := 384

    // Create the Hydrogen pattern and distribute notes
    pattern := hydrogen.NewHydrogenPattern(cfg.Tempo, patternLength)
    for _, note := range notes {
        if totalTicks > 0 {
            scaledPosition := int(float64(note.Position) * float64(patternLength) / float64(totalTicks))
            pattern.AddNoteAtPosition(note.Instrument, note, scaledPosition)
        } else {
            pattern.AddNoteAtPosition(note.Instrument, note, 0)
        }
    }
    fmt.Printf("Generated Hydrogen pattern: %+v\n", pattern)

    // Export pattern to file
    err = pattern.ExportToFile(cfg.OutputFile)
    if err != nil {
        log.Fatalf("Failed to export pattern: %v", err)
    }
    fmt.Println("Pattern exported to", cfg.OutputFile)
}

