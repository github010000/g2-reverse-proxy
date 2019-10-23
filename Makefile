VERSION := $(shell cat VERSION.txt)

#---

build:
	@echo "\n\033[1;33m+ $@\033[0m"
	@env GOOS=linux GOARCH=amd64 go build -v -o out/g2-reverse-proxy cmd/g2_reverse_proxy/g2_reverse_proxy.go

#---

define build_container
docker build -t au.icr.io/$(1)/g2-reverse-proxy:$(VERSION) .
endef

build-container-dev: build
	@echo "\n\033[1;33m+ $@\033[0m"
	@$(call build_container,"zmondev")

build-container-prod: build
	@echo "\n\033[1;33m+ $@\033[0m"
	@$(call build_container,"zmon")

#---

define push
docker push au.icr.io/$(1)/g2-reverse-proxy:$(VERSION)
endef

push-dev: build-container-dev
	@echo "\n\033[1;33m+ $@\033[0m"
	@$(call push,"zmondev")

push-prod: build-container-prod
	@echo "\n\033[1;33m+ $@\033[0m"
	@$(call push,"zmon")

#---

define k8s_deploy
jsonnet \
--ext-str version="$(VERSION)" \
--ext-str env="$(1)" \
docs/k8s/deploy-template.jsonnet
endef

define k8s_secret
jsonnet \
--ext-str env="$(1)" \
docs/k8s/secret-template.jsonnet
endef

k8s-deployment-dev:
	@$(call k8s_deploy,"dev")

k8s-deployment-prod:
	@$(call k8s_deploy,"prod")

k8s-secret-dev:
	@$(call k8s_secret,"dev")

k8s-secret-prod:
	@$(call k8s_secret,"prod")

#---

test:
	@go test -short ./...