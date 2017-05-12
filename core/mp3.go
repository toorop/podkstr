package core

import (
	"os"
	"time"

	"fmt"

	"github.com/tcolgate/mp3"
)

// Mp3 represents an mp3  audio file
type Mp3 struct {
	Path string
}

// NewMp3 return a new Mp3 struct
func NewMp3(path string) (mp3 Mp3, err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return mp3, fmt.Errorf("File %s doesn't exists", path)
	}
	mp3.Path = path
	return
}

// GetDuration return audi duration as time.Duration
func (m *Mp3) GetDuration() (duration time.Duration, err error) {
	skipped := 0
	r, err := os.Open(m.Path)
	if err != nil {
		return
	}
	d := mp3.NewDecoder(r)
	var f mp3.Frame
	for {
		if err := d.Decode(&f, &skipped); err != nil {
			return duration, err
		}
		duration += f.Duration()
	}
}
