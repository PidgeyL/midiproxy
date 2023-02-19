package midi

import (
    "fmt"
    "os"
    "github.com/rakyll/portmidi"
)


func Initialize() {
    err := portmidi.Initialize()
    if err != nil {
        fmt.Println("Error initializing portmidi:", err)
        os.Exit(1)
    }
}


func Deinitialize() {
    portmidi.Terminate()
}


func Get_midi_inputs() []string {
    numDevices := portmidi.CountDevices()
    devices := make([]string, 0, numDevices)

    for i := 0; i < numDevices; i++ {
        info := portmidi.Info(portmidi.DeviceID(i))
        if info.IsInputAvailable {
            devices = append(devices, info.Name)
        }
    }
    return devices
}


func Get_midi_outputs() []string {
    numDevices := portmidi.CountDevices()
    devices := make([]string, 0, numDevices)

    for i := 0; i < numDevices; i++ {
        info := portmidi.Info(portmidi.DeviceID(i))
        if info.IsOutputAvailable {
            devices = append(devices, info.Name)
        }
    }
    return devices
}


func Open_midi_input(name string) []*portmidi.Stream {
    numDevices := portmidi.CountDevices()
    devices := make([]*portmidi.Stream, 0, numDevices)

    for i := 0; i < numDevices; i++ {
        device_id := portmidi.DeviceID(i)
        info      := portmidi.Info(device_id)
        if info.Name == name && info.IsInputAvailable {
            device, err := portmidi.NewInputStream(device_id, 1024)
            if err != nil {
                fmt.Println("Couldn't attach", info.Name)
                continue
            }
            devices = append(devices, device)
        }
    }
    if len(devices) == 0 {
        fmt.Println(name, "not detected")
    }
    return devices
}


func Open_midi_output(name string) []*portmidi.Stream {
    numDevices := portmidi.CountDevices()
    devices := make([]*portmidi.Stream, 0, numDevices)

    for i := 0; i < numDevices; i++ {
        device_id := portmidi.DeviceID(i)
        info      := portmidi.Info(device_id)
        if info.Name == name && info.IsOutputAvailable {
            device, err := portmidi.NewOutputStream(device_id, 1024, 0)
            if err != nil {
                fmt.Println("Couldn't attach", info.Name)
                continue
            }
            devices = append(devices, device)
        }
    }
    if len(devices) == 0 {
        fmt.Println(name, "not detected")
    }
    return devices
}
