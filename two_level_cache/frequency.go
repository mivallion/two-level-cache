package two_level_cache

import (
	"sort"
	"sync"
)

type frequenciesMap struct {
	mux         sync.RWMutex
	frequencies map[string]int
}

func NewFrequenciesMap() *frequenciesMap {
	return &frequenciesMap{frequencies: make(map[string]int)}
}

func (fm *frequenciesMap) GetFrequencies() []int {
	var res []int
	fm.mux.RLock()
	defer fm.mux.RUnlock()
	for _, k := range fm.frequencies {
		res = append(res, k)
	}
	return res
}

func (fm *frequenciesMap) GetFrequency(key string) (int, bool) {
	fm.mux.RLock()
	defer fm.mux.RUnlock()
	freq, ok := fm.frequencies[key]
	return freq, ok
}

func (fm *frequenciesMap) SetFrequency(key string, freq int) {
	fm.mux.Lock()
	defer fm.mux.Unlock()
	fm.frequencies[key] = freq
}

func (fm *frequenciesMap) RemoveFrequency(key string) {
	fm.mux.Lock()
	defer fm.mux.Unlock()
	delete(fm.frequencies, key)
}

func (fm *frequenciesMap) Clear() {
	fm.mux.Lock()
	defer fm.mux.Unlock()
	fm.frequencies = make(map[string]int)
}

func (fm *frequenciesMap) GetMostFrequentlyUsedKeys() []string {
	fm.mux.RLock()
	defer fm.mux.RUnlock()
	var keys []string
	var values []int
	for k, v := range fm.frequencies {
		keys = append(keys, k)
		values = append(values, v)
	}
	reverseMap := map[int]string{}
	var reverseMapValues []string
	var reverseMapKeys []int
	for i := range keys {
		k := keys[i]
		v := values[i]
		reverseMap[v] = k
		reverseMapKeys = append(reverseMapKeys, v)
	}
	sort.Ints(reverseMapKeys)
	for _, v := range reverseMapKeys {
		reverseMapValues = append(reverseMapValues, reverseMap[v])
	}
	return reverseMapValues
}
