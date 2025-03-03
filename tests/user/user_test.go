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
	tests.MakeRequestAndValidateResponse(t, testData["TestCreateUser"])
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
