# ==============================================================================
# Makefile helper functions for golang
#

GO := go
GO_SUPPORTED_VERSIONS ?= 1.11|1.12|1.13|1.14|1.15|1.16
# ldflags 设置version包版本参数 -s省略符号表和调试信息 -w省略DWARF符号表，都可以缩小二进制大小
GO_LDFLAGS += -X $(VERSION_PACKAGE).gitTag=$(GIT_TAG) -X $(VERSION_PACKAGE).version=$(VERSION) -X $(VERSION_PACKAGE).buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') -X $(VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) -X $(VERSION_PACKAGE).gitTreeState=$(GIT_TREE_STATE) -s -w

# blocface应用名前缀
GO_APP_PRE := blocface-

ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS 2>/dev/null)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH 2>/dev/null)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
endif

GOPATH := $(shell go env GOPATH 2>/dev/null)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

PLATFORMS ?= darwin_amd64 windows_amd64 linux_amd64
COMMANDS ?= $(wildcard ${ROOT_DIR}/cmd/*)
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))

ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

GO_EXIST := $(shell type $(GO) >/dev/null 2>&1 || { echo >&1 "not installed"; })

.PHONY: go.build.verify
go.build.verify:
ifneq ($(GO_EXIST),)
	$(error Please install go($(GO_SUPPORTED_VERSIONS)) first)
else
ifneq ($(shell $(GO) version | grep -q -E '\bgo($(GO_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif
endif
	@echo "===========> go version verification passed"

.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@mkdir -p $(OUTPUT_DIR)/$(OS)/$(ARCH)
	@CGO_ENABLED=1 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build -o $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(GO_APP_PRE)$(COMMAND)$(GO_OUT_EXT) -ldflags "$(GO_LDFLAGS)" $(ROOT_PACKAGE)/cmd/$(COMMAND)

.PHONY: go.build
go.build: go.build.verify $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))

# 运行指定应用
.PHONY: go.exec.%
go.exec.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Running binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@$(OUTPUT_DIR)/$(OS)/$(ARCH)/$(GO_APP_PRE)$(COMMAND)$(GO_OUT_EXT) --dir=$(ROOT_DIR)/cmd/$(COMMAND) start --nodaemon

# 运行单个应用
.PHONY: go.run.%
go.run.%: 
	$(eval bin := $(word 1,$(subst ., ,$*))) 
	@$(MAKE) $(addprefix go.exec., $(addprefix $(PLATFORM)., $(bin)))

# for build bin single
# 编译单个app
.PHONY: go.app.%
go.app.%: go.build.verify 
	$(eval bin := $(word 1,$(subst ., ,$*))) 
	@$(MAKE) $(addprefix go.build., $(addprefix $(PLATFORM)., $(bin)))

.PHONY: go.build.all
go.build.all: go.build.verify $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINS))))

.PHONY: go.clean
go.clean:
	@echo "===========> Cleaning all build output"
	@rm -rf $(OUTPUT_DIR)

.PHONY: go.lint.verify
go.lint.verify: go.build.verify
ifeq (,$(wildcard $(GOBIN)/golangci-lint))
	@echo "===========> Installing golangci-lint"
	@GO111MODULE=on $(GO) get github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint
endif

.PHONY: go.lint
go.lint: go.lint.verify
	@echo "===========> Run golangci-lint to lint source codes"
	@$(GOBIN)/golangci-lint run --fix

## 暂时不用go-junit
# .PHONY: go.test.verify
# go.test.verify: go.build.verify
# ifeq (,$(wildcard $(GOBIN)/go-junit-report))
# 	@echo "===========> Installing go-junit-report"
# 	@GO111MODULE=on $(GO) get github.com/jstemmer/go-junit-report
# 	@$(GO) install github.com/jstemmer/go-junit-report
# endif

## 手动计算覆盖率 TODO目前计算有误差,因为只计算了有包含_test.go文件的包
# @$(GO) test -gcflags=-l -count=1 -timeout=10m -short -v ./... 2>&1 | tee >($(GOBIN)/go-junit-report --set-exit-code >$(OUTPUT_DIR)/report.xml)
.PHONY: go.test
go.test: #go.test.verify
	@echo "===========> Run unit test"
	@mkdir -p $(OUTPUT_DIR)
	@sh $(ROOT_DIR)/script/gotest.sh
