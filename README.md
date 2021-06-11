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

USER="user123"
PASSWORD="secret123"
API_ENDPOINT="http://0.0.0.0:8000/api/v1/"

# test it by sending HTTP POST to "$API_ENDPOINT"
curl --user "$USER:$PASSWORD" --data-binary "@README.md" -i "$API_ENDPOINT"

# use supplied key to get back the payload, e.g.
curl --user "$USER:$PASSWORD" -i "${API_ENDPOINT}i1c9761ce00146fa5b7d812626442f604f1bf7ce97c8f9bab3afc325113f8fe5c"
```

