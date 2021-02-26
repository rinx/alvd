module github.com/rinx/alvd

go 1.15

replace (
	cloud.google.com/go => cloud.google.com/go v0.76.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.37.3
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/chzyer/logex => github.com/chzyer/logex v1.1.11-0.20170329064859-445be9e134b2
	github.com/coreos/etcd => go.etcd.io/etcd v3.3.25+incompatible
	github.com/docker/docker => github.com/moby/moby v20.10.3+incompatible
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20210129204804-4364a4b9cfdd
	github.com/gogo/googleapis => github.com/gogo/googleapis v1.4.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp => github.com/google/go-cmp v0.5.4
	github.com/google/pprof => github.com/google/pprof v0.0.0-20210125172800-10e9aeb4a998
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.4
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.15.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/hailocab/go-hostpool => github.com/monzo/go-hostpool v0.0.0-20200724120130-287edbb29340
	github.com/klauspost/compress => github.com/klauspost/compress v1.11.8-0.20210203154158-6c96f3e2a592
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v2.1.2+incompatible
	github.com/vdaas/vald => ./vald
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	google.golang.org/grpc => google.golang.org/grpc v1.35.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.25.0
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
	k8s.io/metrics => k8s.io/metrics v0.20.2
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.8.1
)

require (
	cloud.google.com/go v0.74.0
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/prometheus v0.2.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/aws/aws-sdk-go v1.27.0
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-redis/redis/v8 v8.6.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.7.1
	github.com/gogo/googleapis v0.0.0-20180223154316-0cd9801be74a
	github.com/gogo/protobuf v1.3.1
	github.com/gogo/status v1.1.0
	github.com/google/go-cmp v0.5.4
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-version v1.2.1
	github.com/json-iterator/go v1.1.10
	github.com/klauspost/compress v0.0.0-00010101000000-000000000000
	github.com/kpango/fastime v1.0.16
	github.com/kpango/fuid v0.0.0-20200823100533-287aa95e0641
	github.com/kpango/gache v1.2.5
	github.com/kpango/glg v1.5.5
	github.com/pierrec/lz4/v3 v3.3.2
	github.com/rancher/remotedialer v0.2.5
	github.com/scylladb/gocqlx v1.5.0
	github.com/tensorflow/tensorflow v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v2 v2.3.0
	github.com/vdaas/vald v0.0.0-00010101000000-000000000000
	go.opencensus.io v0.23.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.17.0
	go.opentelemetry.io/otel/metric v0.17.0
	go.uber.org/automaxprocs v1.4.0
	go.uber.org/goleak v1.1.10
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7
	golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073
	google.golang.org/api v0.40.0
	google.golang.org/grpc v1.35.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/metrics v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
