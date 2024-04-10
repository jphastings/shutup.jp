package main

import (
	"fmt"
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

type Date string

func (d Date) Valid() bool {
	_, _, _, err := d.Components()
	return err == nil
}

func (d Date) Components() (year int, month int, day int, err error) {
	_, err = fmt.Sscanf(string(d), "%d-%d-%d", &year, &month, &day)
	return year, month, day, err
}

func (d Date) Time() (time.Time, error) {
	year, month, day, err := d.Components()
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, time.Month(month), day, 23, 0, 0, 0, time.UTC), nil
}

func (d Date) Format(layout string) string {
	t, err := d.Time()
	if err != nil {
		return ""
	}
	return t.Format(layout)
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
func (a BySentOn) Less(i, j int) bool { return a[i].SentOn < a[j].SentOn }
