package gouser

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

type User struct {
	Username string
	Uid      string
	Gid      string
	Home     string
}

func generateRandom(size int) (string, error) {
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return "", err
	}
	return hex.EncodeToString(id)[:size], nil
}

func ListUser() (map[string]User, error) {
	return parseFile("/etc/passwd")
}

func parseFile(path string) (map[string]User, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return parseReader(file)
}

func parseReader(r io.Reader) (map[string]User, error) {
	lines := bufio.NewReader(r)
	entries := make(map[string]User)
	for {
		line, _, err := lines.ReadLine()
		if err != nil {
			break
		}
		name, entry, err := parseLine(string(copyBytes(line)))
		if err != nil {
			return nil, err
		}
		entries[name] = entry
	}
	return entries, nil
}

func parseLine(line string) (string, User, error) {
	fs := strings.Split(line, ":")
	if len(fs) != 7 {
		return "", User{}, errors.New("Unexpected number of fields in /etc/passwd")
	}
	return fs[0], User{fs[0], fs[2], fs[3], fs[5]}, nil
}

func copyBytes(x []byte) []byte {
	y := make([]byte, len(x))
	copy(y, x)
	return y
}

func GeneratePassword() (string, error) {
	return generateRandom(10)
}

func CreateUser(username string, password string) error {
	args := []string{username}
	err := exec.Command("adduser", args...).Start()
	if err != nil {
		return err
	}

	args = []string{password, "|", "passwd", username, "--stdin"}
	err = exec.Command("echo", args...).Start()
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(username string) {
	args := []string{"-f", username}
	err := exec.Command("deluser", args...).Start()
	if err != nil {
		panic(err)
	}
}

func GetUser(username string) (User, error) {
	users, err := ListUser()
	if err != nil {
		return User{}, err
	}
	user, _ := users[username]
	return user, nil
}

func HasUser(username string) bool {
	_, err := GetUser(username)
	return err == nil
}
