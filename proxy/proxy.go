package proxy

import (
    "fmt"
    "os"
    "time"
    "midiproxy/config"
    "midiproxy/midi"
    "github.com/rakyll/portmidi"
)


func Print_available_devices() {
    midi.Initialize()

    // List the available MIDI input devices
    input_devices := midi.Get_midi_inputs()
    fmt.Println("## Input Devices ##")
    for _, device := range input_devices {
        fmt.Println(device)
    }
    // List the available MIDI output devices
    output_devices := midi.Get_midi_outputs()
    fmt.Println("\n## Output Devices ##")
    for _, device := range output_devices {
        fmt.Println(device)
    }

    // exit program
    midi.Deinitialize()
    os.Exit(0)
}


func Loop() {
    // Getting input and output streams
    midi.Initialize()
    defer end_gracefully()

    inputs, outputs, err := config.Get_midi_streams()
    // Check for errors
    if err != nil {
        panic(err)
    }
    if len(inputs) == 0 {
        panic("No input devices found!")
    }
    if len(outputs) == 0 {
        panic("No input devices found!")
    }
    // Make message channels
    input_channel   := make(chan portmidi.Event)
    output_channels := []chan portmidi.Event{}

    // Make input listening routines
    for _, in_device := range inputs {
        go listen(in_device, input_channel)
    }
    // Make output listening channels and routines
    for _, out_device := range outputs {
        out_chan := make(chan portmidi.Event)
        output_channels = append(output_channels, out_chan)
        go send(out_device, out_chan)
    }

    // Main loop
    fmt.Println("Starting midiproxy...")
    for {
        message := <- input_channel
        for _, out_chan := range output_channels {
            out_chan <- message
        }
    }
}


func listen(device midi.Port, input_channel chan <- portmidi.Event) {
    defer device.Stream.Close()
    for {
        avail, err := device.Stream.Poll()
        if err != nil {
            panic(err)
        }
        if avail {
            events, err := device.Stream.Read(1024)
            if err != nil {
                panic(err)
            }
            for _, event := range events {
                // Check if the event type should be proxied
                event_type := event.Status >> 4
                for _, proxy := range device.Proxy {
                    if event_type == proxy {
                        input_channel <- event
                        break
                    }
                }
            }
        }
        time.Sleep(time.Millisecond * 10)
    }
}

func send(device midi.Port, output_channel chan portmidi.Event) {
    defer device.Stream.Close()
    for {
        // Get message
        event := <- output_channel
        event_type := event.Status >> 4
        // Check if the event tpe should be proxied
        for _, proxy := range device.Proxy {
            if event_type == proxy {
                device.Stream.WriteShort(event.Status, event.Data1, event.Data2)
                break
            }
        }
    }
}



func end_gracefully() {
    if r := recover(); r != nil {
        fmt.Println(r)
    }else{
        fmt.Println("Exiting...")
    }
    midi.Deinitialize()
}
