package main
import (
    "fmt"
    "github.com/jessevdk/go-flags"
    "midiproxy/proxy"
)

type Options struct {
    ListDevices bool `short:"l" description:"List available devices"`
}



func main() {
    var opts Options
    // Parse command-line arguments
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        panic(err)
    }

    if opts.ListDevices {
        proxy.Print_available_devices()
    }
    proxy.Loop()
    fmt.Println("Exiting...")
}
