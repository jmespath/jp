jp
==

The ``jp`` command is a command line interface to
[JMESPath](http://jmespath.org), an expression
language for manipulating JSON.


Installing
==========

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


Usage
=====

The most basic usage of ``jp`` is to accept input JSON data through
stdin, apply the JMESPath expression you've provided as an argument,
and print the resulting JSON data stdout.
