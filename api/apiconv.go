// Copyright 2017 YTD Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// apiconv: Converts Decoded Video data to MP3, WEBM or MP4.
// NOTE: To reimplement using Go ffmpeg bindings.
package api

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

//Downloads decoded audio stream
func ApiConvertVideo(file, id, format string, bitrate uint, decStream []stream) error {
	cmd := exec.Command("ffmpeg", "-i", "-", "-ab", fmt.Sprintf("%dk", bitrate), file)
	if err := os.MkdirAll(filepath.Dir(file), 666); err != nil {
		return err
	}
	out, err := os.Create(file)
	if err != nil {
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(decStream)
	_, err = exec.LookPath("ffmpeg")
	if err != nil {
		return errors.New("ffmpeg not found on system")
	}
	cmd.Start()
	stdin.Write(buf.Bytes())
	out.Write(buf.Bytes())
	return nil
}

//Downloads decoded video stream.
func ApiDownloadVideo(path, file, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Http.Get\nerror: %s\nURL: %s\n", err, url)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Reading Output: status code: '%v'", resp.StatusCode)
		return errors.New("Non 200 status code received")
	}
	err = os.MkdirAll(filepath.Dir(file), 666)
	if err != nil {
		return err
	}
	out, err := os.Create(file)
	if err != nil {
		return err
	}

	//saving downloaded file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("Download Error: ", err)
		return err
	}
	return nil
}
