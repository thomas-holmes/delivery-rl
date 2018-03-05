package monsters

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type Repository interface {
	Configure(loadPath string) error
	EnsureLoaded(collections ...string) error
	GetAll() ([]Definition, error)
	Get(collectionName string) ([]Definition, error)
}

var defaultRepository = NewRepository()

func NewRepository() Repository {
	return &monsterRepository{
		collections: make(map[string][]Definition),
	}
}

type monsterRepository struct {
	loadPath   string
	configured bool

	allDefinitions []Definition
	collections    map[string][]Definition
}

func (i *monsterRepository) Configure(loadPath string) error {
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

func (i *monsterRepository) ensureConfigured() error {
	if !i.configured {
		return errors.New("Repository must be configured before use")
	}

	return nil
}

// EnsureLoaded eagerly loads all specified collections, returning an error if any fail to load.
func (i *monsterRepository) EnsureLoaded(collections ...string) error {
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

func (i *monsterRepository) loadIfAbsent(collectionName string) error {
	if err := i.ensureConfigured(); err != nil {
		return err
	}

	_, ok := i.collections[collectionName]
	if !ok {
		collection, err := LoadDefinitions(path.Join(i.loadPath, fmt.Sprintf("%s.yaml", collectionName)))
		if err != nil {
			return err
		}
		i.allDefinitions = nil
		i.collections[collectionName] = collection
	}

	return nil
}

func (i *monsterRepository) GetAll() ([]Definition, error) {
	if err := i.ensureConfigured(); err != nil {
		return nil, err
	}

	if i.allDefinitions != nil {
		return i.allDefinitions, nil
	}

	var allDefinitions []Definition
	for _, collection := range i.collections {
		for _, definition := range collection {
			allDefinitions = append(allDefinitions, definition)
		}
	}

	i.allDefinitions = allDefinitions

	return allDefinitions, nil
}

func (i *monsterRepository) Get(collectionName string) ([]Definition, error) {
	if err := i.ensureConfigured(); err != nil {
		return nil, err
	}

	if err := i.loadIfAbsent(collectionName); err != nil {
		return nil, err
	}

	collection, ok := i.collections[collectionName]
	if !ok {
		return nil, errors.New("This should never happen, but just in case we couldn't load the collection immediately after load")
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

// GetAll operates on the default Repository. Returns a merged list of all definitions.
func GetAllCollections() ([]Definition, error) {
	return defaultRepository.GetAll()
}

func GetCollection(collectionName string) ([]Definition, error) {
	return defaultRepository.Get(collectionName)
}
