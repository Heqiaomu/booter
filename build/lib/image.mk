# ==============================================================================
# Makefile helper functions for docker image
#

DOCKER := docker
DOCKER_SUPPORTED_VERSIONS ?= 17|18|19

# CI_REGISTRY ?= 10.1.19.19:8099 ## From CI builder
# Determine image files by looking into hack/docker/*.Dockerfile
IMAGE_FILES=$(wildcard ${ROOT_DIR}/build/k8s/image/*.Dockerfile)
# Determine images names by stripping out the dir names
IMAGES=$(foreach image,${IMAGE_FILES},$(subst .Dockerfile,,$(notdir ${image})))

ifeq (${IMAGES},)
  $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

DOCKER_EXIST := $(shell type $(DOCKER) >/dev/null 2>&1 || { echo >&1 "not installed"; })

.PHONY: image.build.verify
image.build.verify:
ifneq ($(DOCKER_EXIST),)
	$(error Please install docker($(DOCKER_SUPPORTED_VERSIONS)) first)
else
ifneq ($(shell $(DOCKER) -v | grep -q -E '\bversion ($(DOCKER_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported docker version. Please make install one of the following supported version: '$(DOCKER_SUPPORTED_VERSIONS)')
endif
endif
	@echo "===========> Docker version verification passed"

.PHONY: image.build
image.build: image.build.verify $(addprefix image.build., $(IMAGES))

.PHONY: image.push
image.push: image.build.verify $(addprefix image.push., $(IMAGES))

.PHONY: image.build.%
image.build.%:
	@echo "===========> Building $* $(VERSION) docker image"
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@cat $(ROOT_DIR)/build/k8s/image/$*.Dockerfile >tmp_$*.Dockerfile
	@cat $(ROOT_DIR)/build/k8s/image/$*.ignore >.dockerignore
	@$(DOCKER) build --pull -t $(CI_REGISTRY)/b20/$*:$(VERSION) -f tmp_$*.Dockerfile .
	@rm tmp_$*.Dockerfile
	@rm .dockerignore

.PHONY: image.push.%
image.push.%: image.build.%
	@echo "===========> Pushing $* $(VERSION) image to $(CI_REGISTRY)"
	@$(DOCKER) login -u $(CI_REGISTRY_USER) -p $(CI_REGISTRY_PASSWORD) $(CI_REGISTRY)
	@$(DOCKER) push $(CI_REGISTRY)/b20/$*:$(VERSION)
