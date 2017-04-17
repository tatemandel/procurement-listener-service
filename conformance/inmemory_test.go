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
	"log"
	"os"
	"procurementlistenerservice/inmemory"
	"procurementlistenerservice/server"
	"testing"
)

const (
	TEST_PORT int = 54345
)

var service *inmemory.InMemoryService

type inMemoryTestContext struct {
	t *testing.T
}

var _ TestContext = inMemoryTestContext{}

func (c inMemoryTestContext) Port() int {
	return TEST_PORT
}

func (c inMemoryTestContext) T() *testing.T {
	return c.t
}

func (c inMemoryTestContext) GetEntitlements() []inmemory.EntitlementInfo {
	v := make([]inmemory.EntitlementInfo, 0)

	for _, value := range service.Entitlements {
		v = append(v, value)
	}
	return v
}

func TestInMemoryService(t *testing.T) {
	context := inMemoryTestContext{
		t: t,
	}

	for _, test := range Tests {
		service.Reset()
		t.Run(test.Name, func(t *testing.T) {
			test.Execute(context)
		})
	}
}

func TestMain(m *testing.M) {
	service = inmemory.CreateService(metadata)
	s, err := server.CreateServer(TEST_PORT, service)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	go s.Start()

	m.Run()
}
