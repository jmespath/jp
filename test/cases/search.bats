@test "Has valid help output" {
  run ./jp --help
  [ "$status" -eq 0 ]
  echo $output | grep "\-\-filename"
}

@test "Can display version" {
  run ./jp --version
  [ "$status" -eq 0 ]
}

@test "Can search basic expression" {
  output=$(echo '{"foo": "bar"}' | ./jp foo)
  [ "$output" == "\"bar\"" ]
}

@test "Ignores subsequent data" {
  output=$(echo '{"foo": "bar"}blah' | ./jp foo)
  [ "$output" == "\"bar\"" ]
}

@test "Processes subsequent data in stream mode" {
  output=$(echo '{"foo": "bar"}{"foo": "x"}' | ./jpp foo)
  echo "$output"
  [ "$output" == $'\"bar\"\n\"x\"' ]
}

@test "Test recursive accumulate mode for nested lists" {
  output=$(echo '{"foo": ["a"]}{"foo": ["x"]}' | ./jpp -a -c @)
  echo "$output"
  [ "$output" == '{"foo":["a","x"]}' ]
}

@test "Test recursive accumulate mode for nested lists" {
  output=$(echo '["a"]["x"]' | ./jpp -a -c @)
  echo "$output"
  [ "$output" == '["a","x"]' ]
}

@test "Test recursive accumulate mode for nested lists" {
  output=$(echo '["a"]["x"]' | ./jpp -a -c -u @)
  echo "$output"
  [ "$output" == $'"a"\n"x"' ]
}

@test "Test that recursive accumulate mode coalesces duplicates from different nested lists" {
  output=$(echo '{"foo": ["a", "a"]}{"foo": ["a"]}' | ./jpp -a -c @)
  echo "$output"
  [ "$output" == '{"foo":["a","a"]}' ]
}

@test "Test that recursive accumulate mode preserves duplicates from the same nested list" {
  output=$(echo '{"foo": ["a", "a"]}{"foo": ["a", "b"]}' | ./jpp -a -c @)
  echo "$output"
  [ "$output" == '{"foo":["a","a","b"]}' ]
}

@test "Test raw string input" {
  output=$(echo 'hello world' | ./jpp -R -a -c @)
  echo "$output"
  [ "$output" == '"hello world"' ]
}

@test "Test raw string input" {
  output=$(echo 'hello world' | ./jpp -R -r @)
  echo "$output"
  [ "$output" == 'hello world' ]
}

@test "Test raw string input" {
  output=$(echo 'hello world' | ./jpp -R --unquoted @)
  echo "$output"
  [ "$output" == 'hello world' ]
}

@test "Test multi-line raw string input" {
  output=$(printf -- '%s\n' 'line '{1..3} | ./jpp -R -c @)
  echo "$output"
  [ "$output" == $'"line 1"\n"line 2"\n"line 3"\n' ]
}

@test "Test multi-line raw string input slurp" {
  output=$(printf -- '%s\n' 'line '{1..3} | ./jpp -R -s -c @)
  echo "$output"
  [ "$output" == '"line 1\nline 2\nline 3\n"' ]
}

@test "Can search subexpr expression" {
  output=$(echo '{"foo": {"bar": "baz"}}' | ./jp foo.bar)
  [ "$output" == "\"baz\"" ]
}

@test "Can read input from file" {
  echo '{"foo": "bar"}' > "$BATS_TMPDIR/input.json"
  run ./jp -f "$BATS_TMPDIR/input.json" foo
  [ "$output" == "\"bar\"" ]
}

@test "Can print result unquoted" {
  output=$(echo '{"foo": "bar"}' | ./jp -u foo)
  [ "$output" == "bar" ]
}

@test "Bad JMESPath expression has non zero rc" {
  echo '{"foo": "bar"}' > "$BATS_TMPDIR/input.json"
  run ./jp -f "$BATS_TMPDIR/input.json" "bax[expre]ssion"
  [ "$status" -eq 1 ]
}

@test "Large numbers are not printed with scientific notation" {
  skip
  echo '{"foo": 47268765}' > "$BATS_TMPDIR/input.json"
  run ./jp -f "$BATS_TMPDIR/input.json" "foo"
  [ "$status" -eq 0 ]
  [ "$output" == "47268765" ]
}

@test "Can accept expression from file" {
  echo 'foo.bar' > "$BATS_TMPDIR/expr"
  echo '{"foo": {"bar": "baz"}}' > "$BATS_TMPDIR/input.json"
  run ./jp -u -f "$BATS_TMPDIR/input.json" -e "$BATS_TMPDIR/expr"
  [ "$output" == "baz" ]
}

@test "Can pretty print expr AST" {
  run ./jp --ast "foo"
  expected='
ASTField {
  value: "foo"
}'
  echo "$output"
  echo "$expected"
  [ "$output" == "$expected" ]
}

@test "Can sort int array" {
  echo '[2,1,3,5,4]' > "$BATS_TMPDIR/input.json"
  echo "sort(@) | map(&to_string(@), @) | join('', @)" > "$BATS_TMPDIR/expr"
  run ./jp -u -f "$BATS_TMPDIR/input.json" -e "$BATS_TMPDIR/expr"
  [ "$status" -eq 0 ]
  [ "$output" == "12345" ]
}
