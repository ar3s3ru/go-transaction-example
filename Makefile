GO ?= go
TARGET ?= app
BUILD_DIR ?= ${PWD}/build

TARGET_PKG = ./cmd/${TARGET}
TARGET_BUILD = ${BUILD_DIR}/${TARGET}

info = echo -e "\033[0;32m${1}\033[0m"

install: build
	${call info,Moving app into /bin...}
	@mv ${TARGET_BUILD} /bin

run:
	${call info,Running app...}
	@${GO} run ${APP_PKG}

build:
	${call info,Building app in "${BUILD_DIR}"...}
	@${GO} build -o ${TARGET_BUILD} ${TARGET_PKG}

clean:
	${call info,Cleaning app from /bin and ${BUILD_DIR}...}
	@rm ${TARGET_BUILD}
	@rm /bin/${TARGET}

.PHONY: run build install clean
