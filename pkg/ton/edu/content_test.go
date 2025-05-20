package edu

import (
	"testing"
)

func TestContentFromCell_Offchain(t *testing.T) {
	off := ContentOffchain{
		URI: "https://tonutils.com/22.json",
	}

	c, err := off.ContentCell()
	if err != nil {
		t.Fatal(err)
	}

	content, err := ContentFromCell(c)
	if err != nil {
		t.Fatal(err)
	}

	if content.URI != off.URI {
		t.Fatal("URI not eq:", content.URI)
	}
}
