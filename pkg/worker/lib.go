package worker

import (
	"ccloud_hdd_server/pkg/auth"
	"ccloud_hdd_server/pkg/db_sql"
	"ccloud_hdd_server/pkg/get_db"
	"context"
	"database/sql"
	"net/http"
	"strconv"

	db_user "ccloud_hdd_server/pkg/db_sql/user"
)

type Worker interface {
	Do(w http.ResponseWriter, r *http.Request, key []byte)
}


func authAndGetInfo(r *http.Request) (userkey int, dbConn *sql.Conn, iv []byte, basePath string, err error) {
	ck, ck_err := r.Cookie("session")
	if ck_err != nil {
		return 0, nil, nil, "", ck_err
	}

	using_key, atol_err := strconv.Atoi(ck.Value)
	if atol_err != nil {
		return 0, nil, nil, "", atol_err
	}

	userkey, ck_err = auth.GetUesrId(uint32(using_key))
	if ck_err != nil {
		return 0, nil, nil, "", ck_err
	}

	dbConn, ck_err = get_db.GetDbConn(context.Background())
	if ck_err != nil {
		return 0, nil, nil, "", ck_err
	}

	iv, ck_err = db_user.GetUserIv(dbConn, using_key)
	if ck_err != nil {
		return 0, nil, nil, "", ck_err
	}

	base_id, base_err := db_user.GetBasePathId(dbConn, using_key)

	basePath, base_err = db_sql.GetBasePath(dbConn, base_id)
	if base_err != nil {
		return 0, nil, nil, "", base_err
	}
	err = nil
	return
}
