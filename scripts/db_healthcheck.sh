#!/bin/bash

NOT_FOUND_E=66
NO_PERMISSION_E=77

source db_secret.sh

main () {
  local host="127.0.0.1"
  local port="5432"
  local user=$POSTGRES_USER
  local db=$POSTGRES_DB

  declare -A credentials

  if [[ $# -eq 1 ]]; then
    if ! check_file_exists $1; then
      fatal_error $NOT_FOUND_E "File not found: $1"
    fi
    if ! check_has_read_permission $1; then
      fatal_error $NO_PERMISSION_E "Permission error: $1"
    fi

    secret_path=$(dirname $1)
    read_credentials_from_file $1 credentials secret_path
    shift

    err=""
    replace_credential_files_with_content credentials err
    if [[ $? != 0 ]]; then
      fatal_error $? $err
    fi
  fi

  user="${credentials[POSTGRES_USER]}"
  db="${credentials[POSTGRES_DB]}"

  check_health $host $port $user $db
}

fatal_error() {
  echo $2
  exit $1
}

check_health() {
  pg_isready -h $1 -p $2 -U $3 -d $4
}

main "$@"
