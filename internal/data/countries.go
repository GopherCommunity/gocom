package data

import (
	"fmt"
	"io/ioutil"
	"sort"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

type Countries struct {
	data  yaml.MapSlice
	mutex sync.RWMutex
}

type ByKey yaml.MapSlice

func (slice ByKey) Len() int {
	return len(slice)
}

func (slice ByKey) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice ByKey) Less(i, j int) bool {
	fmt.Printf("%d < %d?\n", i, j)
	a, aok := slice[i].Key.(string)
	b, bok := slice[j].Key.(string)
	if !aok || !bok {
		panic("MapSlice has to have string-keys")
	}
	return a < b
}

func (cl *Countries) Add(country Country) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	if cl.unsafeContainsCode(country.Code) {
		return
	}
	cl.data = append(cl.data, yaml.MapItem{
		Key:   country.Code,
		Value: country,
	})

	sort.Sort(ByKey(cl.data))
}

func (cl *Countries) ContainsCode(code string) bool {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()
	return cl.unsafeContainsCode(code)
}

func (cl *Countries) unsafeContainsCode(code string) bool {
	for _, item := range cl.data {
		if item.Key == code {
			return true
		}
	}
	return false
}

func (cl *Countries) WriteToFile(fpath string) error {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()
	raw, err := yaml.Marshal(cl.data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fpath, raw, 0644)
}

func LoadCountries(fpath string) (*Countries, error) {
	result := &Countries{
		data: make([]yaml.MapItem, 0, 10),
	}
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &result.data); err != nil {
		return nil, err
	}
	return result, nil
}

type Country struct {
	Code      string `yaml:"code"`
	Name      string `yaml:"name"`
	Continent string `yaml:"continent"`
}
