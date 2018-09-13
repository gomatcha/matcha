// Package matcha provides the implementation for the matcha command line tool.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gomatcha/matcha/cmd"
	"github.com/spf13/cobra"
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var RootCmd = &cobra.Command{
	Use:   "matcha",
	Short: "Matcha is a tool for building Android & iOS apps in Go",
	Long:  "Matcha is a tool for building Android & iOS apps in Go",
}

var (
	buildN       bool   // -n
	buildX       bool   // -x
	buildV       bool   // -v
	buildWork    bool   // -work
	buildGcflags string // -gcflags
	buildLdflags string // -ldflags
	buildO       string // -o
	// buildThreaded bool
	// buildBinary  bool   // -binary
	buildTargets string // --targets
)

func init() {
	flags := InitCmd.Flags()
	flags.BoolVarP(&buildN, "dry-run", "n", false, "print the commands but do not run them.")
	flags.BoolVarP(&buildX, "trace", "x", false, "print the commands.")
	flags.BoolVarP(&buildV, "verbose", "v", false, "print the logs verbosely.")
	flags.BoolVar(&buildWork, "work", false, "print the name of the temporary work directory and do not delete it when exiting.")
	flags.StringVar(&buildGcflags, "gcflags", "", "arguments to pass on each go tool compile invocation.")
	flags.StringVar(&buildLdflags, "ldflags", "", "arguments to pass on each go tool link invocation.")
	flags.StringVar(&buildTargets, "target", "", "space separated os/arch. Valid values are: android, ios, android/arm, android/arm64, android/386, android/amd64, ios/arm, ios/arm64, ios/386, ios/amd64.")

	RootCmd.AddCommand(InitCmd)
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Install the mobile compiler toolchain",
	Long:  ``,
	Run: func(command *cobra.Command, args []string) {
		flags := &cmd.Flags{
			Logger:       log.New(os.Stderr, "", 0),
			BuildN:       buildN,
			BuildX:       buildX,
			BuildV:       buildV,
			BuildWork:    buildWork,
			BuildGcflags: buildGcflags,
			BuildLdflags: buildLdflags,
			BuildTargets: buildTargets,
			Threaded:     true,
		}
		if err := cmd.Init(flags); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	flags := BuildCmd.Flags()
	flags.BoolVarP(&buildN, "dry-run", "n", false, "print the commands but do not run them.")
	flags.BoolVarP(&buildX, "trace", "x", false, "print the commands.")
	flags.BoolVarP(&buildV, "verbose", "v", false, "print the logs verbosely.")
	flags.BoolVar(&buildWork, "work", false, "print the name of the temporary work directory and do not delete it when exiting.")
	flags.StringVar(&buildGcflags, "gcflags", "", "arguments to pass on each go tool compile invocation.")
	flags.StringVar(&buildLdflags, "ldflags", "", "arguments to pass on each go tool link invocation.")
	flags.StringVar(&buildTargets, "target", "", "space separated os/arch. Valid values are: android, ios, android/arm, android/arm64, android/386, android/amd64, ios/arm, ios/arm64, ios/386, ios/amd64.")

	RootCmd.AddCommand(BuildCmd)
}

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the Matcha static library",
	Long:  ``,
	Run: func(command *cobra.Command, args []string) {
		flags := &cmd.Flags{
			Logger:       log.New(os.Stderr, "", 0),
			BuildN:       buildN,
			BuildX:       buildX,
			BuildV:       buildV,
			BuildWork:    buildWork,
			BuildGcflags: buildGcflags,
			BuildLdflags: buildLdflags,
			BuildTargets: buildTargets,
			Threaded:     true,
		}
		if err := cmd.Build(flags, args); err != nil {
			fmt.Println(err)
		}
	},
}

/*
func init() {
	flags := InstallCmd.Flags()
	flags.BoolVar(&buildN, "n", false, "print the commands but do not run them.")
	flags.BoolVar(&buildX, "x", false, "print the commands.")
	flags.BoolVar(&buildV, "v", false, "print the names of packages as they are compiled.")
	flags.BoolVar(&buildBinary, "b", false, "builds only the binary.")
	flags.BoolVar(&buildWork, "work", false, "print the name of the temporary work directory and do not delete it when exiting.")
	flags.StringVar(&buildGcflags, "gcflags", "", "arguments to pass on each go tool compile invocation.")
	flags.StringVar(&buildLdflags, "ldflags", "", "arguments to pass on each go tool link invocation.")
	flags.StringVar(&buildO, "output", "", "forces build to write the resulting object to the named output file.")
	flags.StringVar(&buildTargets, "targets", "", "space separated os/arch. Valid values are: android, darwin, android/arm, android/arm64, android/386, android/amd64, darwin/arm, darwin/arm64, darwin/386, darwin/amd64.")

	RootCmd.AddCommand(InstallCmd)
}

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Builds the Matcha static library and copies iOS frameworks to a directory",
	Long:  ``,
	Run: func(command *cobra.Command, args []string) {
		flags := &cmd.Flags{
			BuildN:       buildN,
			BuildX:       buildX,
			BuildV:       buildV,
			BuildWork:    buildWork,
			BuildGcflags: buildGcflags,
			BuildLdflags: buildLdflags,
			BuildO:       buildO,
			BuildBinary:  buildBinary,
			BuildTargets: buildTargets,
		}
		if err := cmd.Bind(flags, args); err != nil {
			fmt.Println(err)
		}
	},
}
*/
