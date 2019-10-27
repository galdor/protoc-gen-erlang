
BIN = protoc-gen-erlang

PROTO_DIR = example
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

PROTOC = protoc
PROTOC_FLAGS += --plugin $(BIN) --erlang_out .

all: $(BIN)

$(BIN): FORCE
	go build .

test:
	go test -race .

example: $(PROTO_FILES)
	$(PROTOC) $(PROTOC_FLAGS) $^

clean:
	$(RM) $(BIN) example/*.erl

FORCE:

.PHONY: all build test example clean
