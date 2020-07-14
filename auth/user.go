package user

import (
	"context"
	"errors"
	"math/rand"
	"database/sql"
	"time"

	"github.com/scylladb/go-set/u32set"

	"ccloud_hdd_server/db"
)




type userCache struct {
	session map[uint32]context.Context
	cancel map[uint32]context.CancelFunc
}

var (
	NotExistKeyError = errors.New("NotExistKeyError")
)

var uc = func() userCache {
	return userCache{
		make(map[uint32]context.Context),
		make(map[uint32]context.CancelFunc),
	}
}()
var sessionKeySet = u32set.New()
var timeOut = time.Hour * 24

func Login(conn *sql.Conn,passwd []byte) (uint32,error) {
	u_id,err := db.GetUserId(conn,passwd)
	if err != nil {return 0,err}

	key := makeSessionKey()
	session := context.WithValue(context.Background(),"uid",u_id)
	sessTime,cancel := setSessionTime(session)

	uc.session[key] = sessTime
	uc.cancel[key] = cancel

	return key,nil
}

func Logout(key uint32) error {
	if !sessionKeySet.Has(key) {return NotExistKeyError}
	uc.cancel[key]()

	deleteSessionKey(key)
	delete(uc.cancel,key)
	delete(uc.session,key)

	return nil
}

func GetUesrId(key uint32) (int,error) {
	if !sessionKeySet.Has(key) {return 0,NotExistKeyError}
	
	ctx := uc.session[key]
	return ctx.Value("uid").(int),nil
}

func setSessionTime(p context.Context) (context.Context,context.CancelFunc) {
	return context.WithTimeout(p,timeOut)
}

func makeSessionKey() uint32 {
	var ran_u uint32
	for {
		ran_u = rand.Uint32()
		if sessionKeySet.Has(ran_u) {continue}

		sessionKeySet.Add(ran_u)
		break
	}
	return ran_u
}
func deleteSessionKey(key uint32) {sessionKeySet.Remove(key)}