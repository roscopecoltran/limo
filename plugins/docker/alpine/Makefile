.PHONY: all deps build

# determine platform
ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  PLATFORM 			:= MacOSX
  PLATFORM_SLUG 	:= macosx
  GO_BUILD_OS 		:= darwin
else
  PLATFORM 			:= Linux
  PLATFORM_SLUG 	:= linux
  GO_BUILD_OS 		:= $(PLATFORM_SLUG)
endif

# env variables shared with docker and make gnu for more clarity
# ENVFILE 					:= $(CURDIR)/*.config.env
ENVFILES 					:= $(shell find $(CURDIR) -name "*.config.env" -type f)

ccred=$(shell echo "\033[0;31m")
ccyellow=$(shell echo "\033[0;33m")
ccend=$(shell echo "\033[0m")
ccgreylight=$(shell echo "\033[0;35m")
cccyan=$(shell echo "\033[0;36m")
ccgreen=$(shell echo "\033[0;31m")

env-files: clear-screen
	@for f in $(shell find $(CURDIR) -name "*.config.env" -type f); do \
			echo "file: $$f" ; \
		done

include $(ENVFILES)
export

# app variables
APP_BIN_NAME=$(APP_BASENAME)_$(GIT_VERSION)

clear-screen:
	@clear
	@echo ""

intro: ## info about the Makefile purpose(s)
	@echo ""
	@echo "This Makefile is a simple example that demonstrates usual steps to build a binary that can be run in the same"
	@echo "architecture that was compiled in. The "ldflags" in the build assure that any needed dependency is included in the"
	@echo "binary and no external dependencies are needed to run the service."
	@echo ""

info: clear-screen intro vars/app vars/git vars/go vars/docker vars/stat

app-info: vars/app
git-info: vars/git
go-info: vars/app_go vars/go

vars/%: ## display variables key/value prefixed by a string (note. uppercase, lowercase only)	
	@$(foreach v, $(filter $(shell echo $* | tr '[:lower:]' '[:upper:]')%, $(.VARIABLES)), $(info $(shell printf "$(ccyellow)%32s$(ccend) $(cccyan)%s$(ccend)" "$(v)" "$($(v))" )))
	@$(foreach v, $(filter $(shell echo $* | tr '[:upper:]' '[:lower:]')%, $(.VARIABLES)), $(info $(shell printf "$(ccgreen)%32s$(ccend) $(ccred)%s$(ccend)" "$(v)" "$($(v))" )))

print-%:  clear-screen ## print user-defined variable key/value
    @echo '$*=$($*)'

default: clear-screen all ## default target

all:  clear-screen deps build

gox: golang-install-deps gox-xbuild ## install missing dependencies and cross-compile app for macosx, linux and windows platforms

darwin:
	@gox -verbose -os="darwin" -arch="amd64" -output="./bin/{{.Dir}}" ./cmd/...	

gox-darwin: ## generate all binaries for Mac/Apple platforms for the project with gox utility
	@gox -verbose -os="darwin" -arch="amd64" -output="$(APP_DIST_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}" $(APP_PKGS) # $(glide novendor)

gox-dist: ## generate all binaries for the project with gox utility
	@gox -verbose -os="darwin linux windows" -arch="amd64" -output="$(APP_DIST_DIR)/{{.OS}}/{{.Dir}}_{{.Os}}_{{.Arch}}" $(APP_PKGS) # $(glide novendor)

glide: glide-create glide-install ## install and manage all project dependencies via glide utility

glide-clean: ## clean glide utility cache (hint: check the contant of dirs available at \$GLIDE_TMP and \$GLIDE_HOME)
	@glide cc

glide-create: ## create the list of used dependencies in this golang project, via glide utility
	@if [ ! -f $(CURDIR)/glide.yaml ]; then glide create --non-interactive ; fi

glide-install: ## install app/pkg dependencies via glide utility
	@if [ -f $(CURDIR)/glide.yaml ]; then glide install --strip-vendor ; fi

golang-install-deps: golang-package-deps golang-embedding-deps golang-test-deps ## install global golang pkgs/deps

golang-package-deps: 
	@if [ ! -f $(GO_GOX) ]; then go get -v github.com/mitchellh/gox ; fi
	@if [ ! -f $(GO_GLIDE) ]; then go get -v github.com/Masterminds/glide ; fi

