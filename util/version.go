package util

import (
	"fmt"
	"runtime"
)

var(
	version = "99.99.99" //value from version file, inherited while building using makefile
	buildDate = "1970-01-01T00:00:00Z" //from `date -u +'%Y-%m-%dT%H:%M:%SZ'`
)
type Version struct{
	Version string
	BuildDate string
	Platform string
}
func (v Version) String() string{
	return v.Version
}
func PrintVersion(){
	version:=GetVersion()
	fmt.Println("Version: ", version.Version)
	fmt.Println("BuildDate: ", version.BuildDate)
	fmt.Println("Platform: ", version.Platform)
}
func GetVersion() Version{
	var versionStr string
	versionStr="v"+version
	return Version{
		Version: versionStr,
		BuildDate: buildDate,
		Platform: fmt.Sprintf("%s/%s",runtime.GOOS, runtime.GOARCH),
	}
}
