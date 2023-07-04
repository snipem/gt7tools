GO := go
BINDIR := build

SOURCES := $(shell find cmd -name '*.go')
BINARY_NAMES := $(patsubst cmd/%/,%,$(dir $(SOURCES)))

.PHONY: all windows darwin_arm64 darwin_amd64 linux clean

all: windows darwin_arm64 darwin_amd64 linux

windows: $(patsubst %,$(BINDIR)/%.exe,$(BINARY_NAMES))

darwin_arm64: $(patsubst %,$(BINDIR)/%-darwin-arm64,$(BINARY_NAMES))

darwin_amd64: $(patsubst %,$(BINDIR)/%-darwin-amd64,$(BINARY_NAMES))

linux: $(patsubst %,$(BINDIR)/%-linux,$(BINARY_NAMES))

$(BINDIR)/%.exe: cmd/%
	@mkdir -p $(dir $@)
	GOOS=windows GOARCH=amd64 $(GO) build -o $@ ./$<

$(BINDIR)/%-darwin-arm64: cmd/%
	@mkdir -p $(dir $@)
	GOOS=darwin GOARCH=arm64 $(GO) build -o $@ ./$<

$(BINDIR)/%-darwin-amd64: cmd/%
	@mkdir -p $(dir $@)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $@ ./$<

$(BINDIR)/%-linux: cmd/%
	@mkdir -p $(dir $@)
	GOOS=linux GOARCH=amd64 $(GO) build -o $@ ./$<

clean:
	rm -rf $(BINDIR)
