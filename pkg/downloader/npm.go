package downloader

import (
	"fmt"
	"os"
	"wio/internal/config/root"
	"wio/pkg/log"
	"wio/pkg/npm/resolve"
	"wio/pkg/util"
	"wio/pkg/util/sys"
)

type NpmDownloader struct{}

func traverseAndDelete(n *resolve.Node) error {
	path := sys.Path(root.GetToolchainPath(), sys.WioFolder, sys.Modules, n.Name+"__"+n.ResolvedVersion.String())

	if err := os.RemoveAll(path); err != nil {
		return err
	}

	for _, dep := range n.Dependencies {
		if err := traverseAndDelete(dep); err != nil {
			return err
		}
	}

	return nil
}

func traverseAndSymlink(n *resolve.Node) error {
	oldPath := sys.Path(root.GetToolchainPath(), sys.WioFolder, sys.Modules, n.Name+"__"+n.ResolvedVersion.String())
	newFilePath := sys.Path(root.GetToolchainPath(), fmt.Sprintf("%s__%s", n.Name, n.ResolvedVersion.String()))

	if sys.Exists(newFilePath) {
		if err := os.RemoveAll(newFilePath); err != nil {
			return err
		}
	}

	if err := os.Symlink(oldPath, newFilePath); err != nil {
		return err
	}

	for _, dep := range n.Dependencies {
		if err := traverseAndSymlink(dep); err != nil {
			return err
		}
	}

	return nil
}

func (npmDownloaded NpmDownloader) DownloadModule(path, name, version string, retool bool) (string, error) {
	info := resolve.NewInfo(path)

	var err error = nil
	if util.IsEmptyString(version) {
		if version, err = info.GetLatest(name); err != nil {
			return "", util.Error("no version found for toolchain %s", name)
		}
	}

	log.Write(log.Cyan, "Fetching toolchain ")
	log.Writeln(log.Green, "%s@%s", name, version)

	node := &resolve.Node{Name: name, ConfigVersion: version}

	if err = info.ResolveTree(node); err != nil {
		return "", err
	}

	pkgPath := sys.Path(root.GetToolchainPath(), sys.WioFolder, sys.Modules, name+"__"+version)

	shouldDownload := false
	if sys.Exists(pkgPath) && retool {
		if err := traverseAndDelete(node); err != nil {
			return "", err
		}

		shouldDownload = true
	} else if !sys.Exists(pkgPath) {
		shouldDownload = true
	}

	log.Write("\033[2K")
	if shouldDownload {
		if err = info.InstallResolved(); err != nil {
			return "", err
		}
	} else {
		log.Writeln(log.Green, "|> Already Exists!")
	}

	if err := traverseAndSymlink(node); err != nil {
		return "", err
	}

	return sys.Path(path, node.Name) + "__" + version, nil
}
