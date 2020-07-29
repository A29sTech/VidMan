package core

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

// Thumbnail Dir Name ;
const (
	ThumbnailDir = "thumbnails"
)

// Create Thumbnail After Video Saved ;
func genThumbnail(videoPath string) {
	filename := path.Join(path.Dir(videoPath), ThumbnailDir, filepath.Base(videoPath)+".png")
	cmdStr := fmt.Sprintf(` -i "%s" -ss 00:00:01.000 -vframes 1 "%s"`, videoPath, filename)

	cmd := exec.Command(`ffmpeg`)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CmdLine:       cmdStr,
		CreationFlags: 0,
	}
	cmd.Run()

}

// SaveUploadedFile : save the file to given path if not exsists ;
func SaveUploadedFile(file io.Reader, pathAndName string) error {

	// Check For file Exsistence ;
	if _, err := os.Stat(pathAndName); err == nil {
		return errors.New("File already exsists")
	}

	// Read BytesArray from io.Reader ;
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Write The File with BytesArray ;
	if err := ioutil.WriteFile(pathAndName, fileBytes, 0644); err != nil {
		return err
	}

	// Craete Thumbnail ;
	go genThumbnail(pathAndName)

	return nil

}
