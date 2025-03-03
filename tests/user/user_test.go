package user

import (
	"concurrent_money_transfer_system/tests"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	tests.Setup()
	setup()
	code := m.Run()
	os.Exit(code)
}

var testData map[string]tests.TestData

func setup() {
	testData = tests.ReadTestData("test_data.json")
}

func TestCreateUser(t *testing.T) {
	tests.MakeRequestAndValidateResponse(t, testData["TestCreateUser"])
}

func TestGetUser(t *testing.T) {
	// Create a user first
	test_data := testData["TestCreateUser"]
	test_data.Request.Body["id"] = "abc2"
	test_data.Response.Body["id"] = "abc2"
	tests.MakeRequestAndValidateResponse(t, test_data)
	// Get the user
	tests.MakeRequestAndValidateResponse(t, testData["TestGetUser"])

}

func TestInvalidEmailFormat(t *testing.T) {
	tests.MakeRequestAndValidateResponse(t, testData["TestInvalidEmailFormat"])
}

func TestInvalidPhoneFormat(t *testing.T) {
	tests.MakeRequestAndValidateResponse(t, testData["TestInvalidPhoneFormat"])
}

func TestPasswordTooShort(t *testing.T) {
	tests.MakeRequestAndValidateResponse(t, testData["TestPasswordTooShort"])
}

func TestEmailIsRequired(t *testing.T) {
	tests.MakeRequestAndValidateResponse(t, testData["TestEmailIsRequired"])
}
