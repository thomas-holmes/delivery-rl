package items

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestRepositoryConfigure(t *testing.T) {
	repo := NewRepository()

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
	repo := NewRepository()

	err := repo.EnsureLoaded("foo")
	if err == nil {
		t.Error("Expected error due to using non-configured repository")
	}

	if !strings.Contains(err.Error(), "Repository must be configured") {
		t.Error("Returned error was incorrect. Expecting \"Repository must be configured...\", instead got", err)
	}
}

func TestRepositoryLoadsDefinitions(t *testing.T) {
	repo := NewRepository()

	appleYAML := `---
 - name: "Apple"
   description: "A delicious green apple"
   glyph: "a"
   color: [0, 255, 0]
   kind: consumeable

 - name: Chocolate Cake
   description: "Yummy chocolate cake"
   glyph: "c"
   color: [150, 150, 50]
   kind: consumeable
`

	tempDir, err := ioutil.TempDir("", "testYamlLoad")
	if err != nil {
		t.Fatal("Failed to create a tempdir", err)
	}

	if err := ioutil.WriteFile(path.Join(tempDir, "consumeables.yaml"), []byte(appleYAML), 0666); err != nil {
		t.Fatal("Failed to write test file", err)
	}
	t.Log("Wrote file", path.Join(tempDir, "consumeables.yaml"))

	if err := repo.Configure(tempDir); err != nil {
		t.Fatal("Failed to configure repo", err)
	}

	collection := repo.Get("consumeables")

	if len(collection.definitions) != 2 {
		t.Error("Expected size of collection to be 2, instead got", len(collection.definitions))
	}

	def, ok := collection.GetByName("chocolate cake")
	if !ok {
		t.Error("Expected to find chocolate cake")
	}

	if def.Name != "Chocolate Cake" {
		t.Error("Expected to retrieve definition for Chocolate Cake, instead got", def.Name)
	}
}
