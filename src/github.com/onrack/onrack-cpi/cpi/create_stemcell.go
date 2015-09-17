package cpi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/nu7hatch/gouuid"
	"github.com/onrack/onrack-cpi/config"
	"github.com/onrack/onrack-cpi/stemcell"
)

func CreateStemcell(config config.Cpi, extInput ExternalInput) (string, error) {
	var imagePath string

	if reflect.TypeOf(extInput[0]) == reflect.TypeOf(imagePath) {
		imagePath = extInput[0].(string)
	} else {
		return "", errors.New("Received unexpected type for stemcell image path")
	}

	s := stemcell.New(imagePath)
	stemcellHandle, err := s.Extract()
	if err != nil {
		return "", err
	}
	defer stemcellHandle.Close()

	uuid, err := uuid.NewV4()
	if err != nil {
		return "", errors.New("Error generating UUID")
	}

	url := fmt.Sprintf("http://%s:8080/api/common/files/%s", config.ApiServer, uuid.String())
	body := ioutil.NopCloser(stemcellHandle)
	request, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Printf("Error building request to api server: %s", err)
		return "", err
	}

	fileInfo, err := stemcellHandle.Stat()
	if err != nil {
		return "", err
	}
	request.ContentLength = fileInfo.Size()

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Error making request to api server: %s", err)
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read response body")
		return "", err
	}

	if resp.StatusCode != 201 {
		log.Printf("Failed uploading stemcell '%s'", resp.Status)
		return "", fmt.Errorf("Error uploading stemcell: %s", string(bodyBytes))
	}

	stemcell_uuid := string(bodyBytes)
	log.Printf("Succeeded uploading stemcell got '%s' as uuid", stemcell_uuid)

	return uuid.String(), nil
}