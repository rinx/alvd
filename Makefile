ORG = rinx
ALTORG = ghcr.io/rinx
REPO = alvd

VERSION ?= unknown

VALD_DIR = vald
VALD_REPO = vdaas/vald
VALD_BRANCH = master
VALD_DEPTH = 1

NGT_VERSION = 1.13.1

GO_VERSION = $(shell go version)
GOARCH = $(shell go env GOARCH)

ifeq ($(GOARCH),amd64)
CFLAGS ?= -mno-avx512f -mno-avx512dq -mno-avx512cd -mno-avx512bw -mno-avx512vl
CXXFLAGS ?= $(CFLAGS)
EXTLDFLAGS ?= -m64
else ifeq ($(GOARCH),arm64)
CFLAGS ?=
CXXFLAGS ?= $(CFLAGS)
EXTLDFLAGS ?= -march=armv8-a
else
CFLAGS ?=
CXXFLAGS ?= $(CFLAGS)
EXTLDFLAGS ?=
endif

DOCKER ?= docker
DOCKER_OPTS ?=

NOAVX_TAG ?= noavx
AVX2_TAG ?= avx2

NGT_BUILD_OPTIONS ?= -DNGT_AVX_DISABLED=ON

.PHONY:
all: build

.PHONY: clean
clean:
	rm -rf \
	    cmd/alvd/alvd \
	    internal \
	    pkg/vald/agent/ngt \
	    $(VALD_DIR)

.PHONY: build
build: \
	cmd/alvd/alvd

.PHONY: docker/build
docker/build: \
	docker/build/noavx \
	docker/build/avx2

.PHONY: docker/build/noavx/name
docker/build/noavx/name:
	@echo $(ORG)/$(REPO):$(NOAVX_TAG)

.PHONY: docker/build/noavx/altname
docker/build/noavx/altname:
	@echo $(ALTORG)/$(REPO):$(NOAVX_TAG)

.PHONY: docker/build/noavx
docker/build/noavx:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -t $(ORG)/$(REPO):$(NOAVX_TAG) \
	    -t $(ALTORG)/$(REPO):$(NOAVX_TAG) . \
	    --build-arg NGT_BUILD_OPTIONS="-DNGT_AVX_DISABLED=ON" \
	    --build-arg VERSION="$(VERSION)"

.PHONY: docker/build/avx2/name
docker/build/avx2/name:
	@echo $(ORG)/$(REPO):$(AVX2_TAG)

.PHONY: docker/build/avx2/altname
docker/build/avx2/altname:
	@echo $(ALTORG)/$(REPO):$(AVX2_TAG)

.PHONY: docker/build/avx2
docker/build/avx2:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -t $(ORG)/$(REPO):$(AVX2_TAG) \
	    -t $(ALTORG)/$(REPO):$(AVX2_TAG) . \
	    --build-arg NGT_BUILD_OPTIONS="" \
	    --build-arg VERSION="$(VERSION)"

cmd/alvd/alvd: \
	ngt/install \
	internal \
	pkg/vald/agent/ngt \
	$(shell find ./cmd/alvd -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export CGO_CXXFLAGS="-g -Ofast -march=native" \
	    && export CGO_FFLAGS="-g -Ofast -march=native" \
	    && export CGO_LDFLAGS="-g -Ofast -march=native" \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static -fPIC -pthread -fopenmp -std=c++17 -lstdc++ -lm $(EXTLDFLAGS)' \
		-X 'main.Version=$(VERSION)' \
		-X 'main.GoVersion=$(GO_VERSION)' \
		-X 'main.NGTVersion=$(NGT_VERSION)'" \
	    -a \
	    -tags "cgo netgo" \
	    -trimpath \
	    -installsuffix "cgo netgo" \
	    -o $@ \
	    $(dir $@)main.go

internal: $(VALD_DIR)
	mkdir -p $(dir $@)
	cp -r $(VALD_DIR)/$@ $@
	find $@ -type f -name "*.go" | xargs sed -i "s:$(VALD_REPO)/internal:$(ORG)/$(REPO)/internal:g"
	find $@ -type f -name "*.go" | xargs sed -i "s:$(VALD_REPO)/pkg/agent/internal:$(ORG)/$(REPO)/pkg/vald/agent/internal:g"
	find $@ -type f -name "*.go" | xargs sed -i "s:$(VALD_REPO)/pkg/agent/core/ngt:$(ORG)/$(REPO)/pkg/vald/agent/ngt:g"

pkg/vald/agent/ngt: \
	$(VALD_DIR) \
	pkg/vald/agent/internal
	mkdir -p $(dir $@)
	cp -r $(VALD_DIR)/pkg/agent/core/ngt $@
	find $@ -type f -name "*.go" | xargs sed -i "s:$(VALD_REPO)/internal:$(ORG)/$(REPO)/internal:g"
	find $@ -type f -name "*.go" | xargs sed -i "s:$(VALD_REPO)/pkg/agent/internal:$(ORG)/$(REPO)/pkg/vald/agent/internal:g"
	find $@ -type f -name "*.go" | xargs sed -i "s:$(VALD_REPO)/pkg/agent/core/ngt:$(ORG)/$(REPO)/pkg/vald/agent/ngt:g"

pkg/vald/agent/internal: $(VALD_DIR)
	mkdir -p $(dir $@)
	cp -r $(VALD_DIR)/pkg/agent/internal $@
	find $@ -type f -name "*.go" | xargs sed -i "s:$(VALD_REPO)/internal:$(ORG)/$(REPO)/internal:g"


$(VALD_DIR):
	git clone \
	    --depth $(VALD_DEPTH) \
	    -b $(VALD_BRANCH) \
	    https://github.com/$(VALD_REPO) \
	    $(VALD_DIR)

.PHONY: ngt/install
## install NGT
ngt/install: /usr/local/include/NGT/Capi.h
/usr/local/include/NGT/Capi.h:
	curl -LO https://github.com/yahoojapan/NGT/archive/v$(NGT_VERSION).tar.gz
	tar zxf v$(NGT_VERSION).tar.gz -C /tmp
	cd /tmp/NGT-$(NGT_VERSION) && \
	    cmake \
	    -DCMAKE_C_FLAGS="$(CFLAGS)" \
	    -DCMAKE_CXX_FLAGS="$(CXXFLAGS)" \
	    $(NGT_BUILD_OPTIONS) \
	    .
	make -j -C /tmp/NGT-$(NGT_VERSION)
	make install -C /tmp/NGT-$(NGT_VERSION)
	rm -rf v$(NGT_VERSION).tar.gz
	rm -rf /tmp/NGT-$(NGT_VERSION)
	ldconfig
