package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type FileStorage struct {
	FileName string
}

func (f FileStorage) Read() EntryList {
	data, err := ioutil.ReadFile(f.FileName)
	if err != nil {
		emptyMap := make(map[string]Entry)
		return EntryList{emptyMap}
	}
	var entryList EntryList
	json.Unmarshal(data, &entryList)

	return entryList
}

func (f FileStorage) Write(e Entry) {
	existing := f.Read()
	existing.Entries[e.Title] = e
	data, _ := json.Marshal(existing)
	ioutil.WriteFile(f.FileName, data, os.ModePerm)
}

func (f FileStorage) DeleteEntry(title string) {
	existing := f.Read()
	if _, ok := existing.Entries[title]; ok == true {
		delete(existing.Entries, title)
		data, _ := json.Marshal(existing)
		ioutil.WriteFile(f.FileName, data, os.ModePerm)
	}
}
