#!/bin/bash

NOT_FOUND_E=66
NO_PERMISSION_E=77

check_file_exists() {
  [[ -f $1 ]]
}

check_has_read_permission() {
  [[ -r $1 ]]
}

read_credentials_from_file() {
  local -n _credentials=$2

  while IFS== read -r key value; do
    if [[ -n "${!key}" ]]; then
      continue
    fi

    if [[ "$key" == *_FILE ]]; then
      value="${!3}/$value"
    fi

    _credentials[$key]=$value
  done < "$1"
}

credentials_to_string() {
  local -n _credentials=$1
  local -n _s=$2

  for key in "${!_credentials[@]}"; do
    _s+="$key=${_credentials[$key]} "
  done

  _s=$(echo $_s | xargs)
}

replace_credential_files_with_content() {
  local -n _credentials=$1
  local -n _err=$2

  for key in "${!_credentials[@]}"; do
    local value="${_credentials[$key]}"

    if [[ "$key" != *_FILE ]]; then
      continue
    fi

    if ! check_file_exists $value; then
      _err="File not found: $value"
      return $NOT_FOUND_E
    fi
    if ! check_has_read_permission $value; then
      _err="Permission error: $value"
      return $NO_PERMISSION_E
    fi

    unset _credentials[$key]
    key="${key%"_FILE"}"
    _credentials[$key]="$(< "$value")"
  done
}
