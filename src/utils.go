package src

import (
	"fmt"
	"github.com/flytam/filenamify"
	"github.com/u2takey/go-utils/uuid"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveToDisk(file []byte, filename ...string) bool {
	var fullOutPath string
	if len(filename) == 0 {
		fullOutPath = filepath.Join(AUTOHNE_APP_CONTEXT, "out", random())
	} else {
		fullOutPath = parsePath(filename[0])
	}

	err := os.MkdirAll(filepath.Dir(fullOutPath), os.ModePerm)
	err = os.WriteFile(fullOutPath, file, 0644)
	if err != nil {
		log.Printf(err.Error())
		return false
	}
	return true
}

func ReadFromDisk(path string) []byte {
	fullPath := parsePath(path)
	file, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func Pathify(filepath string) string {
	result, err := filenamify.Filenamify(filepath, filenamify.Options{})
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	result = strings.ReplaceAll(result, " ", "_")
	result = strings.ReplaceAll(result, "@", "+")
	return result
}

func parsePath(path string) string {
	fullPath := path
	idx := strings.Index(fullPath, "@random")
	for idx != -1 {
		fullPath = fullPath[:idx] + random() + fullPath[idx+7:]
		idx = strings.Index(fullPath, "@random")
	}
	fullPath = strings.ReplaceAll(fullPath, "@id", id())
	fullPath = strings.Replace(fullPath, "@videos", AUTOHNE_VIDEOS_DIR, -1)
	fullPath = strings.Replace(fullPath, "@app", AUTOHNE_APP_CONTEXT, -1)

	dir, _ := filepath.Split(fullPath)
	if dir == "" {
		fullPath = filepath.Join(AUTOHNE_APP_CONTEXT, "out", fullPath)
	}
	fullPath = filepath.Clean(fullPath)

	// @first ... name of first file in directory
	fstIdx := strings.Index(fullPath, "@first")
	if fstIdx != -1 {
		ext := filepath.Ext(fullPath)
		if fstIdx != len(fullPath)-6-len(ext) {
			log.Fatal("The '@first' wildcard is only allowed at the end of the path!")
		}
		entries, err := os.ReadDir(fullPath[:fstIdx])
		if err != nil {
			log.Fatal(err)
		}
		firstname := ""
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if filepath.Ext(entry.Name()) != ext {
				continue
			}
			firstname = entry.Name()
			break
		}
		if firstname == "" {
			log.Fatal("The '@first' wildcard couldn't be resolved. No '" + ext + "' files in directory!")
		}
		fullPath = fullPath[:fstIdx] + firstname
	}

	return fullPath
}

func random() string {
	date := time.Now().Format("2006-01-02")
	return date + "_" + uuid.NewUUID()
}

func id() string {
	date := time.Now().Format("20060102")
	return "id-" + date + "-" + uuid.NewUUID()
}
