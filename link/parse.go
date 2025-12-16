package link

import "io"

// Link represents a link (<a href="">) in an HTML document
type Link struct {
	Href string
	Text string
}

var r io.Reader

func Parse(r io.Reader) ([]Link, error) {
	return nil, nil
}
