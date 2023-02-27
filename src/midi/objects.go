package midi

import (
    "github.com/rakyll/portmidi"
)


var (
    Proxy_Status_Codes = map[string]int64{
        "note_off": 8,
        "note_on": 9,
        "poly_aftertouch": 10,
        "control_change": 11,
        "program_change": 12,
        "chan_aftertouch": 13,
        "pitch_bend": 14}
)


type Port struct {
    Name      string
    Stream    *portmidi.Stream
    Proxy     []int64
}
