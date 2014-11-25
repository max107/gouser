package gouser

import (
	"os/user"
	"testing"
)

func TestGenerateRandom(t *testing.T) {
	_, err := GeneratePassword()
	if err != nil {
		t.Error("Failed to generate random string")
	}
}

func TestHasUser(t *testing.T) {
	user, _ := user.Current()
	has := HasUser(user.Username)
	if has == false {
		t.Error("Failed to get user")
	}
}

func TestCreateUser(t *testing.T) {
	username := "container_test"
	password, _ := GeneratePassword()
	user, err := CreateUser(username, password)
	if err != nil {
		t.Error("Failed to create user", err)
	}

	if username != user.Username {
		t.Errorf("%q != %q", username, user.Username)
	}
}
