package core

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Feed represents a podcast feed
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel
}

// Channel represents a feed channel
type Channel struct {
	XMLName               xml.Name            `xml:"channel"`
	Title                 string              `xml:"title"`
	Link                  string              `xml:"link"`
	AtomLink              *AtomLink           `xml:"atom_link,omitempty"`
	Description           string              `xml:"description,omitempty"`
	Language              string              `xml:"language,omitempty"`
	Category              string              `xml:"category,omitempty"`
	Copyright             string              `xml:"copyright,omitempty"`
	Image                 *Image              `xml:"image,omitempty"`
	ItunesAuthor          string              `xml:"itunes_author,omitempty"`
	ItunesOwner           string              `xml:"itunes_owner,omitempty"`
	ItunesImage           *ItunesImage        `xml:"itunes_image,omitempty"`
	ItunesSubtitle        string              `xml:"itunes_subtitle,omitempty"`
	ItunesSummary         string              `xml:"itunes_summary,omitempty"`
	ItunesCategory        *ItunesCategory     `xml:"itunes_category,omitempty"`
	ItunesExplicit        string              `xml:"itunes_explicit,omitempty"`
	GoogleplayAuthor      string              `xml:"googleplay_author,omitempty"`
	GoogleplayImage       *GoogleplayImage    `xml:"googleplay_image,omitempty"`
	GoogleplayEmail       string              `xml:"googleplay_mail,omitempty"`
	GoogleplayDescription string              `xml:"googleplay_description,omitempty"`
	GoogleplayCategory    *GoogleplayCategory `xml:"googleplay_category,omitempty"`
	GoogleplayExplicit    string              `xml:"googleplay_explicit,omitempty"`
	Item                  []Item              `xml:"item"`
}

// AtomLink represents an atom Channel.AtomLink
type AtomLink struct {
	Href string `xml:"href,attr,omitempty"`
}

// Image represents a Channel.Image
type Image struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url,omitempty"`
	Title   string   `xml:"title,omitempty"`
	Link    string   `xml:"link,omitempty"`
	Width   string   `xml:"width,omitempty"`
	Height  string   `xml:"height,omitempty"`
}

// ItunesImage represents à Channel.ItunesImage
type ItunesImage struct {
	Herf string `xml:"href,attr,omitempty"`
}

// ItunesCategory represents à Channel.ItunesCategory
type ItunesCategory struct {
	Herf string `xml:"text,attr,omitempty"`
}

// GoogleplayImage represents à Channel.GoogleplayImage
type GoogleplayImage struct {
	Herf string `xml:"href,attr,omitempty"`
}

// GoogleplayCategory represents à Channel.GoogleplayImage
type GoogleplayCategory struct {
	Herf string `xml:"text,attr,omitempty"`
}

// Item represents a Channel.Item
type Item struct {
	Title                 string     `xml:"title"`
	Link                  string     `xml:"link"`
	Description           string     `xml:"description,omitempty"`
	GUID                  string     `xml:"guid,omitempty"`
	GUIDisPermalink       string     `xml:"isPermalink,attr,omitempty"`
	PubDate               string     `xml:"pubDate,omitempty"`
	Enclosure             *Enclosure `xml:"enclosure,omitempty"`
	ItunesAuthor          string     `xml:"itunes_author,omitempty"`
	ItunesSubtitle        string     `xml:"itunes_subtitle,omitempty"`
	ItunesSummary         string     `xml:"itunes_summary,omitempty"`
	ItunesDuration        int        `xml:"itunes_duration,omitempty"`
	ItunesKeywords        string     `xml:"itunes_keywords,omitempty"`
	ItunesExplicit        string     `xml:"itunes_explicit,omitempty"`
	GoogleplayAuthor      string     `xml:"googleplay_author,omitempty"`
	GoogleplayDescription string     `xml:"googleplay_description,omitempty"`
	GoogleplayExplicit    string     `xml:"googleplay_explicit,omitempty"`
}

// Enclosure represents Channel.Item.Enclosure
type Enclosure struct {
	URL    string `xml:"url,attr,omitempty"`
	Length string `xml:"length,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
}

// std XML lib doesn't handle XML prefixes...
func prefixWorkaround(in []byte) []byte {
	in = bytes.Replace(in, []byte("<atom:"), []byte("<atom_"), -1)
	in = bytes.Replace(in, []byte("</atom:"), []byte("</atom_"), -1)
	in = bytes.Replace(in, []byte("<itunes:"), []byte("<itunes_"), -1)
	in = bytes.Replace(in, []byte("</itunes:"), []byte("</itunes_"), -1)
	in = bytes.Replace(in, []byte("<googleplay:"), []byte("<googleplay_"), -1)
	return bytes.Replace(in, []byte("</googleplay:"), []byte("</googleplay_"), -1)
}

// NewFeed returns a new feed
func NewFeed(url string) (feed Feed, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = xml.Unmarshal(prefixWorkaround(body), &feed); err != nil {
		return
	}
	return
}
