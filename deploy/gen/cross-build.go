package main

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "github.com/mholt/archiver"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

type Config struct {
    ProjectName string   `yaml:"project_name"`
    SrcFolder   string   `yaml:"src_folder"`
    DestFolder  string   `yaml:"dest_folder"`
    Targets     []string `yaml:"targets"`
    OtherFiles  []string `yaml:"other_files"`
}

func Execute(dir string, name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func createString(appName, platform, architecture, extension string) string {
    return appName + "_" + platform + "_" + architecture + extension
}

func buildBinaries(targets, projectName, srcFolder, destFolder string) error {
    var destFolderAbs string
    var err error

    if destFolderAbs, err = filepath.Abs(destFolder); err != nil {
        return err
    }

    if _, err := os.Stat(destFolderAbs); !os.IsNotExist(err) {
        if err :=  os.RemoveAll(destFolderAbs); err != nil {
            return err
        }
    }

    if err := os.Mkdir(destFolder, os.ModePerm); err != nil {
        return err
    }

    dir, err := os.Getwd()
    if err != nil {
        return err
    }

    err = Execute(dir, "xgo", "--targets="+targets, "-out", destFolder+"/"+projectName, srcFolder)
    if err != nil {
        return err
    }

    if srcFolder, err = filepath.Abs(srcFolder); err != nil {
        return err
    }

    fmt.Println("Renaming the packages...")

    files, err := ioutil.ReadDir(destFolderAbs)
    if err != nil {
        return err
    }

    archMapping := map[string]string{
        "386":   "i386",
        "amd64": "x86_64",
        "arm-5": "arm5",
        "arm-6": "arm6",
        "arm-7": "arm7",
    }

    osMapping := map[string]string{
        "darwin":  "",
        "linux":   "",
        "windows": ".exe",
    }

    for _, f := range files {
        if strings.Contains(f.Name(), projectName) {
            for osName, osExt := range osMapping {
                if strings.Contains(f.Name(), osName) {
                    for archName, archVal := range archMapping {
                        if strings.Contains(f.Name(), archName) {
                            newName := createString(projectName, osName, archVal, osExt)
                            os.Rename(destFolderAbs+"/"+f.Name(), destFolderAbs+"/"+newName)
                        }
                    }
                }
            }
        }
    }

    return nil
}

func packAndCompress(destFolder string, otherFiles []string) error {
    fmt.Println("Compressing all the files needed...")

    var otherFilesAbs []string

    // convert other file paths to abs paths
    for _, file := range otherFiles {
        newPath, err := filepath.Abs(file)
        if err != nil {
            return err
        }

        otherFilesAbs = append(otherFilesAbs, newPath)
    }


    osCompressMapping := map[string]string{
        "darwin":  "tar.gz",
        "linux":   "tar.gz",
        "windows": "zip",
    }

    files, err := ioutil.ReadDir(destFolder)
    if err != nil {
        return err
    }

	var checkSums []string

    for _, f := range files {
        filePath := destFolder + "/" +f.Name()
        compressFiles := append(otherFilesAbs, filePath)

        var fileExtension = filepath.Ext(filePath)
        var filePathNoExtension = filePath[0:len(filePath)-len(fileExtension)]

        for osName, osExt := range osCompressMapping {
            if strings.Contains(f.Name(), osName) {
                compressedPath := filePathNoExtension + "." + osExt
                switch osExt {
                case "zip":
                    fmt.Println("creating build zip file for " + f.Name() + "...")
                    if err := archiver.Archive(compressFiles, compressedPath); err != nil {
                        return err
                    }
                    if err := os.Remove(filePath); err != nil {
                        return err
                    }
                    break
                case "tar.gz":
                    fmt.Println("creating build tar.gz file for " + f.Name() + "...")
                    if err := archiver.Archive(compressFiles, compressedPath); err != nil {
                        return err
                    }
                    if err := os.Remove(filePath); err != nil {
                        return err
                    }
                    break
                }

                data, err := ioutil.ReadFile(compressedPath)
                if err != nil {
                    return err
                }

                hash := sha256.Sum256(data)
                str := hex.EncodeToString(hash[:])
                checkSums = append(checkSums, str + "  " + compressedPath[strings.LastIndex(
                    compressedPath, "/")+1:])
            }
        }
    }

    err = ioutil.WriteFile(destFolder + "/checksums.txt", []byte(strings.Join(checkSums, "\n")), os.ModePerm)
    if err != nil {
        return err
    }

    return nil
}

func main() {
    currPath, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    config := Config{}
    buff, err := ioutil.ReadFile(currPath + "/config.yml")
    if err != nil {
        log.Fatal(err)
    }

    if err := yaml.Unmarshal(buff, &config); err != nil {
        log.Fatal(err)
    }

    targets := "darwin/*,linux/*,linux/*,windows/*"
    if len(config.Targets) > 0 {
        targets = strings.Join(config.Targets, ",")
    }

    if err := buildBinaries(targets, config.ProjectName, config.SrcFolder,
        config.DestFolder); err != nil {
        log.Fatal(err)
    }

    if config.SrcFolder, err = filepath.Abs(config.SrcFolder); err != nil {
        log.Fatal(err)
    }

    if config.DestFolder, err = filepath.Abs(config.DestFolder); err != nil {
        log.Fatal(err)
    }

    // pack toolchain and zip it
    if err := packAndCompress(config.DestFolder, config.OtherFiles); err != nil {
       log.Fatal(err)
    }

    fmt.Println("Project cross compiled and packaged successfully!")
}
