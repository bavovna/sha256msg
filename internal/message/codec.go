package message

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mkorenkov/sha256msg/internal/requestcontext"
	"github.com/mkorenkov/sha256msg/internal/store/s3store"
	"github.com/pkg/errors"
)

//StoreHandler HTTP Handler, that takes HTTP body, stores `[]byte` and returns sha256 hash
func StoreHandler(w http.ResponseWriter, r *http.Request) {
	rctx := requestcontext.FromRequest(r)
	if r.ContentLength == 0 || r.ContentLength > rctx.Config.MaxUploadSizeBytes {
		httpError(w, nil, `{"error": "bad payload"}`, http.StatusBadRequest)
		return
	}

	client, err := rctx.Config.NewS3Client(r.Context())
	if err != nil {
		internalServerError(w, errors.Wrap(err, "s3 client creation error"))
		return
	}

	store := &s3store.API{
		Conf:   rctx.Config,
		Client: client,
	}

	f, err := ioutil.TempFile(rctx.Config.UploadDir, "sha256msg")
	if err != nil {
		internalServerError(w, errors.Wrap(err, "tempfile error"))
		return
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		internalServerError(w, errors.Wrap(err, "error copying req body"))
		return
	}

	key, err := store.Store(r.Context(), f)
	if err != nil {
		internalServerError(w, errors.Wrap(err, "s3 error"))
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"key": "%s"}`, key)))
	if err != nil {
		internalServerError(w, errors.Wrap(err, "write error"))
		return
	}
}

//FetchHandler HTTP Handler, that returns `[]byte` based on sha256 hash key
func FetchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectKey, ok := vars["key"]
	if !ok {
		internalServerError(w, errors.New("no objectKey found"))
		return
	}
	if len(objectKey) != 64 {
		httpError(w, nil, `{"error": "bad key"}`, http.StatusBadRequest)
		return
	}

	rctx := requestcontext.FromRequest(r)
	client, err := rctx.Config.NewS3Client(r.Context())
	if err != nil {
		internalServerError(w, errors.Wrap(err, "s3 client creation error"))
		return
	}

	store := &s3store.API{
		Conf:   rctx.Config,
		Client: client,
	}

	blobReadCloser, err := store.Fetch(r.Context(), objectKey)
	if err != nil {
		if errors.Is(err, s3store.ErrorNoResults) {
			httpError(w, nil, `{"error": "not found"}`, http.StatusNotFound)
			return
		}
		internalServerError(w, errors.Wrap(err, "s3 error"))
		return
	}
	defer blobReadCloser.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="`+objectKey+`"`)

	_, err = io.Copy(w, blobReadCloser)
	if err != nil {
		internalServerError(w, errors.Wrap(err, "client connection error"))
		return
	}
}

func httpError(w http.ResponseWriter, err error, output string, status int) {
	if err != nil {
		log.Printf("[ERROR] StoreHandler %+v\n", err)
	}
	if status != 0 {
		w.WriteHeader(status)
	}
	if output != "" {
		_, _ = w.Write([]byte(output + "\n"))
	}
}

func internalServerError(w http.ResponseWriter, err error) {
	httpError(w, err, `{"error": "internal server error"}`, http.StatusInternalServerError)
}
