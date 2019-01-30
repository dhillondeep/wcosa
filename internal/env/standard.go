package env

import "os"

func GetWioPath() string {
    return os.Getenv("WIOPATH")
}

func GetWioRoot() string {
    return os.Getenv("WIOROOT")
}

func GetOS() string {
    return os.Getenv("WIOOS")
}
