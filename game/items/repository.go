package items

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/MichaelTJones/pcg"
)

type Collection struct {
	definitions []Definition
	nameMap     map[string]int
}

func (c Collection) Sample(rng *pcg.PCG64) Definition {
	i := rng.Bounded(uint64(len(c.definitions)))

	return c.definitions[i]
}

func (c Collection) GetByName(name string) (Definition, bool) {
	d, ok := c.nameMap[strings.ToLower(name)]
	if ok {
		return c.definitions[d], true
	}

	return Definition{}, false
}

type Repository interface {
	Configure(loadPath string) error
	EnsureLoaded(Collection ...string) error
	Get(collectionName string) (Collection, error)
}

var defaultRepository = NewRepository()

func NewRepository() Repository {
	return &itemRepository{collections: make(map[string]Collection)}
}

type itemRepository struct {
	loadPath   string
	configured bool

	collections map[string]Collection
}

func (i *itemRepository) Configure(loadPath string) error {
	fs, err := os.Stat(loadPath)
	if err != nil && os.IsNotExist(err) {
		return errors.New("File path does not exist")
	}

	if !fs.IsDir() {
		return errors.New("Load path must be a directory")
	}

	i.loadPath = loadPath
	i.configured = true

	return nil
}

func (i *itemRepository) ensureConfigured() error {
	if !i.configured {
		return errors.New("Repository must be configured before use")
	}

	return nil
}

// EnsureLoaded eagerly loads all specified collections, returning an error if any fail to load.
func (i *itemRepository) EnsureLoaded(collections ...string) error {
	if err := i.ensureConfigured(); err != nil {
		return err
	}

	for _, c := range collections {
		if err := i.loadIfAbsent(c); err != nil {
			return err
		}
	}

	return nil
}

func (i *itemRepository) loadIfAbsent(collectionName string) error {
	if err := i.ensureConfigured(); err != nil {
		return err
	}

	_, ok := i.collections[collectionName]
	if !ok {
		definitions, err := LoadDefinitions(path.Join(i.loadPath, fmt.Sprintf("%s.yaml", collectionName)))
		if err != nil {
			return err
		}

		nameMap := make(map[string]int)
		for i, d := range definitions {
			nameMap[strings.ToLower(d.Name)] = i
		}

		collection := Collection{
			definitions: definitions,
			nameMap:     nameMap,
		}

		i.collections[collectionName] = collection
	}

	return nil
}

func (i *itemRepository) Get(collectionName string) (Collection, error) {
	if err := i.ensureConfigured(); err != nil {
		return Collection{}, err
	}

	if err := i.loadIfAbsent(collectionName); err != nil {
		return Collection{}, err
	}

	collection, ok := i.collections[collectionName]
	if !ok {
		return Collection{}, errors.New("This should never happen, but just in case we couldn't load the collection immediately after load")
	}

	return collection, nil
}

// Configure operates on the default Repository. Sets the path for loading collections
func Configure(loadPath string) error {
	return defaultRepository.Configure(loadPath)
}

// EnsureLoaded operates on the default Repository. Eagerly load one or more collections, return an error if any fail to load
func EnsureLoaded(collections ...string) error {
	return defaultRepository.EnsureLoaded(collections...)
}

func GetCollection(collectionName string) (Collection, error) {
	return defaultRepository.Get(collectionName)
}
