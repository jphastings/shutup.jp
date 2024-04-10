package main

import (
	"math/big"
	"time"
)

type Postcard struct {
	Name            string
	Location        Location `json:"location,omitempty"`
	Flip            Flip     `json:"flip" yaml:"flip"`
	SentOn          Date     `json:"sentOn,omitempty" yaml:"sent_on"`
	Front           Side     `json:"front,omitempty"`
	Back            Side     `json:"back,omitempty"`
	FrontDimensions Size     `json:"frontSize" yaml:"front_size,omitempty"`
	Context         Context  `json:"context,omitempty"`
}

type Location struct {
	Name      string   `json:"name"`
	Latitude  *float64 `json:"lat,omitempty" yaml:"latitude,omitempty"`
	Longitude *float64 `json:"long,omitempty" yaml:"longitude,omitempty"`
}

type Side struct {
	Description   string `json:"description,omitempty"`
	Transcription string `json:"transcription,omitempty"`
}

type Context struct {
	Description string `json:"description"`
}

type Flip string

const (
	FlipBook      Flip = "book"
	FlipLeftHand  Flip = "left-hand"
	FlipCalendar  Flip = "calendar"
	FlipRightHand Flip = "right-hand"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	str := string(b)
	if str == "null" {
		d.Time = time.Time{}
		return nil
	}

	// Specify your custom layout here
	d.Time, err = time.Parse(`"2006-01-02"`, str)
	return err
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(d.Time.Format(`"2006-01-02"`)), nil
}

type Size struct {
	CmWidth  *big.Rat `json:"cmW,omitempty"`
	CmHeight *big.Rat `json:"cmH,omitempty"`
	PxWidth  int      `json:"pxW"`
	PxHeight int      `json:"pxH"`
}

type BySentOn []Postcard

// Implement the sort.Interface - Len, Less, and Swap methods
func (a BySentOn) Len() int           { return len(a) }
func (a BySentOn) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySentOn) Less(i, j int) bool { return a[i].SentOn.Unix() > a[j].SentOn.Unix() }
