TOPDIR:=$(shell pwd)
BINDIR:=$(TOPDIR)/bin
SRVDIR:=$(TOPDIR)/cmd

TARGETS:=loginserver gameserver masterserver

ADD_EXT:=$(if $(findstring $(GOOS),windows),.exe)

.PHONY: pre all clean

all: pre
	$(foreach server,$(TARGETS), \
		go build -v -o $(BINDIR)/$(server)$(ADD_EXT) $(SRVDIR)/$(server);)

pre:
	mkdir -p $(BINDIR)

clean:
	rm -r $(BINDIR)

%:
	$(MAKE) -C $(TOPDIR) pre
	go build -v -o $(BINDIR)/$@$(ADD_EXT) $(SRVDIR)/$@