# sha256msg

## Description

Stores blobs `[]byte` on S3 by their `sha256` hash.

### API

HTTP POST `/api/v1/`:
- Takes HTTP body `[]byte` as an input a message. Stores the message and returns the SHA256 hash of that message.

HTTP GET `/api/v1/sha256hash`:
- returns `[]byte` message associated with that hash

### Implementation details

Storage system is abstracted with the following interface.

```
type store interface {
    Fetch(ctx context.Context, objectKey string) (io.ReadCloser, error)
    Store(ctx context.Context, r io.ReadSeeker) (string, error)
}
```

Let's say, one needs to add Azure, GCP or IPFS-backed API, then the only thing they need is to implement `Fetch` and `Store` methods for that storage provider.
Use `internal/store/s3store/s3.go` as an example.


### Local setup

Check `docker-compose.test.yml` for username, secret and server endpoint.

```
docker-compose -f docker-compose.test.yml up -d --no-deps --build

USER="username"
PASSWORD="secret"
API_ENDPOINT="http://0.0.0.0:8000/api/v1"

curl --user username:secret --data-binary "@filename_here" -i "$API_ENDPOINT"
curl --user "$USER:$PASSWORD" -i "$API_ENDPOINT/sha256key"
```
