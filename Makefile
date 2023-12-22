TOPDIR:=$(shell pwd)
BINDIR:=$(TOPDIR)/bin
SRVDIR:=$(TOPDIR)/cmd

TARGETS:=loginserver gameserver masterserver

# detect HOST OS
# however, we can produce other builds with overrided HOST_OS variable
ifeq ($(OS),Windows_NT)
	HOST_OS ?= windows
else
	HOST_OS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
endif

# add *.exe extension on Windows builds
ADD_EXT:=$(if $(findstring $(HOST_OS),windows),.exe)

.PHONY: pre all clean

all: pre
	$(foreach server,$(TARGETS), \
		GOOS=$(HOST_OS) go build -v $(DEV_FLAGS) -o $(BINDIR)/$(server)$(ADD_EXT) $(SRVDIR)/$(server);)

pre:
	mkdir -p $(BINDIR)

clean:
	rm -r $(BINDIR)

%:
	$(MAKE) -C $(TOPDIR) pre
	GOOS=$(HOST_OS) go build -v $(DEV_FLAGS) -o $(BINDIR)/$@$(ADD_EXT) $(SRVDIR)/$@
