package walker

import (
	"os"
	"testing"
)

func Test_Search(t *testing.T) {
	content, err := os.ReadFile("./responses/init.html")
	if err != nil {
		t.Fatalf("Could not read file init.html %v", err)
	}

	anchors := search(string(content), "http://localhost")
	if len(anchors) != 2 {
		t.Errorf("Expecting 3 anchors. Found %d", len(anchors))
	}

	if anchors[0] != "http://localhost/gb" {
		t.Errorf("Expected first anchror to be %s. Found %s", "http://localhost/gb", anchors[0])
	}

	if anchors[1] != "http://localhost/sec" {
		t.Errorf("Expected first anchror to be %s. Found %s", "http://localhost/sec", anchors[1])
	}
}
