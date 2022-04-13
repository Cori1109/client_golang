module github.com/prometheus/client_golang

require (
	github.com/beorn7/perks v1.0.1
	github.com/cespare/xxhash/v2 v2.1.2
	github.com/davecgh/go-spew v1.1.1
	github.com/golang/protobuf v1.5.2
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.33.0
	github.com/prometheus/procfs v0.7.3
	golang.org/x/sys v0.0.0-20220328115105-d36c6a25d886
	google.golang.org/protobuf v1.28.0
)

exclude github.com/prometheus/client_golang v1.12.1

go 1.16
