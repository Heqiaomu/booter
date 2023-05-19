# ==============================================================================
# Makefile helper functions for common
#

SHELL := /bin/bash

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif
ifeq ($(origin OUTPUT_PREFIX),undefined)
OUTPUT_PREFIX := output
endif
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/$(OUTPUT_PREFIX)
endif
ifeq ($(origin TOOLS_DIR),undefined)
TOOLS_DIR := $(OUTPUT_DIR)/tools
endif
ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
endif

# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --dirty --always --tags | sed 's/-/./2' | sed 's/-/./2' )
endif
export VERSION

# set the git tag number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin GIT_TAG), undefined)
GIT_TAG := $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
endif
export GIT_TAG

# set the git commit hash. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin GIT_COMMIT), undefined)
GIT_COMMIT := $(shell git log --pretty=format:'%H' -n 1)
endif
export GIT_COMMIT

# set the git tree state. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin GIT_TREE_STATE), undefined)
GIT_TREE_STATE := $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
endif
export GIT_TREE_STATE

COMMA := ,
SPACE :=
SPACE +=
