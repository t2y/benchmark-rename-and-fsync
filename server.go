package main

import (
	"fmt"
	"net/http"
)

type fsyncFadviceHandler struct {
	size   int
	pathCh chan string
}

func (h *fsyncFadviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, ok := <-h.pathCh
	if !ok {
		return
	}
	createFsyncFadviceFile(path, h.size*KiB, argSyncFileRange, flagSyncFileRange)
	w.Write([]byte(fmt.Sprintf("%s was created\n", path)))
}

type nosyncHandler struct {
	size   int
	pathCh chan string
}

func (h *nosyncHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, ok := <-h.pathCh
	if !ok {
		return
	}
	createNoSyncFile(path, h.size*KiB)
	w.Write([]byte(fmt.Sprintf("%s was created\n", path)))
}
