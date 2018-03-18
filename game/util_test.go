package main

import (
	"log"
	"strings"
	"testing"
)

func TestWrapTextShortString(t *testing.T) {
	content := "Hello world, I am short"

	result := wrapText(content, 0, 0, 40)

	if len(result) != 1 {
		t.Error("Expected only one line of text but got", len(result))
	}

	content = "Short text, to be indented 4 spaces"
	//brb

	result = wrapText(content, 4, 0, 40)

	if len(result) != 1 {
		t.Error("Expected only one line of text but got", len(result))
	}

	if result[0][:5] != "    S" {
		t.Errorf("Expected to see indentation of 4, instead saw >%s<", result[0][:5])
	}
}

func TestWrappingWithIndent(t *testing.T) {

	content := "Short text, to be indented 4 spaces"
	result := wrapText(content, 4, 0, 20)

	log.Printf("result: %v", strings.Join(result, "%"))
	if len(result) != 2 {
		t.Error("Expected two lines of text but got", len(result))
	}

	if result[0][:5] != "    S" {
		t.Errorf("Expected to see indentation of 4, instead saw >%s<", result[0][:5])
	}

	if result[1][:2] != "be" {
		t.Error("Expected second lien to start with 'indented' but instead got", result[1][:8])
	}
}

func TestWrapTextSubsequentIndent(t *testing.T) {
	content := "Short text, to be to be indented 2 spaces with a trailing indent of 6 spaces"
	result := wrapText(content, 2, 6, 20)
	log.Printf("result: %v", strings.Join(result, "%"))

	if len(result) != 7 {
		t.Error("Expected to get 7 lines of text after wrapping but got", len(result))
	}

	if result[0][:3] != "  S" {
		t.Errorf("Expected beginning of first line to be '  S' but instead got '%s'", result[0][:3])
	}

	if result[1][:8] != "      to" {
		t.Errorf("Expected beginning of first line to be '      to' but instead got '%s'", result[1][:8])
	}
}
