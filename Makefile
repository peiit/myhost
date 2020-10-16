#初始化项目目录变量
HOMEDIR := $(shell pwd)
OUTDIR  := $(HOMEDIR)/output
OUTDIR_BIN  := $(OUTDIR)/bin
APPNAME := $(shell basename `pwd`)

#初始化命令变量
GO      := go
GOMOD   := $(GO) mod
GOBUILD := $(GO) build
GOTEST  := $(GO) test
GOCLEAN := $(GO) clean
GOPKGS  := $$($(GO) list ./...| grep -vE "vendor")

#产出结果
BIN_LINUX := $(APPNAME)_linux
BIN_ARM := $(APPNAME)_arm
BIN_DARWIN := $(APPNAME)_darwin

#执行编译，可使用命令 make 或 make all 执行, 顺序执行prepare -> compile -> test -> package 几个阶段
all: prepare compile test package


# prepare阶段, 下载非Go依赖，可单独执行命令: make prepare
prepare: prepare-dep
	$(shell rm -rf $(OUTDIR))
    $(shell mkdir -p $(OUTDIR_BIN))

prepare-dep:
	git config --global http.sslVerify false #设置git， 保证github mirror能够下载

set-env:
	$(GO) env -w GOPROXY=https://goproxy.cn
	$(GO) env -w GONOSUMDB=\*

#complile阶段，执行编译命令,可单独执行命令: make compile
compile:build

build: set-env
	$(GOMOD) tidy #下载Go依赖
	# Build for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUTDIR_BIN)/$(BIN_LINUX) main.go
	$(GOCLEAN)
	# Build for darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(OUTDIR_BIN)/$(BIN_DARWIN) main.go
	$(GOCLEAN)
	# Build for arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o $(OUTDIR_BIN)/$(BIN_ARM) main.go
	$(GOCLEAN)

#test阶段，进行单元测试， 可单独执行命令: make test
test: test-case

test-case: set-env
	$(GOTEST) -v -cover $(GOPKGS)

#与覆盖率平台打通，输出测试结果到文件中
#@$(GOTEST) -v -json -coverprofile=coverage.out $(GOPKGS) > testlog.out
#package阶段，对编译产出进行打包，输出到output目录, 可单独执行命令: make package

package: package-bin

package-bin:

#clean阶段，清除过程中的输出, 可单独执行命令: make clean
clean:
	rm -rf $(OUTDIR)

# avoid filename conflict and speed up build
.PHONY: all prepare compile test package  clean build
