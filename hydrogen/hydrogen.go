package hydrogen

import (
    "encoding/xml"
    "fmt"
    "mxml2patt/parser"
    "os"
)

type HydrogenPattern struct {
    Tempo  int
    Tracks map[string][]HydrogenNote
    Length int
}

type HydrogenNote struct {
    Position int
    Volume   int
    Duration int
}

type drumkitPattern struct {
    XMLName           xml.Name          `xml:"drumkit_pattern"`
    Version           string            `xml:"version,attr"`
    PatternForDrumkit string            `xml:"pattern_for_drumkit"`
    Pattern           pattern           `xml:"pattern"`
}

type pattern struct {
    Name     string   `xml:"pattern_name"`
    Category string   `xml:"category"`
    Size     int      `xml:"size"`
    NoteList noteList `xml:"noteList"`
}

type noteList struct {
    Notes []instrumentNote `xml:"note"`
}

type instrumentNote struct {
    Position   int     `xml:"position"`
    LeadLag    int     `xml:"leadlag"`
    Velocity   float64 `xml:"velocity"`
    PanL       float64 `xml:"pan_L"`
    PanR       float64 `xml:"pan_R"`
    Pitch      int     `xml:"pitch"`
    Key        string  `xml:"key"`
    Length     int     `xml:"length"`
    Instrument int     `xml:"instrument"`
}

// NewHydrogenPattern creates a new pattern with the given tempo and length.
func NewHydrogenPattern(tempo, length int) *HydrogenPattern {
    return &HydrogenPattern{
        Tempo:  tempo,
        Tracks: make(map[string][]HydrogenNote),
        Length: length,
    }
}

// AddNoteAtPosition adds a note with a manually defined position to the specified instrument track.
func (hp *HydrogenPattern) AddNoteAtPosition(instrument string, note parser.ParsedNote, position int) {
    hydrogenNote := HydrogenNote{
        Position: position,
        Volume:   note.Volume,
        Duration: -1,
    }
    hp.Tracks[instrument] = append(hp.Tracks[instrument], hydrogenNote)
}

// ExportToFile exports the HydrogenPattern to a Hydrogen-compatible .h2pattern file.
func (hp *HydrogenPattern) ExportToFile(filePath string) error {
    pattern := drumkitPattern{
        Version:           "0.2",
        PatternForDrumkit: "GMkit",
        Pattern: pattern{
            Name:     "Generated Pattern",
            Category: "Uncategorized",
            Size:     384,
            NoteList: noteList{},
        },
    }

    // Map instrument names to unique numeric IDs
    instrumentIDMap := make(map[string]int)
    idCounter := 0
    for instrument := range hp.Tracks {
        instrumentIDMap[instrument] = idCounter
        idCounter++
    }

    // Convert tracks and notes with proper scaling and values
    for instrument, notes := range hp.Tracks {
        for _, note := range notes {
            noteEntry := instrumentNote{
                Position:   note.Position,
                LeadLag:    0,
                Velocity:   float64(note.Volume) / 127.0,
                PanL:       0.5,
                PanR:       0.5,
                Pitch:      0,
                Key:        "C0",
                Length:     -1,
                Instrument: instrumentIDMap[instrument],
            }
            pattern.Pattern.NoteList.Notes = append(pattern.Pattern.NoteList.Notes, noteEntry)
        }
    }

    // Open file for writing
    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer file.Close()

    // Write XML declaration manually
    _, err = file.WriteString("<?xml version=\"1.0\" ?>\n")
    if err != nil {
        return fmt.Errorf("failed to write XML declaration: %v", err)
    }

    // Serialize to XML with proper indentation
    encoder := xml.NewEncoder(file)
    encoder.Indent("", "  ")
    if err := encoder.Encode(pattern); err != nil {
        return fmt.Errorf("failed to encode XML: %v", err)
    }

    return nil
}

