.POSIX:

PNAME = pgot

RTEMPLATE ?= ../repo-template

all: goUtil doc

doc: docMain docMan

clean: cleanGoUtil cleanCheck

cleanDoc: cleanDocMain cleanDocMan

install: installGoUtil installMan

uninstall: uninstallGoUtil uninstallMan

.DEFAULT_GOAL := all

.PHONY: all doc clean cleanDoc install uninstall

#---Helper Macros to Remove Files---

RMDIR_IF_EMPTY := sh -c '\
if test -d $$0 && ! ls -1qA $$0 | grep -q . ; then \
	rmdir $$0; \
fi'

RM ?= rm -f

#---Generate Golang Utility---

PREFIX ?= /usr/local
_INSTDIR = $(DESTDIR)$(PREFIX)
BINDIR ?= $(_INSTDIR)/bin
GO ?= go
GOFLAGS ?=

VERSION != git describe --first-parent 2> /dev/null

GOSRC != find . -name '*.go'

goUtil: $(PNAME)

$(PNAME): $(GOSRC)
	$(GO) build $(GOFLAGS) \
		-ldflags "-X main.Version=$(VERSION)" \
		-o $@

cleanGoUtil:
	$(RM) $(PNAME)

installGoUtil: $(PNAME)
	strip $(PNAME)
	mkdir -m755 -p $(BINDIR)
	install -m755 $(PNAME) $(BINDIR)/$(PNAME)

uninstallGoUtil:
	$(RM) $(BINDIR)/$(PNAME)
	$(RMDIR_IF_EMPTY) $(BINDIR)

.PHONY: goUtil cleanGoUtil installGoUtil uninstallGoUtil

#---Generate Man Page(s)---

.SUFFIXES:
.SUFFIXES: .1 .5 .1.scd .5.scd

PREFIX ?= /usr/local
_INSTDIR = $(DESTDIR)$(PREFIX)
MANDIR ?= $(_INSTDIR)/share/man

DOCMAN := doc/$(PNAME).1 doc/$(PNAME).5

.1.scd.1:
	scdoc < $< > $@

.5.scd.5:
	scdoc < $< > $@

docMan: $(DOCMAN)

cleanDocMan:
	$(RM) $(DOCMAN)

installMan: $(DOCMAN)
	mkdir -m755 -p $(MANDIR)/man1 $(MANDIR)/man5
	install -m644 doc/$(PNAME).1 $(MANDIR)/man1/$(PNAME).1
	install -m644 doc/$(PNAME).5 $(MANDIR)/man5/$(PNAME).5

uninstallMan:
	$(RM) $(MANDIR)/man1/$(PNAME).1
	$(RM) $(MANDIR)/man5/$(PNAME).5
	$(RMDIR_IF_EMPTY) $(MANDIR)/man1
	$(RMDIR_IF_EMPTY) $(MANDIR)/man5
	$(RMDIR_IF_EMPTY) $(MANDIR)

.PHONY: installMan uninstallMan

#---Test/Check Section---

TESTDIR = tests

check: $(PNAME)
	cd $(TESTDIR) && go test -v

cleanCheck:
	find $(TESTDIR) -name '*.result' -delete

.PHONY: check cleanCheck

#---Generate Main Documents---

DOCMAIN := README.md LICENSE

README.md: template/README.md.got
	pgot -i ":$(RTEMPLATE)" -o $@ $<

LICENSE: $(RTEMPLATE)/LICENSE.src/BSD-2-clause.got
	pgot -i ":$(RTEMPLATE)" -o $@ $<

docMain: $(DOCMAIN)

cleanDocMain:
	$(RM) $(DOCMAIN)

.PHONY: docMain, cleanDocMain

#---Generate Makefile---

Makefile: template/Makefile.got
	pgot -i ":$(RTEMPLATE)" -o $@ $<

mkFile: Makefile

regenMkFile:
	pgot -i ":$(RTEMPLATE)" -o Makefile template/Makefile.got

.PHONY: mkFile regenMkFile

#---Lint Helper Target---

lint:
	@find . -path ./.git -prune -or \
		-type f -and -not -name 'Makefile' \
		-exec grep -Hn '<no value>' '{}' ';'

#---TODO Helper Target---

todo:
	@find . -path ./.git -prune -or \
		-type f -and -not -name 'Makefile' \
		-exec grep -Hn TODO '{}' ';'

# vim:set noet tw=80:
