# Contributing

We work hard to provide a high-quality and useful CLI, and we greatly value
feedback and contributions from our community. Whether it's a new feature,
correction, or additional documentation, we welcome your pull requests. Please
submit any [issues](https://github.com/jmespath/jp/issues)
or [pull requests](https://github.com/jmespath/jp/pulls)
through GitHub.

This document contains guidelines for contributing code and filing issues.

## Contributing Code

This list below are guidelines to use when submitting pull requests.
These are the same set of guidelines that the core contributors use
when submitting changes, and we ask the same of all community
contributions as well:

* We maintain a high percentage of code coverage in our tests.  As a general
  rule of thumb, code changes should not lower the overall code coverage
  percentage for the project. In practice, this means that every bug fix and
  feature addition should include tests.
* All PRs must run cleanly through `make test`.
  in more detail in the sections below.


## Feature Development

This CLI is designed to be a reference CLI implementation of JMESPath,
and implements the official JMESPath specification.  As such, we do not
accept feature requests that are not part of the official JMESPath
specification.  If you would like to propose changes to the JMESPath
language itself, which includes new syntax, functions, operators, etc.,
you can create a JMESPath enhancement proposal in the
[jmespath.jep repository](https://github.com/jmespath/jmespath.jep).
Once the proposal is accepted, we can then implement the changes in this
CLI.

## Running Tests

To run the tests, you can run `make test`.