golang-embedding-deps: 
	@if [ ! -f $(GO_BINDATA) ]; then go get -v github.com/jteeuwen/go-bindata/... ; fi
	@if [ ! -f $(GO_BINDATA_ASSETFS) ]; then go get -v github.com/elazarl/go-bindata-assetfs/... ; fi

golang-test-deps: ## install unit-tests/debugging dependencies
	@if [ ! -f $(GO_VENDORCHECK) ]; then go get -u github.com/FiloSottile/vendorcheck ; fi
	@if [ ! -f $(GO_LINT) ]; then go get -u github.com/golang/lint/golint ; fi
	@if [ ! -f $(GO_DEP) ]; then go get -u github.com/golang/dep/cmd/dep ; fi
	@if [ ! -f $(GO_ERRCHECK) ]; then go get -u github.com/kisielk/errcheck ; fi
	@if [ ! -f $(GO_UNCONVERT) ]; then go get -u github.com/mdempsky/unconvert ; fi
	@if [ ! -f $(GO_INTERFACER) ]; then go get -u github.com/mvdan/interfacer/cmd/interfacer ; fi
	go get -u github.com/opennota/check/...
	go get -u github.com/yosssi/goat/...
	go get -u honnef.co/go/tools/...

# golang-go-github-fix:
# 	@if [ -d $(CURDIR)/vendor/github.com/google/go-github/github ]; then find . -type f -name "*.go" -exec sed -i 's/Starred/Topics/g' {} + ; fi

golang-fix-all: golang-fork-fix golang-logrus-fix

pkg-uri-clean:
	@rm -fR $(CURDIR)/vendor/github.com/roscopecoltran/sniperkit-sift
	@rm -fR $(CURDIR)/vendor/github.com/svent/sift

pkg-uri-fix: install-ag pkg-uri-clean clear-screen ## fix sniperkit-sift pkg uri for golang package import
	@echo "fix sniperkit-sift pkg uri for golang package import"
	@$(AG_EXEC) -l 'github.com/svent/sift' --ignore Makefile --ignore *.md . | xargs sed -i -e 's/svent\/sift/roscopecoltran\/sniperkit-sift/g'
	@find . -name "*-e" -exec rm -f {} \; 

pkg-uri-revert: install-ag pkg-uri-clean clear-screen ## fix sift, fork, pkg uri for golang package import
	@echo "fix sift, fork, pkg uri for golang package import"
	@$(AG_EXEC) -l 'github.com/roscopecoltran/sniperkit-sift' --ignore Makefile --ignore *.md . | xargs sed -i -e 's/roscopecoltran\/sniperkit-sift/svent\/sift/g'
	@find . -name "*-e" -exec rm -f {} \;

golang-logrus-fix: install-ag clear-screen ## fix logrus case for golang package import
	@if [ -d $(CURDIR)/vendor/github.com/Sirupsen ]; then rm -fr vendor/github.com/Sirupsen ; fi
	@$(AG_EXEC) -l 'github.com/Sirupsen/logrus' vendor | xargs sed -i 's/Sirupsen/sirupsen/g'

go-github-fix:
	@if [ -d ./vendor/github.com/google/go-github/github ]; then find ./vendor/github.com/google/go-github/github -name activity_star.go -exec sed -i 's/mediaTypeStarringPreview/mediaTypeTopicsPreview/g' {} + ; fi

install-ag: install-ag-$(PLATFORM_SLUG) ## install the silver searcher (aka. ag)

