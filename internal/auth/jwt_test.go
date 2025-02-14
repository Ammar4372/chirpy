package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func jwtTest(t *testing.T) {
	id := uuid.New()
	secret := "ewjfheihfiehfoiehfo2i3"
	token, err := MakeJWT(id, secret, 5*time.Second)
	if err != nil {
		t.Errorf("Error Creating JWT")
	}
	_, err = ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("Error Validating JWT")
	}
	_, err = ValidateJWT(token, "whdwehfihweifhiwefh")
	if err == nil {
		t.Errorf("did not reject token signed with wrong signature")
	}
	time.Sleep(10 * time.Second)
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("did not reject expired token")
	}
}
