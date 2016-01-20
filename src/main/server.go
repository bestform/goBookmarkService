package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type SaveResponse struct {
	Status, Msg string
}

type Entry struct {
	Title, Url string
}

type EntryList struct {
	Entries map[string](Entry)
}

var storage ListStorage

func main() {
	//storage = FileStorage{FileName: "data.json"}
	s, err := NewSqliteStorage("data.db")
	if err != nil {
		log.Fatal(err)
	}
	storage = s
	http.HandleFunc("/save", save)
	http.HandleFunc("/list", list)
	http.HandleFunc("/delete", deleteEntry)
	http.ListenAndServe(":8080", nil)
}

func deleteEntry(w http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	if "" == title {
		w.Write(getFormattedSaveResponse("1", "No title given"))
		return
	}
	storage.DeleteEntry(title)
	w.Write(getFormattedSaveResponse("0", "Entry deleted"))
}

func list(w http.ResponseWriter, req *http.Request) {
	resp := storage.Read()
	jsonResp, _ := json.MarshalIndent(resp, "", "    ")
	w.Write([]byte(jsonResp))
}

func save(w http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	if "" == title {
		w.Write(getFormattedSaveResponse("1", "No title given"))
		return
	}
	url := req.FormValue("url")
	if "" == url {
		w.Write(getFormattedSaveResponse("1", "No url given"))
		return
	}
	storage.Write(Entry{Title: title, Url: url})
	w.Write(getFormattedSaveResponse("0", "Link Saved"))
}

func getFormattedSaveResponse(status, msg string) []byte {
	resp, _ := json.MarshalIndent(SaveResponse{Status: status, Msg: msg}, "", "    ")
	return []byte(resp)
}
