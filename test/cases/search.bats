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
  echo '{"foo": 47268765}' > "$BATS_TMPDIR/input.json"
  run ./jp -f "$BATS_TMPDIR/input.json" "foo"
  [ "$status" -eq 0 ]
  [ "$output" == "47268765" ]
}
