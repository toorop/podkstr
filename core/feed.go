package core

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
)

// ErrFeedNotAPodcast when feed is not related to a podcast
var ErrFeedNotAPodcast = errors.New("this feed is not related to a podcast")

// Feed represents a podcast feed
type Feed struct {
	XMLName     xml.Name `xml:"rss"`
	Version     string   `xml:"version,attr"`
	XMLnsItunes string   `xml:"xmlns_itunes,attr"`
	Channel     Channel
}

// Channel represents a feed channel
type Channel struct {
	XMLName               xml.Name           `xml:"channel"`
	Title                 string             `xml:"title"`
	LastBuildDate         string             `xml:"lastBuildDate"`
	Link                  string             `xml:"link"`
	AtomLink              AtomLink           `xml:"atom_link"`
	Description           string             `xml:"description,omitempty"`
	Language              string             `xml:"language,omitempty"`
	Category              string             `xml:"category,omitempty"`
	Copyright             string             `xml:"copyright,omitempty"`
	Image                 FeedImage          `xml:"image"`
	ItunesAuthor          string             `xml:"itunes_author,omitempty"`
	ItunesOwner           ItunesOwner        `xml:"itunes_owner,omitempty"`
	ItunesImage           ItunesImage        `xml:"itunes_image,omitempty"`
	ItunesSubtitle        string             `xml:"itunes_subtitle,omitempty"`
	ItunesSummary         string             `xml:"itunes_summary,omitempty"`
	ItunesCategory        ItunesCategory     `xml:"itunes_category"`
	ItunesExplicit        string             `xml:"itunes_explicit,omitempty"`
	GoogleplayAuthor      string             `xml:"googleplay_author,omitempty"`
	GoogleplayImage       GoogleplayImage    `xml:"googleplay_image"`
	GoogleplayEmail       string             `xml:"googleplay_mail,omitempty"`
	GoogleplayDescription string             `xml:"googleplay_description,omitempty"`
	GoogleplayCategory    GoogleplayCategory `xml:"googleplay_category,omitempty"`
	GoogleplayExplicit    string             `xml:"googleplay_explicit,omitempty"`
	Items                 []Item             `xml:"item"`
}

// AtomLink represents an atom Channel.AtomLink
type AtomLink struct {
	XMLName xml.Name `xml:"atom_link"`
	Href    string   `xml:"href,attr,omitempty"`
}

// FeedImage represents a Channel.Image
type FeedImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url,omitempty"`
	Title   string   `xml:"title,omitempty"`
	Link    string   `xml:"link,omitempty"`
	Width   string   `xml:"width,omitempty"`
	Height  string   `xml:"height,omitempty"`
}

// ItunesOwner represents a Channel.ItunesOwner
type ItunesOwner struct {
	XMLName xml.Name `xml:"itunes_owner"`
	Name    string   `xml:"itunes_name,omitempty"`
	Email   string   `xml:"itunes_email,omitempty"`
}

// ItunesImage represents à Channel.ItunesImage
type ItunesImage struct {
	XMLName xml.Name `xml:"itunes_image"`
	Href    string   `xml:"href,attr,omitempty"`
}

// ItunesCategory represents à Channel.ItunesCategory
type ItunesCategory struct {
	XMLName xml.Name `xml:"itunes_category"`
	Href    string   `xml:"text,attr,omitempty"`
}

// GoogleplayImage represents à Channel.GoogleplayImage
type GoogleplayImage struct {
	XMLName xml.Name `xml:"googleplay_image"`
	Href    string   `xml:"href,attr,omitempty"`
}

// GoogleplayCategory represents à Channel.GoogleplayImage
type GoogleplayCategory struct {
	XMLName xml.Name `xml:"googleplay_category"`
	Href    string   `xml:"text,attr,omitempty"`
}

// Item represents a Channel.Item
type Item struct {
	Title                 string        `xml:"title"`
	Link                  string        `xml:"link"`
	Description           string        `xml:"description,omitempty"`
	GUID                  string        `xml:"guid,omitempty"`
	GUIDisPermalink       string        `xml:"isPermaLink,attr,omitempty"`
	PubDate               string        `xml:"pubDate,omitempty"`
	Enclosure             ItemEnclosure `xml:"enclosure,omitempty"`
	Image                 ItemImage     `xml:"image"`
	ItunesImage           string        `xml:"itunes_image,omitempty"`
	ItunesAuthor          string        `xml:"itunes_author,omitempty"`
	ItunesSubtitle        string        `xml:"itunes_subtitle,omitempty"`
	ItunesSummary         string        `xml:"itunes_summary,omitempty"`
	ItunesDuration        string        `xml:"itunes_duration,omitempty"`
	ItunesKeywords        string        `xml:"itunes_keywords,omitempty"`
	ItunesExplicit        string        `xml:"itunes_explicit,omitempty"`
	GoogleplayAuthor      string        `xml:"googleplay_author,omitempty"`
	GoogleplayDescription string        `xml:"googleplay_description,omitempty"`
	GoogleplayExplicit    string        `xml:"googleplay_explicit,omitempty"`
}

// ItemImage represents Channel.Item.Image
type ItemImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url,omitempty"`
	Title   string   `xml:"title,omitempty"`
	Link    string   `xml:"link,omitempty"`
}

// ItemEnclosure represents Channel.Item.Enclosure
type ItemEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr,omitempty"`
	Length  string   `xml:"length,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
}

// std XML lib doesn't handle XML prefixes...
func prefixWorkaround(in []byte) []byte {
	in = bytes.Replace(in, []byte("<atom:"), []byte("<atom_"), -1)
	in = bytes.Replace(in, []byte("</atom:"), []byte("</atom_"), -1)
	in = bytes.Replace(in, []byte("<itunes:"), []byte("<itunes_"), -1)
	in = bytes.Replace(in, []byte("</itunes:"), []byte("</itunes_"), -1)
	in = bytes.Replace(in, []byte("<googleplay:"), []byte("<googleplay_"), -1)
	in = bytes.Replace(in, []byte("xmlns:"), []byte("xmlns_"), -1)
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
	if feed.XMLnsItunes == "" {
		return Feed{}, ErrFeedNotAPodcast
	}
	return
}
