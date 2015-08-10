jp
==

The ``jp`` command is a command line interface to
[JMESPath](http://jmespath.org), an expression
language for manipulating JSON.


# Installing

Check the Release page to download the latest ``jp`` executable.
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

Notice the output of the above command is ``"baz"``, that is,
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
variable to make this the default behavior.   Also keep in mind that this behavior
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
