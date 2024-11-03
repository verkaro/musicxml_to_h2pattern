package parser

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "log"
    "mxml2patt/config"
)

type ParsedNote struct {
    Instrument string
    Position   int
    Volume     int
    Duration   int
}

type ScorePartwise struct {
    Parts []Part `xml:"part"`
}

type Part struct {
    Measures []Measure `xml:"measure"`
}

type Measure struct {
    Notes []Note `xml:"note"`
}

type Note struct {
    Pitch    Pitch   `xml:"pitch"`
    Duration int     `xml:"duration"`
    Rest     *Rest   `xml:"rest"`
}

type Pitch struct {
    Step   string `xml:"step"`
    Octave int    `xml:"octave"`
}

type Rest struct{}

func ParseMusicXML(filePath string, cfg *config.Config) ([]ParsedNote, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var score ScorePartwise
    err = xml.Unmarshal(data, &score)
    if err != nil {
        return nil, fmt.Errorf("failed to parse XML: %v", err)
    }

    var notes []ParsedNote
    for _, part := range score.Parts {
        for _, measure := range part.Measures {
            for _, note := range measure.Notes {
                if note.Rest != nil {
                    continue
                }

                // Combine step and octave to form the full pitch (e.g., C4)
                pitch := fmt.Sprintf("%s%d", note.Pitch.Step, note.Pitch.Octave)
                instrument, exists := cfg.DrumMappings[pitch]
                if !exists {
                    log.Printf("No mapping for pitch: %s", pitch)
                    continue
                }

                parsedNote := ParsedNote{
                    Instrument: instrument,
                    Position:   0,
                    Volume:     cfg.VolumeLevels[instrument],
                    Duration:   note.Duration,
                }
                notes = append(notes, parsedNote)
            }
        }
    }

    return notes, nil
}

