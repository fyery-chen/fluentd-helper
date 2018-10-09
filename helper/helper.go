package helper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"github.com/Sirupsen/logrus"

	"bytes"
	"github.com/rancher/fluentd-helper/config"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// ReloadPath full "http://127.0.0.1:24444/api/config.reload"
const ReloadPath = "/api/config.reload"

func ReloadFluentd() error {
	logrus.Infof("sending request to reload fluentd")
	ReloadURL := fmt.Sprintf("http://%s%s", config.FluentdAddress, ReloadPath)
	req, err := http.NewRequest("GET", ReloadURL, nil)
	if err != nil {
		return err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		logrus.Errorf("get error when call  %s, details: %s", ReloadURL, err.Error())
		return err
	}
	if res.StatusCode != 200 {
		responseData, err := ioutil.ReadAll(res.Body)
		logrus.Infof("reload fluentd fail, response from %s, status code: %v, res %s", ReloadURL, res.Status, string(responseData))
		if err != nil {
			logrus.Errorf("read fluentd response fail: %s", err)
			return err
		}
	}

	logrus.Infof("reponse 200 from reload fluentd")
	defer res.Body.Close()

	return nil
}

func RenewTicket() error {
	logrus.Infof("renew ticket for new users")
	cmd := exec.Command("/kerberos-init.sh")
	bufErr := &bytes.Buffer{}
	cmd.Stderr = bufErr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
