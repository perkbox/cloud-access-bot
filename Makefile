export GO111MODULE=on

NAME = cloud-access-bot
BIN_DIR := $(CURDIR)/bin

clean:
	rm -rf ./bin ./build ./dist

prepare:
	GOBIN=$(BIN_DIR) go install github.com/mitchellh/gox

build:
	$(BIN_DIR)/gox \
	-osarch="darwin/amd64 linux/amd64 windows/amd64" \
	-output "build/{{.OS}}-{{.Arch}}/$(NAME)" \
	${SOURCE_FILES}


dist:
	$(eval FILES := $(shell ls build/))
	mkdir dist
	for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		echo $$f; \
	done

git-release-artifacts: clean prepare build dist

