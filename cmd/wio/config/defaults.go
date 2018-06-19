package config

type defaults struct {
    Ide           string
    Framework     string
    Port          string
    AVRBoard      string
    Baud          int
    DefaultTarget string
    AppTargetName string
    PkgTargetName string
}

var ProjectDefaults = defaults{
    Ide:           "none",
    Framework:     "cosa",
    Port:          "none",
    AVRBoard:      "uno",
    Baud:          9600,
    DefaultTarget: "default",
    AppTargetName: "main",
    PkgTargetName: "test",
}
