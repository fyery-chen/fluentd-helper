package helper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/rancher/fluentd-helper/config"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// ReloadPath full "http://127.0.0.1:24444/api/config.reload"
const ReloadPath = "/api/config.reload"

func ReloadFluentd() error {
	ReloadURL := fmt.Sprintf("http://%s%s", config.FluentdAddress, ReloadPath)
	req, err := http.NewRequest("GET", ReloadURL, nil)
	if err != nil {
		return err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		logrus.Infof("get error when call  %s, details: %s", ReloadURL, err.Error())
		return err
	}
	if res.StatusCode != 200 {
		responseData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		logrus.Infof("response from %s, status code: %v, res %s", ReloadURL, res.Status, string(responseData))
	}

	logrus.Infof("reponse 200 from reload fluentd")
	defer res.Body.Close()

	return nil
}
