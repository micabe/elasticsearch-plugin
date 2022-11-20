package server

import (
	"errors"
	"net/http"
	"testing"

	"github.com/halverneus/static-file-server/config"
	"github.com/halverneus/static-file-server/handle"
)

func TestRun(t *testing.T) {
	listenerError := errors.New("listener")
	selectListener = func() handle.ListenerFunc {
		return func(string, http.HandlerFunc) error {
			return listenerError
		}
	}

	config.Get.Debug = false
	if err := Run(); listenerError != err {
		t.Errorf("Without debug expected %v but got %v", listenerError, err)
	}

	config.Get.Debug = true
	if err := Run(); listenerError != err {
		t.Errorf("With debug expected %v but got %v", listenerError, err)
	}
}

func TestHandlerSelector(t *testing.T) {
	// This test only exercises function branches.
	testFolder := "/web"
	testPrefix := "/url/prefix"
	var ignoreReferrer []string
	testReferrer := []string{"http://localhost"}
	testAccessKey := "access-key"

	testCases := []struct {
		name      string
		folder    string
		prefix    string
		listing   bool
		debug     bool
		refer     []string
		cors      bool
		accessKey string
	}{
		{"Basic handler w/o debug", testFolder, "", true, false, ignoreReferrer, false, ""},
		{"Prefix handler w/o debug", testFolder, testPrefix, true, false, ignoreReferrer, false, ""},
		{"Basic and hide listing handler w/o debug", testFolder, "", false, false, ignoreReferrer, false, ""},
		{"Prefix and hide listing handler w/o debug", testFolder, testPrefix, false, false, ignoreReferrer, false, ""},
		{"Basic handler w/debug", testFolder, "", true, true, ignoreReferrer, false, ""},
		{"Prefix handler w/debug", testFolder, testPrefix, true, true, ignoreReferrer, false, ""},
		{"Basic and hide listing handler w/debug", testFolder, "", false, true, ignoreReferrer, false, ""},
		{"Prefix and hide listing handler w/debug", testFolder, testPrefix, false, true, ignoreReferrer, false, ""},
		{"Basic handler w/o debug w/refer", testFolder, "", true, false, testReferrer, false, ""},
		{"Prefix handler w/o debug w/refer", testFolder, testPrefix, true, false, testReferrer, false, ""},
		{"Basic and hide listing handler w/o debug w/refer", testFolder, "", false, false, testReferrer, false, ""},
		{"Prefix and hide listing handler w/o debug w/refer", testFolder, testPrefix, false, false, testReferrer, false, ""},
		{"Basic handler w/debug w/refer w/o cors", testFolder, "", true, true, testReferrer, false, ""},
		{"Prefix handler w/debug w/refer w/o cors", testFolder, testPrefix, true, true, testReferrer, false, ""},
		{"Basic and hide listing handler w/debug w/refer w/o cors", testFolder, "", false, true, testReferrer, false, ""},
		{"Prefix and hide listing handler w/debug w/refer w/o cors", testFolder, testPrefix, false, true, testReferrer, false, ""},
		{"Basic handler w/debug w/refer w/cors", testFolder, "", true, true, testReferrer, true, ""},
		{"Prefix handler w/debug w/refer w/cors", testFolder, testPrefix, true, true, testReferrer, true, ""},
		{"Basic and hide listing handler w/debug w/refer w/cors", testFolder, "", false, true, testReferrer, true, ""},
		{"Prefix and hide listing handler w/debug w/refer w/cors", testFolder, testPrefix, false, true, testReferrer, true, ""},
		{"Access Key and Basic handler w/o debug", testFolder, "", true, false, ignoreReferrer, false, testAccessKey},
		{"Access Key and Prefix handler w/o debug", testFolder, testPrefix, true, false, ignoreReferrer, false, testAccessKey},
		{"Access Key and Basic and hide listing handler w/o debug", testFolder, "", false, false, ignoreReferrer, false, testAccessKey},
		{"Access Key and Prefix and hide listing handler w/o debug", testFolder, testPrefix, false, false, ignoreReferrer, false, testAccessKey},
		{"Access Key and Basic handler w/debug", testFolder, "", true, true, ignoreReferrer, false, testAccessKey},
		{"Access Key and Prefix handler w/debug", testFolder, testPrefix, true, true, ignoreReferrer, false, testAccessKey},
		{"Access Key and Basic and hide listing handler w/debug", testFolder, "", false, true, ignoreReferrer, false, testAccessKey},
		{"Access Key and Prefix and hide listing handler w/debug", testFolder, testPrefix, false, true, ignoreReferrer, false, testAccessKey},
		{"Access Key and Basic handler w/o debug w/refer", testFolder, "", true, false, testReferrer, false, testAccessKey},
		{"Access Key and Prefix handler w/o debug w/refer", testFolder, testPrefix, true, false, testReferrer, false, testAccessKey},
		{"Access Key and Basic and hide listing handler w/o debug w/refer", testFolder, "", false, false, testReferrer, false, testAccessKey},
		{"Access Key and Prefix and hide listing handler w/o debug w/refer", testFolder, testPrefix, false, false, testReferrer, false, testAccessKey},
		{"Access Key and Basic handler w/debug w/refer w/o cors", testFolder, "", true, true, testReferrer, false, testAccessKey},
		{"Access Key and Prefix handler w/debug w/refer w/o cors", testFolder, testPrefix, true, true, testReferrer, false, testAccessKey},
		{"Access Key and Basic and hide listing handler w/debug w/refer w/o cors", testFolder, "", false, true, testReferrer, false, testAccessKey},
		{"Access Key and Prefix and hide listing handler w/debug w/refer w/o cors", testFolder, testPrefix, false, true, testReferrer, false, testAccessKey},
		{"Access Key and Basic handler w/debug w/refer w/cors", testFolder, "", true, true, testReferrer, true, testAccessKey},
		{"Access Key and Prefix handler w/debug w/refer w/cors", testFolder, testPrefix, true, true, testReferrer, true, testAccessKey},
		{"Access Key and Basic and hide listing handler w/debug w/refer w/cors", testFolder, "", false, true, testReferrer, true, testAccessKey},
		{"Access Key and Prefix and hide listing handler w/debug w/refer w/cors", testFolder, testPrefix, false, true, testReferrer, true, testAccessKey},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Get.Debug = tc.debug
			config.Get.Folder = tc.folder
			config.Get.ShowListing = tc.listing
			config.Get.URLPrefix = tc.prefix
			config.Get.Referrers = tc.refer
			config.Get.Cors = tc.cors
			config.Get.AccessKey = tc.accessKey

			handlerSelector()
		})
	}
}

func TestListenerSelector(t *testing.T) {
	// This test only exercises function branches.
	testCert := "file.crt"
	testKey := "file.key"

	testCases := []struct {
		name string
		cert string
		key  string
	}{
		{"HTTP", "", ""},
		{"HTTPS", testCert, testKey},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Get.TLSCert = tc.cert
			config.Get.TLSKey = tc.key
			listenerSelector()
		})
	}
}
