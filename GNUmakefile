
BIN = protoc-gen-erlang

all: $(BIN)

$(BIN):
	go build .

test:
	go test -race .

clean:
	$(RM) $(BIN)

.PHONY: all build test clean
