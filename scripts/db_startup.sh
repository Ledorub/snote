#!/bin/bash

source db_secret.sh

NOT_FOUND_E=66
NO_PERMISSION_E=77

main () {
  declare -A credentials

  if [[ $# -eq 3 ]]; then
    if ! check_file_exists $1; then
      fatal_error $NOT_FOUND_E "File not found: $1"
    fi
    if ! check_has_read_permission $1; then
      fatal_error $NO_PERMISSION_E "Permission error: $1"
    fi

    secret_path=$(dirname $1)
    read_credentials_from_file $1 credentials secret_path
    shift
  fi

  if ! check_file_exists $1; then
    fatal_error $NOT_FOUND_E "File not found: $1"
  fi
  if ! check_has_read_permission $1; then
    fatal_error $NO_PERMISSION_E "Permission error: $1"
  fi

  credentials_string=""
  credentials_to_string credentials credentials_string

  start_db "$@"
}

fatal_error() {
  echo $2
  exit $1
}

start_db() {
  exec env $credentials_string "$@"
}

main "$@"
