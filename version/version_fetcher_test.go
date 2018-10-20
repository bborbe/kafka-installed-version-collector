package version_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"github.com/bborbe/kafka-installed-version-collector/mocks"
	"github.com/bborbe/kafka-installed-version-collector/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version Fetcher", func() {
	var httpClient *mocks.HttpClient
	var fetcher *version.Fetcher
	BeforeEach(func() {
		httpClient = &mocks.HttpClient{}
		fetcher = &version.Fetcher{
			HttpClient: httpClient,
			App:        "Banana",
			Url:        "http://www.example.com",
			Regex:      `<meta\s+name="ajs-version-number"\s+content="([^"]+)">`,
		}
	})
	It("return error if empty", func() {
		httpClient.DoReturnsOnCall(0, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
		}, nil)

		_, err := fetcher.Fetch(context.Background())
		Expect(err).To(HaveOccurred())
	})
	It("returns version", func() {
		httpClient.DoReturnsOnCall(0, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`<meta name="ajs-version-number" content="1.3.37">`)),
		}, nil)
		version, err := fetcher.Fetch(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(version.App).To(Equal("Banana"))
		Expect(version.Url).To(Equal("http://www.example.com"))
		Expect(version.Version).To(Equal("1.3.37"))
	})
})
