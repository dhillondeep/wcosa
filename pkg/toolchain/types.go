package toolchain

type FrameworkConfig struct {
    Name          string   `json:"name"`
    Url           string   `json:"url"`
    Keywords      []string `json:"keywords"`
    ToolchainFile string   `json:"toolchain-file"`
    Platform      string   `json:"platform"`
    Framework     string   `json:"framework"`
    Requirements  []string `json:"requirements"`
}
