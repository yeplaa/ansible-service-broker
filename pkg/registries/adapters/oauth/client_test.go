//
// Copyright (c) 2018 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package oauth

import (
	"net/url"
	"strings"
	"testing"

	logging "github.com/op/go-logging"
)

var headerCases = map[string]string{
	"Bearer realm=\"http://foo/a/b/c\",service=\"bar\"":  "http://foo/a/b/c?service=bar",
	"Bearer service=\"bar\",realm=\"http://foo/a/b/c\"":  "http://foo/a/b/c?service=bar",
	"Bearer realm=\"http://foo/a/b/c/\",service=\"bar\"": "http://foo/a/b/c/?service=bar",
	"Bearer realm=\"https://foo\",service=\"bar\"":       "https://foo?service=bar",
	"Bearer realm=\"http://foo/a/b/c\"":                  "http://foo/a/b/c",
}

var headerErrorCases = map[string]string{
	"Bearer service=\"bar\"": "Could not parse www-authenticate header:",
	"Bearer realm=\"\"":      "",
}

var tokenCases = map[string]string{
	"{\"access_token\": \"abc123\"}":                        "abc123",
	"{\"token\": \"abc123\"}":                               "abc123",
	"{\"access_token\": \"abc123\", \"token\": \"def456\"}": "abc123",
	"{}": "",
}

var tokenErrorCases = map[string]string{
	"{\"token\": {}":          "unexpected end of JSON input",
	"{\"access_token\": {}":   "unexpected end of JSON input",
	"{\"token\": null":        "unexpected end of JSON input",
	"{\"access_token\": null": "unexpected end of JSON input",
}

func TestParseAuthHeader(t *testing.T) {
	for in, out := range headerCases {
		result, err := parseAuthHeader(in)
		if err != nil {
			t.Error(err.Error())
		}
		if result.String() != out {
			t.Errorf("Expected %s, got %s", out, result.String())
		}
	}
}

func TestParseAuthHeaderErrors(t *testing.T) {
	for in, out := range headerErrorCases {
		_, err := parseAuthHeader(in)
		if err == nil {
			t.Errorf("Expected an error parsing %s", in)
		} else if strings.HasPrefix(err.Error(), out) == false {
			t.Errorf("Expected prefix %s, got %s", out, err.Error())
		}
	}
}

func TestParseAuthToken(t *testing.T) {
	for in, out := range tokenCases {
		result, err := parseAuthToken([]byte(in))
		if err != nil {
			t.Error(err.Error())
		}
		if result != out {
			t.Errorf("Expected %s, got %s", out, result)
		}
	}
}

func TestParseAuthTokenErrors(t *testing.T) {
	for in, out := range tokenErrorCases {
		_, err := parseAuthToken([]byte(in))
		if err == nil {
			t.Errorf("Expected an error parsing %s", in)
		} else if strings.HasPrefix(err.Error(), out) == false {
			t.Errorf("Expected prefix %s, got %s", out, err.Error())
		}
	}
}

func TestNewRequest(t *testing.T) {
	log := &logging.Logger{}
	u, _ := url.Parse("http://automationbroker.io")
	c := NewClient("foo", "bar", false, u, log)
	c.token = "letmein"
	req, err := c.NewRequest("/v2/")
	if err != nil {
		t.Error(err.Error())
		return
	}
	accepth := req.Header.Get("Accept")
	if accepth != "application/json" {
		t.Errorf("incorrect or missing accept header: %s", accepth)
		return
	}
	authh := req.Header.Get("Authorization")
	if authh != "Bearer letmein" {
		t.Errorf("incorrect or missing authorization header: %s", authh)
		return
	}
}
