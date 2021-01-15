module github.com/rinx/alvd

go 1.15

replace (
	cloud.google.com/go => cloud.google.com/go v0.66.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.34.25
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/chzyer/logex => github.com/chzyer/logex v1.1.11-0.20170329064859-445be9e134b2
	github.com/coreos/etcd => go.etcd.io/etcd v3.3.25+incompatible
	github.com/docker/docker => github.com/moby/moby v1.13.1
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20200815110948-5378c8f664e9
	github.com/gogo/googleapis => github.com/gogo/googleapis v1.4.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp => github.com/google/go-cmp v0.5.2
	github.com/google/pprof => github.com/google/pprof v0.0.0-20200905233945-acf8798be1f7
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.12.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/hailocab/go-hostpool => github.com/monzo/go-hostpool v0.0.0-20200724120130-287edbb29340
	github.com/klauspost/compress => github.com/klauspost/compress v1.11.1-0.20200908135004-a2bf5b1ec3aa
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v2.1.0+incompatible
	github.com/vdaas/vald => ./vald
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	google.golang.org/grpc => google.golang.org/grpc v1.32.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.25.0
	k8s.io/api => k8s.io/api v0.18.9
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.9
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.9
	k8s.io/client-go => k8s.io/client-go v0.18.9
	k8s.io/metrics => k8s.io/metrics v0.18.9
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.3
)

require (
	cloud.google.com/go v0.72.0
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/prometheus v0.2.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.4
	github.com/aws/aws-sdk-go v1.23.20
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.7.1
	github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp v0.5.4
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-version v1.2.1
	github.com/json-iterator/go v1.1.10
	github.com/klauspost/compress v0.0.0-00010101000000-000000000000
	github.com/kpango/fastime v1.0.16
	github.com/kpango/fuid v0.0.0-20200823100533-287aa95e0641
	github.com/kpango/gache v1.2.3
	github.com/kpango/glg v1.5.1
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/pierrec/lz4/v3 v3.3.2
	github.com/rancher/remotedialer v0.2.5
	github.com/scylladb/gocqlx v1.5.0
	github.com/tensorflow/tensorflow v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v2 v2.2.0
	github.com/vdaas/vald v0.0.62
	go.opencensus.io v0.22.5
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.15.0
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/goleak v1.1.10
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/sys v0.0.0-20210113181707-4bcb84eeeb78
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.34.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/metrics v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
