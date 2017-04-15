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

	for _, test := range tests {
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
