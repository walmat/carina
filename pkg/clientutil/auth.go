package clientutil

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	"nebula/pkg/api/model"
	"nebula/pkg/util"
	"os"
	"path"
	"sync"
)

func init() {
	dataDir, err := util.GetDataDirectory()
	if err != nil {
		panic(err)
	}
	sessionPath = path.Join(dataDir, "session.neb")

	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return
	}

	dec := gob.NewDecoder(bytes.NewReader(data))
	if err = dec.Decode(&sessionId); err != nil {
		return
	}
}

var (
	user   *model.AuthUser
	userMu sync.Mutex

	sessionId   *uuid.UUID
	sessionPath string
)

func GetUser() *model.AuthUser {
	return user
}

func SetUser(newUser *model.AuthUser) {
	// themida.Macro(themida.EAGLE_BLACK_START)

	userMu.Lock()
	defer userMu.Unlock()
	user = newUser

	// themida.Macro(themida.EAGLE_BLACK_END)
}

func ModifyUser(cb func(*model.AuthUser)) {
	// themida.Macro(themida.EAGLE_BLACK_START)

	if user == nil {
		return
	}
	userMu.Lock()
	defer userMu.Unlock()
	cb(user)

	// themida.Macro(themida.EAGLE_BLACK_END)
}

func GetSession() *uuid.UUID {
	return sessionId
}

func SetSession(newSess *uuid.UUID) {
	// themida.Macro(themida.EAGLE_BLACK_START)

	sessionId = newSess

	var outBuf bytes.Buffer
	enc := gob.NewEncoder(&outBuf)
	if err := enc.Encode(sessionId); err == nil {
		_ = os.WriteFile(sessionPath, outBuf.Bytes(), os.ModePerm)
	}

	// themida.Macro(themida.EAGLE_BLACK_END)
}