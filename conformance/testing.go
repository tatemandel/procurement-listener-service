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

// Represents a conformance test
type ConformanceTest struct {
	Name     string
	metadata inmemory.Metadata
	actions  []testAction
}

type ConformanceTestBuilder struct {
	test *ConformanceTest
}

type ConformanceTestResponseBuilder struct {
	payload string
	test    *ConformanceTest
}

type TestContext interface {
	Port() int
	T() *testing.T
	GetEntitlements() []inmemory.EntitlementInfo
}

type testAction interface {
	execute(TestContext) error
}

type expectResponseCodeAction struct {
	path            string
	payload         string
	expectedCode    int
	expectedPayload string
}

var _ testAction = expectResponseCodeAction{}

func (e expectResponseCodeAction) execute(c TestContext) error {
	response, err := post(c.Port(), e.path, e.payload)
	if err != nil {
		return err
	}

	if response.StatusCode != e.expectedCode {
		return fmt.Errorf(
			"Unexpected HTTP response code: actual='%d', expected='%d'", response.StatusCode, e.expectedCode)
	}

	return nil
}

type expectEntitlementCountAction struct {
	count int
}

var _ testAction = expectEntitlementCountAction{}

func (e expectEntitlementCountAction) execute(c TestContext) error {
	actualCount := len(c.GetEntitlements())
	if actualCount != e.count {
		return fmt.Errorf("Entitlement count mismatch: actual='%d', expected='%d'.", actualCount, e.count)
	}

	return nil
}

type expectEntitlementAction struct {
	entitlement inmemory.EntitlementInfo
}

var _ testAction = expectEntitlementAction{}

func (e expectEntitlementAction) execute(c TestContext) error {

	for _, ent := range c.GetEntitlements() {

		if ent.Id == e.entitlement.Id {
			if !reflect.DeepEqual(ent, e.entitlement) {
				return fmt.Errorf("Entitlement mismatch: actual='%+v', expected='%+v'.", ent, e.entitlement)
			}
			return nil
		}
	}

	return fmt.Errorf("Entitlement not found: id='%s'", e.entitlement.Id)
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

func BeginTestCase(name string, metadata inmemory.Metadata) ConformanceTestBuilder {
	return ConformanceTestBuilder{
		test: &ConformanceTest{
			Name:     name,
			metadata: metadata,
			actions:  make([]testAction, 0),
		},
	}
}

func (b ConformanceTestBuilder) Build() ConformanceTest {
	return *b.test
}

func (c ConformanceTest) Execute(context TestContext) {
	for _, action := range c.actions {
		err := action.execute(context)
		if err != nil {
			context.T().Error(err)
			context.T().Fail()
			return
		}
	}
}

func (b ConformanceTestBuilder) PostEntitlementEvent(payload string) ConformanceTestResponseBuilder {
	return ConformanceTestResponseBuilder{
		payload: payload,
		test:    b.test,
	}
}

func (b ConformanceTestResponseBuilder) ExpectResponseCode(code int) ConformanceTestBuilder {
	action := expectResponseCodeAction{
		path:            "entitlementEvents",
		payload:         b.payload,
		expectedCode:    code,
		expectedPayload: "",
	}
	b.test.actions = append(b.test.actions, action)
	return ConformanceTestBuilder{
		test: b.test,
	}
}

func (b ConformanceTestBuilder) ExpectEntitlementCount(count int) ConformanceTestBuilder {
	action := expectEntitlementCountAction{
		count: count,
	}

	b.test.actions = append(b.test.actions, action)
	return ConformanceTestBuilder{
		test: b.test,
	}
}

func (b ConformanceTestBuilder) ExpectEntitlement(e inmemory.EntitlementInfo) ConformanceTestBuilder {
	action := expectEntitlementAction{
		entitlement: e,
	}

	b.test.actions = append(b.test.actions, action)
	return ConformanceTestBuilder{
		test: b.test,
	}
}
