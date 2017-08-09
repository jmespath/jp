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
  output=$(echo '{"foo": "bar"}{"foo": "x"}' | ./jp -s foo)
  echo "$output"
  [ "$output" == $'\"bar\"\n\"x\"' ]
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
