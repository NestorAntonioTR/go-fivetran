package fingerprints_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	testutils "github.com/fivetran/go-fivetran/test_utils"
	"github.com/fivetran/go-fivetran/tests"
	"github.com/fivetran/go-fivetran/tests/mock"
)

func TestNewConnectorFingerprintRevokeMock(t *testing.T) {
	// arrange
	testConnectorId := "connector_id"
	testHash := "hash"

	ftClient, mockClient := tests.CreateTestClient()
	handler := mockClient.When(http.MethodDelete, fmt.Sprintf("/v1/connectors/%v/fingerprints/%v", testConnectorId, testHash)).ThenCall(
		func(req *http.Request) (*http.Response, error) {
			response := mock.NewResponse(req, http.StatusOK,
				`{
					"code": "Success", 
					"message": "The certificate has been revoked."
				}`)

			return response, nil
		})

	// act & assert
	response, err := ftClient.NewConnectorFingerprintRevoke().
		ConnectorID(testConnectorId).
		Hash(testHash).
		Do(context.Background())

	if err != nil {
		t.Logf("%+v\n", response)
		t.Error(err)
	}

	interactions := mockClient.Interactions()
	testutils.AssertEqual(t, len(interactions), 1)
	testutils.AssertEqual(t, interactions[0].Handler, handler)
	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Code, "Success")
	testutils.AssertNotEmpty(t, response.Message)
}