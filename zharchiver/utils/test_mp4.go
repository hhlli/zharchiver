package utils

import (
	"fmt"
	"github.com/abema/go-mp4"
)

func TestPrint() {
	var t *mp4.TrackInfo
    fmt.Println(t.AVC)
}
