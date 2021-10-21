package testutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/onsi/ginkgo"
)

func GinkgoGetHorusecCmd(horusecCmd string) *exec.Cmd {
	bin := ginkgoBuildHorusecBinary()
	args := setLogLevelArgsToHorusecCmd(horusecCmd)
	return exec.Command(bin, args...)
}

func GinkgoGetHorusecCmdWithFlags(cmdArg string, flags map[string]string) *exec.Cmd {
	bin := ginkgoBuildHorusecBinary()
	args := setLogLevelArgsToHorusecCmd(cmdArg)
	for flag, value := range flags {
		args = append(args, fmt.Sprintf("%s=%s", flag, value))
	}
	return exec.Command(bin, args...)
}

func ginkgoBuildHorusecBinary(customArgs ...string) string {
	binary := filepath.Join(os.TempDir(), getBinaryNameBySystem())
	args := []string{
		"build",
		`-ldflags=-X 'github.com/ZupIT/horusec/cmd/app/version.Version=vTest'`,
		fmt.Sprintf("-o=%s", binary), filepath.Join(RootPath, "cmd", "app"),
	}
	args = append(args, customArgs...)
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Error on build Horusec binary for e2e test %v", err))
	}
	return binary
}

func setLogLevelArgsToHorusecCmd(horusecCmd ...string) []string {
	return append(horusecCmd, fmt.Sprintf("%s=%s", "--log-level", "debug"))
}

func getBinaryNameBySystem() string {
	binaryName := "horusec-e2e"
	if runtime.GOOS == "windows" {
		binaryName = fmt.Sprintf("%s.exe", binaryName)
	}
	return binaryName
}
