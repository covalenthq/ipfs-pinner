GO = go # if using docker, should not need to be installed/linked
GOBIN = $(CURDIR)/build/bin/
DBG_CGO_CFLAGS += -DMDBX_DEBUG=1



GO_DBG_BUILD = $(DBG_CGO_CFLAGS) $(GO) build -tags $(BUILD_TAGS),debug -gcflags=all="-N -l"  # see delve docs

binary-dbg:
	$(GO_DBG_BUILD) -o $(GOBIN)/ ./binary/

server-dbg:
	$(GO_DBG_BUILD) -o $(GOBIN) ./server/

run:
	$(GOBIN)/server

clean:
	rm -rf build binary/binary server/server
