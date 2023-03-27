package chatbot

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var userList = []*User{NewUser("temporary_user", "default_password", []string{"request"})}
var usernameTable = map[string]*User{userList[0].Name: userList[0]}
var useridTable = map[string]*User{userList[0].Id: userList[0]}

type User struct {
	Id       string
	Name     string
	password string
	access   []string
}

func NewUser(username, password string, access []string) *User {
	id := generateUUID(username)
	return &User{
		Id:       id,
		Name:     username,
		password: hashPassword(password),
		access:   access,
	}
}

func (u *User) String() string {
	return fmt.Sprintf("User(id='%s')", u.Id)
}

func CreateUser(username, password string, access []string) {
	userList = append(userList, NewUser(username, password, access))
	usernameTable[username] = userList[len(userList)-1]
	useridTable[userList[len(userList)-1].Id] = userList[len(userList)-1]
}

func IndexUserWithID(userID string) *User {
	return useridTable[userID]
}

func IndexUserWithName(userName string) *User {
	return usernameTable[userName]
}

func Authenticate(username, password string) *User {
	user := usernameTable[username]
	if user != nil && comparePassword(user.password, password) {
		return user
	}
	return nil
}

func Identity(userid string) *User {
	return useridTable[userid]
}

func generateUUID(username string) string {
	id := uuid.NewMD5(uuid.NameSpaceDNS, []byte(username))
	return id.String()
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func comparePassword(hash, password string) bool {
	h1, err1 := hex.DecodeString(hash)
	h2, err2 := hex.DecodeString(hashPassword(password))
	if err1 != nil || err2 != nil || len(h1) != len(h2) {
		return false
	}
	time.Sleep(5 * time.Millisecond) // Delay to prevent timing attacks
	return compareDigest(h1, h2)
}

func compareDigest(h1, h2 []byte) bool {
	return len(h1) == len(h2) && compareBytes(h1, h2)
}

func compareBytes(a, b []byte) bool {
	diff := byte(0)
	for i := 0; i < len(a); i++ {
		diff |= a[i] ^ b[i]
	}
	return diff == 0
}
