##################################
# booter make helper
#

# booter 二进制名
BOOTER := booter

$(OUTPUT_DIR)/$(OS)/$(ARCH)/$(GO_APP_PRE)$(BOOTER)$(GO_OUT_EXT):
	$(if $(wildcard $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(GO_APP_PRE)$(BOOTER)$(GO_OUT_EXT),,cd $(shell $(MAKE) build.booter)))

# % = Action.Platform.Command
# booter usage:
# ./booter --dir=$(dir) --work=$(target workdir) <action> <app>
# action: start/stop/restart/status/tail
# app: <app-name>/all
# 如果make build.booter未执行过,则自动执行make build.booter
.PHONY: go.booter.%
go.booter.%:
	$(eval COMMAND := $(word 3,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 2,$(subst ., ,$*)))
	$(eval ACTION := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@if [ A$(wildcard $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(GO_APP_PRE)$(BOOTER)$(GO_OUT_EXT)) != A$(OUTPUT_DIR)/$(OS)/$(ARCH)/$(GO_APP_PRE)$(BOOTER)$(GO_OUT_EXT) ]; \
	then \
		$(MAKE) build.booter; \
	fi
	@echo "===========> $(ACTION) binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@$(OUTPUT_DIR)/$(OS)/$(ARCH)/$(GO_APP_PRE)$(BOOTER)$(GO_OUT_EXT) --dir=$(OUTPUT_DIR)/$(OS)/$(ARCH) --work=$(ROOT_DIR)/cmd $(ACTION) $(COMMAND)


	
