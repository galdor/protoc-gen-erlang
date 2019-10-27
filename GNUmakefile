
BIN = protoc-gen-erlang

all: $(BIN)

$(BIN): FORCE
	go build .

test:
	go test -race .

clean:
	$(RM) $(BIN)

FORCE:

.PHONY: all build test clean
