package config

import (
    "errors"
    "fmt"
    "os"
    "gopkg.in/yaml.v3"
    "path/filepath"
    "midiproxy/midi"
)

const (
    config_file = "config.yaml"
)

// Read the physical config file and return the data
func read_config() (map[string]interface{}, error) {
    currentDir, err := os.Getwd()
    if err != nil {
        fmt.Println("Failed to get current working directory:", err)
        return nil, err
    }
    // Read file
    file_path := filepath.Join(currentDir, config_file)

    file, err := os.Open(file_path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Parse the YAML file
    decoder := yaml.NewDecoder(file)
    var data map[string]interface{}
    if err := decoder.Decode(&data); err != nil {
        return nil, err
    }

    // Validate mandatory input/output
    _, ok := data["inputs"]
    if !ok {
        return nil, errors.New("Invalid config! No inputs defined")
    }
    _, ok = data["outputs"]
    if !ok {
        return nil, errors.New("Invalid config! No outputs defined")
    }
    // Return data
    return data, nil
}


// Parse the config data, and return inputs and outputs as port ojbects
func parse_conf(data map[string]interface{}) ([]midi.Port, []midi.Port, error) {
    // Inputs
    raw_devices, ok := data["inputs"].([]interface{})
    if !ok {
        return nil, nil, errors.New("Could not parse the config file! (inputs)")
    }
    inputs, err := parse_device_conf(raw_devices)
    if err != nil {
        return nil, nil, err
    }
    // Outputs
    raw_devices, ok = data["outputs"].([]interface{})
    if !ok {
        return nil, nil, errors.New("Could not parse the config file! (outputs)")
    }
    outputs, err := parse_device_conf(raw_devices)
    if err != nil {
        return nil, nil, err
    }
    return inputs, outputs, nil
}


// Parse a list of device config and return a list of Port types
func parse_device_conf(raw_config []interface{}) ([]midi.Port, error) {
    var devices []midi.Port
    for _, raw_port := range raw_config {
        port, ok := raw_port.(map[string]interface{})
        if !ok {
            return nil, errors.New("Could not parse the config file! (inputs)")
        }
        raw_name, ok := port["name"]
        if !ok {
            return nil, errors.New("Could not parse the config file! (inputs)")
        }
        name, ok := raw_name.(string)
        if !ok {
            return nil, errors.New("Could not parse the config file! (inputs)")
        }
        // Convert proxy to list of strings
        raw_proxy, ok := port["proxy"]
        str_proxy := get_map_keys(midi.Proxy_Status_Codes)
        if ok {
            if proxySlice, ok := raw_proxy.([]interface{}); ok {
                str_proxy = convert_to_string_slice(proxySlice)
            }
        }
        // Turn the strings into their status codes
        var proxy []int64
        for _, i := range str_proxy {
            proxy = append(proxy, midi.Proxy_Status_Codes[i])
        }
        device := midi.Port{Name: name, Proxy: proxy}
        devices = append(devices, device)
    }
    return devices, nil
}


// Get the midi streams for input and output
func Get_midi_streams() ([]midi.Port, []midi.Port, error){
    var inputs  []midi.Port
    var outputs []midi.Port

    // Read config
    conf, err := read_config()
    raw_inputs, raw_outputs, err := parse_conf(conf)
    if err != nil {
        return nil, nil, err
    }
    // Open input ports for devices
    for _, device := range raw_inputs{
        name := device.Name
        // Some devices may have the same name, and open both streams
        for _, stream_pointer := range midi.Open_midi_input(name) {
            stream_dev := midi.Port{Name: name, Stream: stream_pointer, Proxy: device.Proxy}
            inputs = append(inputs, stream_dev)
        }
    }
    // Open output ports for devices
    for _, device := range raw_outputs{
        name := device.Name
        // Some devices may have the same name, and open both streams
        for _, stream_pointer := range midi.Open_midi_output(name) {
            stream_dev := midi.Port{Name: name, Stream: stream_pointer, Proxy: device.Proxy}
            outputs = append(outputs, stream_dev)
        }
    }
    return inputs, outputs, nil
}
