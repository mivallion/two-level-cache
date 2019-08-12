package two_level_cache

import (
	"encoding/gob"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path"
)

type MemoryCache struct {
	*filenamesMap
	*frequenciesMap
}

func NewMemoryCache() (*MemoryCache, error) {
	m := &MemoryCache{filenamesMap: NewFilenamesMap(), frequenciesMap: NewFrequenciesMap()}
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		err := os.Mkdir("temp", os.ModeDir)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (c *MemoryCache) Cache(key string, value Marshaler) error {
	filename, err := uuid.NewV4()
	if err != nil {
		return err
	}
	filenameString := filename.String()
	err = writeGob("./temp/"+filenameString, value)
	if err != nil {
		return err
	}
	c.filenamesMap.Set(key, filenameString)
	c.SetFrequency(key, 1)
	return nil
}

func (c *MemoryCache) Get(key string, marshaler Marshaler) (Marshaler, bool) {
	freq, ok := c.GetFrequency(key)
	if !ok {
		return marshaler, false
	}
	c.SetFrequency(key, freq+1)
	filename, ok := c.filenamesMap.Get(key)
	if !ok {
		return marshaler, ok
	}
	err := readGob("./temp/"+filename, marshaler)
	if err != nil {
		return marshaler, false
	}
	return marshaler, ok
}

func (c *MemoryCache) Remove(key string) error {
	filename, ok := c.filenamesMap.Get(key)
	if !ok {
		return nil
	}
	err := os.Remove("temp/" + filename + ".gob")
	if err != nil {
		return err
	}
	c.RemoveFrequency(key)
	c.filenamesMap.Remove(key)
	return nil
}

func (c *MemoryCache) Clear() error {
	dir, err := ioutil.ReadDir("temp")
	if err != nil {
		return err
	}
	for _, d := range dir {
		err := os.RemoveAll(path.Join([]string{"temp", d.Name()}...))
		if err != nil {
			return err
		}
	}
	c.frequenciesMap.Clear()
	return nil
}

func (c *MemoryCache) Contains(key string) bool {
	_, ok := c.GetFrequency(key)
	return ok
}

func (c *MemoryCache) Size() int {
	return len(c.filenames)
}

func (c *MemoryCache) Keys() []string {
	var keys []string
	for k := range c.filenames {
		keys = append(keys, k)
	}
	return keys
}

func writeGob(filePath string, object Marshaler) error {
	file, err := os.Create(filePath + ".gob")
	defer file.Close()
	if err == nil {
		encoder := gob.NewEncoder(file)
		err := encoder.Encode(object)
		if err != nil {
			return err
		}
	}
	return err
}

func readGob(filePath string, object Marshaler) error {
	file, err := os.Open(filePath + ".gob")
	defer file.Close()
	if err == nil {
		decoder := gob.NewDecoder(file)
		err := decoder.Decode(object)
		if err != nil {
			return err
		}
	}
	return err
}
