#!/bin/bash

set -x
source db_secret.sh

NOT_FOUND_E=66
NO_PERMISSION_E=77
MIGRATION_FAILED_E=125

main() {
  db_type="postgresql"
  host=$1
  port=$2
  migration_type=$4
  migration_dir=$5

  declare -A credentials

  echo -n "Checking the main secret file... "
  if ! check_file_exists "$3"; then
    fatal_error $NOT_FOUND_E "File not found: $3"
  elif ! check_has_read_permission $3; then
    fatal_error $NO_PERMISSION_E "Permission error: $3"
  fi
  echo "Success"

  secret_path=$(dirname "$3")
  echo -n "Reading credentials from the main secret file... "
  read_credentials_from_file "$3" credentials secret_path
  echo "Success"
  shift

  echo -n "Parsing credentials from nested files... "
  err=""
  replace_credential_files_with_content credentials err
  if [[ $? != 0 ]]; then
    fatal_error $? "$err"
  fi
  echo "Success"

  user="${credentials[POSTGRES_USER]}"
  pwd="${credentials[POSTGRES_PASSWORD]}"
  db="${credentials[POSTGRES_DB]}"
  echo "${#user}"

  echo -n "Checking migration dir... "
  if ! check_dir_exists "$migration_dir"; then
    fatal_error $NOT_FOUND_E "Migration dir not found: $migration_dir"
  elif ! check_dir_readable "$migration_dir"; then
    fatal_error $NOT_FOUND_E "Permission error: $migration_dir"
  fi
  echo "Success"

  dsn=""
  build_dsn dsn "$db_type" "$user" "$pwd" "$host" "$port" "$db"

  echo -n "Starting migration... "
  migrate_db "$dsn" "$migration_type" "$migration_dir"
  if [[ $? != 0 ]]; then
    fatal_error $MIGRATION_FAILED_E "Fail"
  fi
  echo "Success"
}

check_dir_exists() {
  [[ -d $1 ]]
}

check_dir_readable() {
  [[ -r $1 ]] && [[ -x $1 ]]
}

build_dsn() {
  local -n var=$1
  local url="$2://"

  encoded=""
  url_encode "encoded" "$3"
  url+="$encoded:"

  encoded=""
  url_encode "encoded" "$4"
  url+="$encoded@"

  encoded=""
  url_encode "encoded" "$5"
  url+="$encoded:"

  encoded=""
  url_encode "encoded" "$6"
  url+="$encoded/"

  encoded=""
  url_encode "encoded" "$7"
  url+="$encoded"

  url+="?sslmode=disable"

  var=$url
}

url_encode() {
  local -n var=$1
  var="$(jq -nsRr --arg s "$2" '$s|@uri')"
}

migrate_db() {
  migrate -path "$3" -database "$1" -verbose "$2"
}

fatal_error() {
  echo "$2"
  exit "$1"
}

main "$@"
