package rexLib

import "embed"

//go:embed "version"
var VersionF embed.FS

func GetVersion() string {
	// note: 读取包中的version文件，并返回
	versionFile, err := VersionF.ReadFile("version")
	if err != nil {
		panic(err)
	}
	return string(versionFile)
}
