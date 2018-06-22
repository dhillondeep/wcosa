package config

type avrDefaults struct {
    Ide           string
    Framework     string
    Port          string
    AVRBoard      string
    Baud          int
    DefaultTarget string
    AppTargetName string
    PkgTargetName string
}

var ProjectDefaults = avrDefaults{
    Ide:           "none",
    Framework:     "cosa",
    Port:          "none",
    AVRBoard:      "uno",
    Baud:          9600,
    DefaultTarget: "default",
    AppTargetName: "main",
    PkgTargetName: "test",
}
