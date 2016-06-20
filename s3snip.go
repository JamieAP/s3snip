package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"

	"github.com/atotto/clipboard"
	"github.com/rlmcpherson/s3gof3r"
)

const s3UrlTemplate = "https://s3-%s.amazonaws.com/%s/%s"
const bitlyBase = "https://api-ssl.bitly.com/v3/shorten?access_token=%s&longUrl=%s.png&format=txt&domain=j.mp"

type config struct {
	AwsRegion        string `json:"awsRegion"`
	AwsAccessKey     string `json:"awsAccessKey"`
	AwsSecretKey     string `json:"awsSecretKey"`
	AwsBucket        string `json:"awsBucket"`
	BitlyAccessToken string `json:"bitlyAccessToken"`
}

func getUserHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func readConfig() config {
	confFile, err := os.Open(getUserHome() + "/.s3snip/conf.json")
	if err != nil {
		log.Fatal(err)
	}
	config := config{}
	decoder := json.NewDecoder(confFile)
	decodeErr := decoder.Decode(&config)
	if decodeErr != nil {
		log.Fatal(err)
	}
	return config
}

func readCloserToString(reader io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	string := buf.String()
	return string
}

func takeScreenshot() []byte {
  err := exec.Command("screencapture", "-s", "/tmp/screenshot.png").Run()
  if err != nil {
    log.Fatal(err)
  }

  file, err := ioutil.ReadFile("/tmp/screenshot.png")
  if err != nil {
    log.Fatal(err)
  }
  return file
}

func main() {
	conf := readConfig()

  screenshot := takeScreenshot()
	hashBytes := sha1.Sum(screenshot)
	hashString := hex.EncodeToString(hashBytes[:])

  awsKeys := s3gof3r.Keys{
		AccessKey:     conf.AwsAccessKey,
		SecretKey:     conf.AwsSecretKey,
		SecurityToken: "",
	}

	s3 := s3gof3r.New(fmt.Sprintf("s3-%s.amazonaws.com", conf.AwsRegion), awsKeys)
	bucket := s3.Bucket(conf.AwsBucket)
	header := make(http.Header)
	header.Add("Content-Type", "image/png")

	writer, err := bucket.PutWriter(hashString+".png", header, nil)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = io.Copy(writer, bytes.NewBuffer(screenshot)); err != nil {
		log.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	s3Url := fmt.Sprintf(s3UrlTemplate, conf.AwsRegion, conf.AwsBucket, hashString)
	bitlyUrl := fmt.Sprintf(bitlyBase, conf.BitlyAccessToken, s3Url)

	resp, err := http.Get(bitlyUrl)
	if err != nil {
		log.Fatal(err)
	}

	clipboard.WriteAll(readCloserToString(resp.Body))
}
