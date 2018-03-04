package items

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestRepositoryConfigure(t *testing.T) {
	repo := itemRepository{collections: make(map[string][]ItemDefinition)}

	err := repo.Configure(path.Join("this", "is", "a", "fake", "path"))
	if err == nil {
		t.Error("Expected error when trying to configure repo on invalid path")
	}

	tempDir, err := ioutil.TempDir("", "repotests")
	if err != nil {
		t.Fatalf("Failed to make a tempdir")
	}
	err = repo.Configure(tempDir)
	if err != nil {
		t.Error("Directory exists, should not return error")
	}

	os.Create(path.Join(tempDir, "foo.yaml"))
	os.Create(path.Join(tempDir, "bar.yaml"))

	if err = repo.EnsureLoaded("foo", "bar"); err != nil {
		t.Error("Failed to load foo and bar", err)
	}

	if err = repo.EnsureLoaded("baz"); err == nil {
		t.Error("Expected to get an error trying to load non-existed collection baz")
	}
}

func TestRepositoryErrorBeforeConfigured(t *testing.T) {
	repo := itemRepository{collections: make(map[string][]ItemDefinition)}

	_, err := repo.Get("foo")

	if err == nil {
		t.Error("Expected error due to using non-configured repository")
	}

	if !strings.Contains(err.Error(), "Repository must be configured") {
		t.Error("Returned error was incorrect. Expecting \"Repository must be configured...\", instead got", err)
	}
}
