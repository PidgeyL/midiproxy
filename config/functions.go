package config

import (
    "fmt"
)


// interface list -> string list
func  convert_to_string_slice(data []interface{}) []string {
    result := make([]string, len(data))
    for i, v := range data {
        result[i] = fmt.Sprintf("%v", v)
    }
    return result
}


// map -> list of string list
func get_map_keys(my_map map[string]int64) ([]string) {
    keys := make([]string, 0, len(my_map))
    for k := range my_map {
        keys = append(keys, k)
    }
    return keys
}
