package midi

import (
    "github.com/rakyll/portmidi"
)


var (
    Proxy_Status_Codes = map[string]int64{"note_on": 9, "note_off": 8, "program_change": 12}
)


type Port struct {
    Name      string
    Stream    *portmidi.Stream
    Proxy     []int64
}
