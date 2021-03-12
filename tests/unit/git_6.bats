#!/usr/bin/env bats

load test_helper

setup() {
  global_setup
  touch /home/dokku/.ssh/known_hosts
  chown dokku:dokku /home/dokku/.ssh/known_hosts
  mkdir -p "$DOKKU_LIB_ROOT/data/git/$TEST_APP"
  chown dokku:dokku "$DOKKU_LIB_ROOT/data/git/$TEST_APP"
}

teardown() {
  rm -f /home/dokku/.ssh/id_rsa.pub || true
  global_teardown
}

@test "(git) git:unlock [success]" {
  run /bin/bash -c "dokku git:unlock $TEST_APP --force"
  echo "output: $output"
  echo "status: $status"
  assert_success
}

@test "(git) git:unlock [missing arg]" {
  run /bin/bash -c "dokku git:unlock"
  echo "output: $output"
  echo "status: $status"
  assert_failure
}
