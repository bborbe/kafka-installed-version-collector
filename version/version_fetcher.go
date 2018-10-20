// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"context"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/bborbe/kafka-installed-version-collector/avro"
	"github.com/pkg/errors"
)

//go:generate counterfeiter -o ../mocks/http_client.go --fake-name HttpClient . httpClient
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Fetcher struct {
	HttpClient httpClient
	Url        string
	App        string
	Regex      string
}

func (f *Fetcher) Fetch(ctx context.Context) (*avro.ApplicationVersionInstalled, error) {
	req, err := http.NewRequest(http.MethodGet, f.Url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create request failed")
	}
	resp, err := f.HttpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}
	if resp.StatusCode/100 != 2 {
		return nil, errors.New("request status code != 2xx")
	}
	defer resp.Body.Close()
	re, err := regexp.Compile(f.Regex)
	if err != nil {
		return nil, errors.Wrap(err, "compile regex failed")
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body failed")
	}
	result := re.FindSubmatch(content)
	if len(result) != 2 {
		return nil, errors.New("find regex failed")
	}
	version := string(result[1])
	glog.V(3).Infof("found version: %s", version)
	return &avro.ApplicationVersionInstalled{
		App:     f.App,
		Url:     f.Url,
		Version: version,
	}, nil
}
