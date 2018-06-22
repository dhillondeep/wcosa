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

var AvrProjectDefaults = avrDefaults{
    Ide:           "none",
    Framework:     "cosa",
    Port:          "none",
    AVRBoard:      "uno",
    Baud:          9600,
    DefaultTarget: "default",
    AppTargetName: "main",
    PkgTargetName: "test",
}

type nativeDefaults struct {
    Ide           string
    Language      string
    Port          string
    DefaultTarget string
    AppTargetName string
    PkgTargetName string
}

var NativeProjectDefaults = nativeDefaults{
    Ide:           "none",
    Language:      "none",
    Port:          "none",
    DefaultTarget: "default",
    AppTargetName: "main",
    PkgTargetName: "test",
}