# if [ "$choice" == 'y' ] && [ "$choice1" == 'y' ]; then
install-ag-macosx: clear-screen ## install the silver searcher on Apple/MacOSX platforms
	@echo "install the silver searcher on Apple/MacOSX platforms"
	@if [ -f $(PKG_BREW_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 		then $(PKG_BREW_EXEC) install the_silver_searcher; fi 
	@if [ -f $(PKG_MACPORTS_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 	then $(PKG_MACPORTS_EXEC) install the_silver_searcher ; fi	

install-ag-linux: clear-screen ## install the silver searcher on Linux platforms
	@echo "install the silver searcher on Linux platforms"
	@if [ -f $(PKG_APK_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 			then $(PKG_APK_EXEC) add --no-cache --update the_silver_searcher ; fi 
	@if [ -f $(PKG_APT_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 			then $(PKG_APT_EXEC) install -f --no-recommend silversearcher-ag ; fi 
	@if [ -f $(PKG_YUM_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 			then $(PKG_YUM_EXEC) install the_silver_searcher ; fi
	@if [ -f $(PKG_DNF_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 			then $(PKG_DNF_EXEC) install the_silver_searcher ; fi
	@if [ -f $(PKG_EMERGE_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 		then $(PKG_EMERGE_EXEC) -a sys-apps/the_silver_searcher ; fi
	@if [ -f $(PKG_PACMAN_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 		then $(PKG_PACMAN_EXEC) -S the_silver_searcher ; fi
	@if [ -f $(PKG_SLACKWARE_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 	then $(PKG_SLACKWARE_EXEC) -i the_silver_searcher ; fi
	@if [ -f $(PKG_ZYPPER_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 		then $(PKG_ZYPPER_EXEC) install the_silver_searcher ; fi
	@if [ -f $(PKG_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 			then $(PKG_EXEC) install the_silver_searcher ; fi
	@if [ -f $(PKG_ADD_EXEC) ] && [ ! -f $(SH_AG_EXEC) ]; 		then $(PKG_ADD_EXEC) the_silver_searcher ; fi

push-tag:
	$(GIT_EXEC) checkout ${APP_BRANCH}
	$(GIT_EXEC) pull origin ${APP_BRANCH}
	$(GIT_EXEC) tag ${GIT_VERSION}
	$(GIT_EXEC) push origin ${APP_BRANCH} --tags

# golang-go-github-fix:
# 	@if [ -d $(CURDIR)/vendor/github.com/google/go-github/github ]; then find . -type f -name "*.go" -exec sed -i 's/Starred/Topics/g' {} + ; fi

prepare: git golang-install-deps vars/app ## prepare required destination folders
	@mkdir -p $(APP_DIST_DIR) $(APP_DIST_DIR)/darwin $(APP_DIST_DIR)/linux $(APP_DIST_DIR)/windows

build: vars/app ## build a local executable with the git commit version
	go build -a -ldflags="-X $(APP_VCS_URI)/version.SnkPluginApkFiledVersion=${GIT_VERSION}" -o $(BIN_NAME) $(SNK_PLUGIN_GOLANG_CMD_DIR)
	@echo "You can now use ./${APP_BIN_NAME}"

travis: vars/travis install-deps golang-install-deps golang-logrus-fix ## travis-ci builc
	@echo "building..."
	@go build

STAT_WGET := $(shell command -v wget 2> /dev/null)
STAT_CURL := $(shell command -v curl 2> /dev/null)
ifneq ("$(wildcard $(CURDIR)/.UUID)","")
STAT_UUID := $(shell cat $(CURDIR)/.UUID | base64)
else
STAT_UUID := $(shell cat /dev/urandom | LC_ALL=C tr -dc A-Z0-9 | head -c32 > $(CURDIR)/.UUID && cat $(CURDIR)/.UUID | base64)
endif
STAT_UNAME := $(shell uname -srmp | base64)
STAT_PATH := $(shell git rev-parse --show-prefix | base64)
ifeq ("$(STAT_PATH)","Cg==")
STAT_PATH := $(shell echo $$(cd .. && git rev-parse --show-prefix)$$(basename `pwd`) | base64)
endif
STAT_GIT_VERSION := $(shell git describe --abbrev=8 --dirty --always --tags | base64)
STAT_GIT_REMOTE := $(shell git config --get remote.origin.url | base64)

# STAT_URL := http://172.16.121.1:8888/statistics/
# STAT_DATA := "uuid=$(STAT_UUID)&uname=$(STAT_UNAME)&path=$(STAT_PATH)&version=$(STAT_GIT_VERSION)&remote=$(STAT_GIT_REMOTE)"

help: ## show this help
	@printf "\033[1;34mCleaning \033[0m\033[34mup...\033[0m"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@printf "\033[1;34mCleaning \033[0m\033[34mup...\033[0m"

