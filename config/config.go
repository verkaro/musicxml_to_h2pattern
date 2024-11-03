package config

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Tempo          int               `yaml:"tempo"`
    DrumMappings   map[string]string `yaml:"drum_mappings"`
    PatternLength  int               `yaml:"pattern_length"`
    Division       int               `yaml:"division"`
    Swing          float64           `yaml:"swing"`
    VolumeLevels   map[string]int    `yaml:"volume_levels"`
    InputFile      string            `yaml:"input_file"`
    OutputFile     string            `yaml:"output_file"`
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(filePath string) (*Config, error) {
    var config Config
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }
    return &config, nil
}

