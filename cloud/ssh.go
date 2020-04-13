package cloud

import (
	"encoding/base64"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"path/filepath"
)

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
	if len(pubKeys)>1 && pubKeys[len(pubKeys)-1:len(pubKeys)] == "\n" {
		pubKeys = pubKeys[:len(pubKeys)-2]
	}
	return pubKeys
}

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

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
	keys := ""
	baseDirectory, _ := homedir.Dir()
	scanDirGlob := fmt.Sprintf("%s/%s/*.pub", baseDirectory, ".ssh")
	log.Printf("Scanning %s for SSH public keys", scanDirGlob)
	files, err := filepath.Glob(scanDirGlob)
	if err != nil {
		log.Fatalf("Failed to look for SSH public keys: %s -- maybe try --no-ssh", err)
	}
	for _, fileName := range(files) {
		log.Printf("Adding SSH public key: %s", fileName)
		fileContents, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatalf("Failed reading SSH public key %s: %s", fileName, err)
		}
		keys += "cloud-vm-docker:" + string(fileContents) + "\n"
	}
	return keys
}
