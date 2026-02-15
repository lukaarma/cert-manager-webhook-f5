GO ?= $(shell which go)
OS ?= $(shell $(GO) env GOOS)
ARCH ?= $(shell $(GO) env GOARCH)

IMAGE_REGISTRY = "docker.io"
IMAGE_NAME := "lukaarma/cert-manager-webhook-f5"
IMAGE_TAG := "latest"

OUT := $(shell pwd)/_out

KUBEBUILDER_VERSION=1.28.0

test: _test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)/etcd _test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)/kube-apiserver _test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)/kubectl
	TEST_ASSET_ETCD=_test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)/etcd \
	TEST_ASSET_KUBE_APISERVER=_test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)/kube-apiserver \
	TEST_ASSET_KUBECTL=_test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)/kubectl \
	$(GO) test -v .

_test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH).tar.gz: | _test
	curl -fsSL https://go.kubebuilder.io/test-tools/$(KUBEBUILDER_VERSION)/$(OS)/$(ARCH) -o $@

_test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)/etcd _test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)/kube-apiserver _test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)/kubectl: _test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH).tar.gz | _test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH)
	tar xfO $< kubebuilder/bin/$(notdir $@) > $@ && chmod +x $@

.PHONY: clean
clean:
	rm -r _test $(OUT)

.PHONY: build
build:
	docker build -t "$(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)" .

.PHONY: rendered-manifest.yaml
rendered-manifest.yaml:
	helm template \
	    cert-manager-webhook-f5 \
        --set image.repository=$(IMAGE_NAME) \
        --set image.tag=$(IMAGE_TAG) \
		--namespace cert-manager \
        deploy/cert-manager-webhook-f5 > "$(OUT)/rendered-manifest.yaml"

_test $(OUT) _test/kubebuilder-$(KUBEBUILDER_VERSION)-$(OS)-$(ARCH):
	mkdir -p $@
