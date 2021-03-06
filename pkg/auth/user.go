package auth

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/scylladb/go-set/u32set"

	"ccloud_hdd_server/pkg/use_hash"
	db_user "ccloud_hdd_server/pkg/db_sql/user"
)

type userCache struct {
	session map[uint32]context.Context
	key map[uint32][]byte
	cancel  map[uint32]context.CancelFunc
}

var (
	NotExistKeyError = errors.New("NotExistKeyError")
	FullLoginCountErr = errors.New("FullLoginCountErr")
)

var uc = func() userCache {
	return userCache{
		make(map[uint32]context.Context),
		make(map[uint32][]byte),
		make(map[uint32]context.CancelFunc),
	}
}()
var sessionKeySet = u32set.New()
var timeOut = time.Hour * 24
const max = 1

func Login(conn *sql.Conn, passwd []byte) (uint32, error) {
	u_id, err := db_user.GetUserId(conn, passwd)
	if err != nil {
		return 0, err
	}
	if len(uc.session) >= max {
		return 0,FullLoginCountErr
	}

	crypto_key, cerr := makeCryptoKey(conn,u_id,passwd)
	if cerr != nil {
		return 0,cerr
	}

	key := makeSessionKey()
	session := context.WithValue(context.Background(), "uid", u_id)
	sessTime, cancel := setSessionTime(session)


	uc.session[key] = sessTime
	uc.key[key] = crypto_key
	uc.cancel[key] = cancel

	return key, nil
}

func Logout(key uint32) error {
	if !sessionKeySet.Has(key) {
		return NotExistKeyError
	}
	uc.cancel[key]()

	deleteSessionKey(key)
	delete(uc.cancel, key)
	delete(uc.session, key)

	return nil
}
func GetUserCryptoKey(key uint32) []byte {
	return uc.key[key]
}
func GetUesrId(key uint32) (int, error) {
	if !sessionKeySet.Has(key) {
		return 0, NotExistKeyError
	}

	ctx := uc.session[key]
	return ctx.Value("uid").(int), nil
}
func GetLoginUserCount() int {return len(uc.session)}

func setSessionTime(p context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(p, timeOut)
}

func makeCryptoKey(conn *sql.Conn,id int,passwd []byte) ([]byte,error) {
	iv,err := db_user.GetUserIv(conn,id)
	if err != nil {return nil,err}
	raw := append(iv,passwd...)
	return use_hash.Sum(raw),nil
}

func makeSessionKey() uint32 {
	var ran_u uint32
	for {
		ran_u = rand.Uint32()
		if sessionKeySet.Has(ran_u) {
			continue
		}

		sessionKeySet.Add(ran_u)
		break
	}
	return ran_u
}
func deleteSessionKey(key uint32) { sessionKeySet.Remove(key) }
