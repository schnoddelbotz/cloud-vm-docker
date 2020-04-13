package cloud

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// GetUserSSHPublicKeys ...
func GetUserSSHPublicKeys(fromFile string, letItBe bool) string {
	pubKeys := ""
	if letItBe {
		log.Printf("Skipping SSH pubkey collection/installation as requested")
		return pubKeys
	}
	if fromFile != "" {
		pubKeys = getSSHKeyFromFile(fromFile)
	} else {
		pubKeys = getAllSSHPublicKeys()
	}
	if len(pubKeys) > 1 && pubKeys[len(pubKeys)-1:] == "\n" {
		pubKeys = pubKeys[:len(pubKeys)-2]
	}
	return pubKeys
}

// Base64Encode ...
func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// Gzip could be applied on cloud_init data, then Base64 encoded, see README.md->readthedocs/cloud_init
func Gzip(input string) []byte {
	return []byte("ABRA/CAaDabrA")
}

func getSSHKeyFromFile(fileName string) string {
	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Failed reading SSH public key %s: %s", fileName, err)
	}
	return "cloud-vm-docker:" + string(fileContents)
}

func getAllSSHPublicKeys() string {
	// FIXME BUG. expected format is: user:<PUBKEY>\n user:<PUBKEY>
	// so ... must pick a single pubkey, as using same username multiple times will not work? verify!
	keys := ""
	baseDirectory, _ := homedir.Dir()
	scanDirGlob := fmt.Sprintf("%s/%s/*.pub", baseDirectory, ".ssh")
	log.Printf("Scanning %s for SSH public keys", scanDirGlob)
	files, err := filepath.Glob(scanDirGlob)
	if err != nil {
		log.Fatalf("Failed to look for SSH public keys: %s -- maybe try --no-ssh", err)
	}
	for _, fileName := range files {
		log.Printf("Adding SSH public key: %s", fileName)
		fileContents, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatalf("Failed reading SSH public key %s: %s", fileName, err)
		}
		keys += "cloud-vm-docker:" + string(fileContents) + "\n"
	}
	return keys
}
