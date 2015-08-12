help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  test     to run all the tests."

test:
	test/vendor/bats/libexec/bats test/cases

.PHONY: help test
