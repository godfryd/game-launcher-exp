
$(GOPATH)/bin/game-launcher:
	/usr/bin/go build -v -x -o $(GOPATH)/bin/game-launcher github.com/godfryd/game-launcher-exp/engo

.PHONY: $(GOPATH)/bin/game-launcher

play: $(GOPATH)/bin/game-launcher
	$(GOPATH)/bin/game-launcher

clean:
	rm -rf $(GOPATH)/bin/game-launcher
