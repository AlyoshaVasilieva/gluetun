// Package storage defines interfaces to interact with the files persisted such as the list of servers.
package storage

import (
	"path/filepath"

	"github.com/qdm12/gluetun/internal/models"
)

//go:generate mockgen -destination=infoerrorer_mock_test.go -package $GOPACKAGE . InfoErrorer

type Storage struct {
	mergedServers    models.AllServers
	hardcodedServers models.AllServers
	logger           InfoErrorer
	filepath         string
}

type InfoErrorer interface {
	Info(s string)
}

// New creates a new storage and reads the servers from the
// embedded servers file and the file on disk in directory+/servers.json.
// Passing an empty filepath disables writing servers to a file.
func New(logger InfoErrorer, directory string) (storage *Storage, err error) {
	// error returned covered by unit test
	harcodedServers, _ := parseHardcodedServers()

	const serversFilename = "servers.json"
	filePath := filepath.Join(directory, serversFilename)

	storage = &Storage{
		hardcodedServers: harcodedServers,
		logger:           logger,
		filepath:         filePath,
	}

	if err := storage.SyncServers(); err != nil {
		return nil, err
	}

	return storage, nil
}
