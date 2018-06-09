package config

type defaults struct {
    Ide           string
    Framework     string
    Platform      string
    Port          string
    Board         string
    DefaultTarget string
}

var ProjectDefaults = defaults{
    Ide:           "none",
    Framework:     "cosa",
    Platform:      "avr",
    Port:          "none",
    Board:         "uno",
    DefaultTarget: "default",
}
