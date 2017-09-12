jp
==

The ``jp`` command is a command line interface to
[JMESPath](http://jmespath.org), an expression
language for manipulating JSON:

```
$ echo '{"foo": {"bar": ["a", "b", "c"]}}' | jp foo.bar[1]
"b"
```


# Installing

If you're a Mac user, you can install via homebrew from the JMESPath
Homebrew tap:

```
brew tap jmespath/jmespath
brew install jp
```

You can download prebuilt binaries if you prefer.
Check the [Release page](https://github.com/jmespath/jp/releases)
 to download the latest ``jp`` executable.  There are binaries
available for Windows, Linux, Mac, FreeBSD.

For example, to install version 0.1.2 on a 64 bit Linux environment use:

```
sudo wget https://github.com/jmespath/jp/releases/download/0.1.2/jp-linux-amd64 -O /usr/local/bin/jp \
&& sudo chmod +x /usr/local/bin/jp
```

## Building from Source

If you have a Go environment installed you can also run:
``go get -u github.com/jmespath/jp`` to get the latest version
of jmespath.  If you have the repo checked out locally you can also
just ``go build`` the project:

```
git clone git://github.com/jmespath/jp
cd jp
go build
./jp --help
```

And finally, if you have a go environment setup, but don't
have a workspace/GOPATH configured, you can just run
``scripts/build-self-contained`` and it will build the
``jp`` executable for you:

```
git clone git://github.com/jmespath/jp
cd jp
scripts/build-self-contained
./jp --help
```

## Cross platform binaries

If you have a go 1.5 environment setup, you can build cross platform binaries
by running ``scripts/build-all-platforms``.
This will put executables in the ``build/`` directory and each
executable will be named ``jp-<platform>``, e.g
``jp-darwin-amd64``,  ``jp-linux-386``, etc.


# Usage

The most basic usage of ``jp`` is to accept input JSON data through stdin,
apply the JMESPath expression you've provided as an argument to ``jp``, and
print the resulting JSON data to stdout.

```
$ echo '{"key": "value"}' | jp key
"value"

$ echo '{"foo": {"bar": ["a", "b", "c"]}}' | jp foo.bar[1]
"b"
```

Note the argument after ``jp``.  This is a JMESPath expression.
If you have no idea what that is, there's a
[JMESPath Tutorial](http://jmespath.org/tutorial.html) that
will take you through the JMESPath language.

## Input From a File

In addition to this basic usage, there's also other ways to use
``jp``.  First, instead of reading from stdin, you can provide
a JSON file as input using the ``-f/--filename`` option:


```
$ echo '{"foo": {"bar": "baz"}}' > /tmp/input.json
$ jp -f /tmp/input.json foo.bar
"baz"
```

## Unquoted Output

[[Notice]] the output of the above command is ``"baz"``, that is,
a double quote ``"``, followed by baz, followed by another
a final double quote.  This can be problematic if you're
trying to use this with other commands that just want
the string and *not* the quoted string.  For example:


```
$ curl -s https://api.github.com/repos/golang/go/events | jp [0].actor.url
"https://api.github.com/users/robpike"
```

Now let's suppose we want to then curl the above URL.  Our first
attempt might look something like this:

```
$ curl $(curl -s https://api.github.com/repos/golang/go/events | ./jp [0].actor.url)
```

And it would fail with:

```
curl: (1) Protocol "https not supported or disabled in libcurl
```

To fix this, we can use the ``-u/--unquoted`` option to specify that
any result that is a string will be printed without quotes.  Note
that the result is not surrounded by double quotes:

```
$ curl -s https://api.github.com/repos/golang/go/events | jp --unquoted [0].actor.url
https://api.github.com/users/robpike
```

If this is a common enough occurance for you, you can set the ``JP_UNQUOTED`` environment
variable to make this the default behavior:

```
$ export JP_UNQUOTED=true
$ curl -s https://api.github.com/repos/golang/go/events | jp --unquoted [0].actor.url
https://api.github.com/users/robpike
```


Also keep in mind that this behavior
only applies if the result of evaluating the JMESPath expression is a string:

```
$ echo '{"foo": ["bar", "baz"]}' | jp -u foo[0]
bar
# But -u does nothing here because the result is an array, not a string:
$ echo '{"foo": ["bar", "baz"]}' | jp -u foo
[
  "bar",
  "baz"
]
```

You can also use the ``-u/--unquoted`` option along with the
[join](http://jmespath.org/specification.html#join) function to create a list
of strings that can be piped into other POSIX text tools. For example:

```
$ echo '{"foo": {"bar": ["a", "b", "c"]}}' | jp foo.bar
[
  "a",
  "b",
  "c"
]
```

Suppose we want to iterate over the 3 values in the list and
run some bash code for each value.  We can do this by running:

```
$ for name in $(echo '{"foo": {"bar": ["a", "b", "c"]}}' | \
  jp -u 'join(`"\n"`, foo.bar)');
  do
      echo "Processing: $name";
  done
Processing: a
Processing: b
Processing: c
```


## Examples

If you're new to the JMESPath language, or just want to see what the language is
capable of, you can check out the [JMESPath tutorial](http://jmespath.org/tutorial.html)
as well as the [JMESPath examples](http://jmespath.org/examples.html), which contains
a curated set of JMESPath examples.  But for now, here's a real world example.
Let's say you wanted to see what the latest activity was with regard to the issue
tracker for one of your github issues.  Here's a simple way to do this:

```
$ curl -s https://api.github.com/repos/golang/go/events | jp \
"[?type=='IssuesEvent'].payload.\
{Title: issue.title, URL: issue.url, User: issue.user.login, Event: action}"

[
  {
    "Event": "opened",
    "Title": "release: cherry pick changes for 1.5 to release branch",
    "URL": "https://api.github.com/repos/golang/go/issues/12093",
    "User": "adg"
  },
  {
    "Event": "closed",
    "Title": "fmt: x format verb for []byte fails in a recursive call to Fscanf from a scanln call in go1.5rc1",
    "URL": "https://api.github.com/repos/golang/go/issues/12090",
    "User": "hubslave"
  },
  {
    "Event": "closed",
    "Title": "doc: release notes recommend wrong version of NaCl",
    "URL": "https://api.github.com/repos/golang/go/issues/12062",
    "User": "davecheney"
  },
  {
    "Event": "opened",
    "Title": "cmd/godoc: show internal packages when explicitly requested",
    "URL": "https://api.github.com/repos/golang/go/issues/12092",
    "User": "jacobsa"
  }
]
```

Try it for your own repo, instead of ``/golang/go``, replace it with your own
``/owner/repo`` value.  In words, this expression says:

* For each element in the top level list, select only the elements where the
``type`` key is equal to the string ``IssueEvent``
* For each of those filtered elements select the ``payload`` hash.
* Each each ``payload`` hash, we're going to create our own hash that has
4 keys: ``Title``, ``URL``, ``User``, ``Event``.  The value for each of key
is the result of evaluating these expressions in their respective order:
``issue.title``, ``issue.url``, ``issue.user.login``, ``action``.

Ensure that if your expression has spaces you surround the expression
in quotes, as shown in the example above.

## Testing

The parsing and evaluation of JMESPath expression is done in the
go-jmespath library, which is a dependency of this project.  ``go-jmespath``
has extensive testing to ensure it is parsing and evaluating JMESPath
expressions correctly.

To ensure that there are no regressions between `go-imespath` and `jp`,
the entire suite of [JMESPath compliance tests](https://github.com/jmespath/jmespath.test)
are run against the `jp` executable.

This repo also include CLI specific test that verify the command line
options and output work as intended.

You can run all of these tests for `jp` by running `make test`:

```
$ make test
# CLI specific test cases.
test/vendor/bats/libexec/bats test/cases
 ✓ Has valid help output
 ✓ Can display version
 ✓ Can search basic expression
 ✓ Can search subexpr expression
 ✓ Can read input from file
 ✓ Can print result unquoted
 ✓ Bad JMESPath expression has non zero rc
 - Large numbers are not printed with scientific notation (skipped)
 ✓ Can accept expression from file
 ✓ Can pretty print expr AST
 ✓ Can sort int array

X tests, 0 failures, 1 skipped
# JMESPath compliance tests, using the jp-compliance
# runner from github.com/jmespath/jmespath.test
test/jp-compliance -d test/compliance/ -e ./jp
............................................................
............................................................
............................................................
............................................................
............................................................
............................................................
............................................................
............................................................
OK
```
