package cache

import (
	"time"

	"ti-ticket/DAO"
	"ti-ticket/utils"

	"github.com/twinj/uuid"
)

const _expire_gap = time.Minute * 3

var (
	cache [32]User
	cp    int = 0
)

func Init() error {
	cp = 0
	return nil
}

func AddUser(account string) (*User, bool) {
	if cp >= 32 {
		return nil, false
	}
	cache[cp].Uid = uuid.NewV4().String()
	cache[cp].Account = account
	refreshUser(&cache[cp])
	cp++
	return &cache[cp-1], true
}

func FetchUser(uid string) (*User, bool) {
	for it := 0; it < cp; it++ {
		if cache[it].Uid == uid {
			refreshUser(&cache[it])
			return &cache[it], true
		}
	}
	return nil, false
}

func FetchUserByAccount(account string) (*User, bool) {
	for it := 0; it < cp; it++ {
		if cache[it].Account == account {
			refreshUser(&cache[it])
			return &cache[it], true
		}
	}
	return nil, false
}

func DropUser(u *User) {
	DAO.DeleteUser((*u).Account)
	// Remove from cache
}

type User struct {
	Uid         string
	Account     string
	Password    string
	Expire_time int64
}

func refreshUser(up *User) {
	passwd := utils.GetSecret(16)
	(*up).Password = passwd
	(*up).Expire_time = time.Now().Add(_expire_gap).UTC().Unix()
	DAO.UpdateUser((*up).Account, (*up).Password)
}
