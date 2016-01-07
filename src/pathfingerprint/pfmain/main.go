package main

import (
    "os"
    "fmt"
    
    "runtime/pprof"

    "code.google.com/p/getopt"

    "pathfingerprint/pfinternal"
)

const (
    ProfilerOutputFilename = "manifest.prof"
)

//var doProfile = getopt.BoolLong("profile", 'p', "Enable profiling")

var opts = getopt.CommandLine

func main() {
// TODO(dustin): Still debugging argument parsing (this interferes with 
//               mandatory arguments).
    var doProfile bool = false

    l := pfinternal.NewLogger()

    if opts.Parse(os.Args); opts.NArgs() < 2 {
        fmt.Println("Please provide at least a scan-path and a catalog-path.")
        os.Exit(1)
    }

    opts.Parse(opts.Args())

    scanPath := os.Args[1]
    catalogPath := os.Args[2]

    if doProfile {
        fmt.Println("Profiling enabled.")

        f, err := os.Create(ProfilerOutputFilename)
        if err != nil {
            l.DieIf(err, "Could not create profiler profile.")
        }

        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    p := pfinternal.NewPath()

    c, err := pfinternal.NewCatalog(&catalogPath, &scanPath, true)
    if err != nil {
        fmt.Printf("Could not open catalog: %s\n", err.Error())
        os.Exit(1)
    }

    hash, err := p.GeneratePathHash(&scanPath, c)
    if err != nil {
        fmt.Printf("Could not generate hash: %s\n", err.Error())
        os.Exit(2)
    }

    fmt.Printf("%s\n", hash)
}