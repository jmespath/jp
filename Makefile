JP_VERSION=""

help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  test     to run all the tests."

test:
	test/vendor/bats/libexec/bats test/cases

# This will create/tag a new release locally, but not push anything.
# The workflow for a new release is:
#
# $ make new-release JP_VERSION=1.0.0
# < you'll get prompted for a few things >
# $ git push origin master --tags
# Go to github and create a release
#
# Right now, the last step isn't automated.  You still
# have to manually upload the release assets from build/.
new-release:
	scripts/bump-version $(JP_VERSION)
	git add jp.go && git commit -m "Bumping version to $(JP_VERSION)"
	git tag -s -m "Tagging $(JP_VERSION) release"
	scripts/build-all-platforms.sh
	scripts/sign-all.sh

.PHONY: help test
