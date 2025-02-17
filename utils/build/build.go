package build

import "log"

var (
	currentBuild    BuildType
	availableBuilds []BuildType = []BuildType{Debug, Release}
)

const (
	Debug BuildType = iota
	Release
)

type BuildType int

func SetBuildType(build BuildType) {
	for _, availableBuild := range availableBuilds {
		if build == availableBuild {
			return
		}
	}

	log.Fatal("Error! Unavailable build")
}

func GetBuildType() BuildType {
	return currentBuild
}
