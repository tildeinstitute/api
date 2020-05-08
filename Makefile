PREFIX?=/var/www
_INSTDIR=$(PREFIX)
BINDIR?=$(_INSTDIR)/api

api: main.go go.mod
	@printf "\n%s\n\n" "Building tilde.institute API"
	go build -o $@
	@printf "\n%s\n\n" "...Done!"

.PHONY: clean
clean:
	@printf "\n%s\n\n" "Cleaning build and module caches..."
	go clean
	go clean -cache -modcache
	@printf "\n%s\n\n" "...Done!"

.PHONY: update
update:
	@printf "\n%s\n\n" "Updating from upstream repository..."
	git pull --rebase origin master
	@printf "\n%s\n\n" "...Done!"

.PHONY: install
install:
	@printf "\n%s\n\n%s\n" "Installing ..." "Creating Directories ..."
	mkdir -p $(BINDIR)/web
	@printf "\n%s\n" "Copying files..."
	install -m755 api $(BINDIR)
	install -m644 web/* $(BINDIR)
	@printf "\n%s\n" "Setting ownership..."
	chown -R www:www $(BINDIR)
	@printf "\n%s\n\n" "...Done!"

.PHONY: uninstall
uninstall:
	@printf "\n%s\n\n" "Removing files..."
	rm -rf $(BINDIR)
	@printf "\n%s\n\n" "...Done!"
