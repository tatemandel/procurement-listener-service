// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conformance

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"procurementlistenerservice/inmemory"
	"reflect"
	"testing"
)

// ConformanceTest is the top-level data structure that represents a conformance test. It contains the service/plan
// metadata, as well as individual actions that make up the test.
type ConformanceTest struct {
	Name     string
	Metadata inmemory.Metadata
	Actions  []Action
}

// Action defines the interface of a test action that will be performed as part of a conformance test.
type Action interface {
	execute(TestContext) error
}

// TestContext abstracts the externally supplied data to the conformance test suite. Each particular test implementation
// will need to implement and supply a context.
type TestContext interface {
	Port() int
	T() *testing.T
	GetEntitlements() []inmemory.EntitlementInfo
}

// PostEntitlementEvent is a test action that will post an entitlement event to the service and will expect a particular
// response.
type PostEntitlementEvent struct {
	Request      string
	ExpectedCode int
}

var _ Action = PostEntitlementEvent{}

func (a PostEntitlementEvent) execute(c TestContext) error {
	response, err := post(c.Port(), "entitlementEvents", a.Request)
	if err != nil {
		return err
	}

	if response.StatusCode != a.ExpectedCode {
		return fmt.Errorf(
			"Unexpected HTTP response code: actual='%d', expected='%d'", response.StatusCode, a.ExpectedCode)
	}

	return nil
}

type ExpectEntitlements struct {
	Entitlements []inmemory.EntitlementInfo
}

var _ Action = ExpectEntitlements{}

func (a ExpectEntitlements) execute(c TestContext) error {
	actual := c.GetEntitlements()
	if len(actual) != len(a.Entitlements) {
		return fmt.Errorf("Entitlement count mismatch: actual='%d', expected='%d'.",
			len(actual), len(a.Entitlements))
	}

	for _, e := range a.Entitlements {
		match := false
		for _, a := range actual {
			if reflect.DeepEqual(e, a) {
				match = true
				break
			}
		}
		if !match {
			return fmt.Errorf("No match for expected entitlement: '%+v'", e)
		}
	}

	return nil
}

func post(port int, path string, payload string) (*http.Response, error) {
	url := createUrl(path, port)

	var reader io.Reader
	if payload != "" {
		buffer := bytes.NewBufferString(payload)
		reader = bytes.NewReader(buffer.Bytes())
	}
	return http.Post(url, "application/json", reader)
}

func createUrl(path string, port int) string {
	return fmt.Sprintf("http://localhost:%d/%s", port, path)
}

func (c ConformanceTest) Execute(context TestContext) {
	for _, action := range c.Actions {
		err := action.execute(context)
		if err != nil {
			context.T().Error(err)
			context.T().Fail()
			return
		}
	}
}
