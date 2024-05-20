package credentials

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	t.Run("When all credentials are present", func(t *testing.T) {
		os.Setenv("VAW_TWILIO_ACCOUNT_SID", "test")
		os.Setenv("VAW_TWILIO_AUTH_TOKEN", "test")
		os.Setenv("VAW_TWILIO_NUMBER", "test")
		os.Setenv("VAW_META_ACCOUNT_ID", "test")
		os.Setenv("VAW_META_TOKEN", "test")
		err := Init()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		os.Unsetenv("VAW_TWILIO_ACCOUNT_SID")
		os.Unsetenv("VAW_TWILIO_AUTH_TOKEN")
		os.Unsetenv("VAW_TWILIO_NUMBER")
		os.Unsetenv("VAW_META_ACCOUNT_ID")
		os.Unsetenv("VAW_META_TOKEN")
	})
	t.Run("When missing credentials", func(t *testing.T) {
		err := Init()
		expected_error := "Missing credential: TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_NUMBER, META_ACCOUNT_ID, META_TOKEN"
		if err.Error() != expected_error {
			t.Errorf("Expected error %s, got %s", expected_error, err.Error())
		}
	})
}

func TestMustInit(t *testing.T) {
	t.Run("When all credentials are present", func(t *testing.T) {
		os.Setenv("VAW_TWILIO_ACCOUNT_SID", "test")
		os.Setenv("VAW_TWILIO_AUTH_TOKEN", "test")
		os.Setenv("VAW_TWILIO_NUMBER", "test")
		os.Setenv("VAW_META_ACCOUNT_ID", "test")
		os.Setenv("VAW_META_TOKEN", "test")
		err := Init()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		os.Unsetenv("VAW_TWILIO_ACCOUNT_SID")
		os.Unsetenv("VAW_TWILIO_AUTH_TOKEN")
		os.Unsetenv("VAW_TWILIO_NUMBER")
		os.Unsetenv("VAW_META_ACCOUNT_ID")
		os.Unsetenv("VAW_META_TOKEN")
	})
	t.Run("When missing credentials", func(t *testing.T) {
		defer func() {
			result := recover()
			if result == nil {
				t.Errorf("The code did not panic")
			}
			panic_message := result.(error).Error()
			expected_error := "Missing credential: TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_NUMBER, META_ACCOUNT_ID, META_TOKEN"
			if panic_message != expected_error {
				t.Errorf("Expected panic with message %s, got panic with %s", expected_error, panic_message)
			}
		}()

		MustInit()
	})
}
