package main

type ListStorage interface {
	ListStorageReader
	ListStorageWriter
	ListStorageDeleter
}

type ListStorageReader interface {
	Read() EntryList
}

type ListStorageWriter interface {
	Write(e Entry)
}

type ListStorageDeleter interface {
	DeleteEntry(title string)
}
