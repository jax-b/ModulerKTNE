#!/bin/bash

PROBUF_DIR="$HOME/tmp/protobuf"

if [ ! -d "$PROBUF_DIR" ]; then
    git clone https://github.com/protocolbuffers/protobuf.git $PROBUF_DIR
fi

protoc --go_out=. --go_opt=paths=source_relative wiremessage.proto -I $PROBUF_DIR/src -I .
protoc --go_out=. --go_opt=paths=source_relative modulestate.proto -I $PROBUF_DIR/src -I .
protoc --go_out=. --go_opt=paths=source_relative gamestate.proto -I $PROBUF_DIR/src -I .