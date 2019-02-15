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

func GetArch() string {
	return os.Getenv("WIOARCH")
}

func GetMinWioVersion() string {
	return os.Getenv("CONFIG_MIN_WIO_VER")
}
