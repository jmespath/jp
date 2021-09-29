JP_VERSION=""


all: build test

build:
	go build ./...

test:
	# CLI specific test cases.
	test/vendor/bats/libexec/bats test/cases
	# JMESPath compliance tests, using the jp-compliance
	# runner from github.com/jmespath/jmespath.test
	test/jp-compliance -d test/compliance/ -e ./jp

# This will create/tag a new release locally, but not push anything.
# The workflow for a new release is:
#
# 1. $ make new-release JP_VERSION=1.0.0
# < you'll get prompted for a few things >
#

# 2. $ git push origin master --tags
#
#
# 3. Go to github and create a release
#    Right now, the last step isn't automated.  You still
#    have to manually upload the release assets from build/.
new-release:
	scripts/bump-version $(JP_VERSION)
	git add jp.go && git commit -m "Bumping version to $(JP_VERSION)"
	git tag -s -m "Tagging $(JP_VERSION) release" $(JP_VERSION)
	scripts/build-all-platforms
	scripts/sign-all

.PHONY: help test build
