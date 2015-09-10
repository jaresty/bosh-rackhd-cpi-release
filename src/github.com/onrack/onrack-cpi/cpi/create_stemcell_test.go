package cpi_test

import (
	"github.com/onrack/onrack-cpi/cpi"
	"github.com/onrack/onrack-cpi/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
	"net/http"

	"fmt"
	"io/ioutil"
	"os"

)

var _ = Describe("CreateStemcell", func() {
	Context("With valid CPI v1 input", func() {
		It("Extracts and uploads a VMDK from a vSphere stemcell", func() {
			apiServerIp := os.Getenv("ON_RACK_API_URI")
			Expect(apiServerIp).ToNot(BeEmpty())

			config := config.Cpi{ApiServer: apiServerIp}
			input := `["../spec_assets/stemcell.tgz"]`

			uuid, err := cpi.CreateStemcell(config, input)
			Expect(err).ToNot(HaveOccurred())

			Expect(uuid).ToNot(BeEmpty())
			url := fmt.Sprintf("http://%s:8080/api/common/files/metadata/%s", config.ApiServer, uuid)
			resp, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())

			respBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())

			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			fileMetadataResp := cpi.FileMetadataResponse{}
			err = json.Unmarshal(respBytes, &fileMetadataResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(fileMetadataResp).To(HaveLen(1))

			fileMetadata := fileMetadataResp[0]
			Expect(fileMetadata.Basename).To(Equal(uuid))
		})
	})
	Context("With invalid CPI v1 input", func() {
		It("Returns an error", func() {
			config := config.Cpi{}
			input := `[{"foo": "bar"}]`

			uuid, err := cpi.CreateStemcell(config, input)
			Expect(err).To(MatchError("Received unexpected type for stemcell image path"))
			Expect(uuid).To(BeEmpty())
		})
	})
})



