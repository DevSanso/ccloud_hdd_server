package worker

import (
	"net/http"
)


type Worker interface {
	Do(w http.ResponseWriter,r *http.Request,next Worker)
}

