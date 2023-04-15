#!/bin/bash

MOP_CONFIG="$HOME/.mop.json"
GH_MOP="gh mop --config=${MOP_CONFIG}"

# For development
# MOP_CONFIG=".mop.json"
# GH_MOP="go run main.go --config=${MOP_CONFIG}"

# Start an operation.
# $1: The issue number for the operation.
function mop-start() {
  if [ -z $1 ]; then
    echo "Error: The issue number is required. Usage: \"mop-start <issue number>\"."
    return 1
  fi

  if [ -n "$MOP_ISSUE" ]; then
    echo "Error: The operation is already started. Issue: $MOP_ISSUE."
    return 1
  fi

  if [ ! -f $MOP_CONFIG ]; then
    echo "Error: The config file \"$MOP_CONFIG\" is not found."
    return 1
  fi

  if ! command -v gh &> /dev/null
  then
      echo "gh could not be found. Please install it with `sudo apt install gh`"
      return 1
  fi
  if ! gh auth status &> /dev/null
  then
      echo "You are not logged in to github. Please run `gh auth login`"
      return 1
  fi
  if ! command -v jq &> /dev/null
  then
      echo "jq could not be found. Please install it with `sudo apt install jq`"
      return 1
  fi
  if ! command -v script &> /dev/null
  then
      echo "script could not be found. Please install it with `sudo apt install script`"
      return 1
  fi

  export MOP_REPO=$(cat $MOP_CONFIG | jq -r .repository)
  if [ -z "$MOP_REPO" ]; then
    echo "Error: The repository name is not set in the config file."
    mop-clear-env
    return 1
  fi
  export MOP_DATADIR=$(cat $MOP_CONFIG | jq -r .datadir)
  if [ -z "$MOP_DATADIR" ]; then
    echo "Error: The data directory is not set in the config file."
    mop-clear-env
    return 1
  fi
  export MOP_ISSUE=$1
  export MOP_STEP=0
  export MOP_OFFSET=0

  mkdir -p ${MOP_DATADIR}/${MOP_REPO}/${MOP_ISSUE} || return $?

  $GH_MOP operation $MOP_ISSUE
  status=$?
  if [ $status -ne 0 ]; then
    mop-clear-env
    return $status
  fi
  export MOP_COMMAND_COUNT=$(cat ${MOP_DATADIR}/${MOP_REPO}/${MOP_ISSUE}/operation.json | jq '.commands | length')

  local record=${MOP_DATADIR}/${MOP_REPO}/${MOP_ISSUE}/typescript
  if [ -f $record ]; then
    mv $record ${record}_$(date +%Y%m%d%H%M%S)
  fi
  script -q -f -a $record
  status=$?
  if [ $status -ne 0 ]; then
    mop-clear-env
    return $status
  fi
  mop-clear-env
  return 0
}

# Clear environment variables
function mop-clear-env() {
  unset MOP_REPO
  unset MOP_DATADIR
  unset MOP_ISSUE
  unset MOP_STEP
  unset MOP_OFFSET
  unset MOP_COMMAND_COUNT
}

if [ -n "$MOP_ISSUE" ]; then
  export PS1="[\${MOP_REPO}#\${MOP_ISSUE}:Step\${MOP_STEP}]$ "
  trap 'finish' EXIT

  echo "The operation is started. Issue: $MOP_ISSUE."
  echo "If you want to read the help, type \"help\"."

  # Move to the next step in the operation.
  function next() {
    if [ $MOP_STEP -eq $(($MOP_COMMAND_COUNT - 1)) ]; then
      echo "The last step is reached."
      return 1
    fi
    MOP_STEP=$(($MOP_STEP + 1))
  }

  # Move to the previous step in the operation.
  function prev() {
    if [ $MOP_STEP -eq 0 ]; then
      echo "The first step is reached."
      return 1
    fi
    MOP_STEP=$(($MOP_STEP - 1))
  }

  # Insert the next command into the current line.
  function insert() {
    local command=$($GH_MOP command $MOP_ISSUE $MOP_STEP)
    READLINE_LINE="$command"
    let READLINE_POINT+=${#command}
  }

  # Upload the results of executed commands.
  # The result is uploaded to the issue's comment.
  function upload() {
    local result
    result=$($GH_MOP upload --offset $MOP_OFFSET $MOP_ISSUE ${MOP_DATADIR}/${MOP_REPO}/${MOP_ISSUE}/typescript)
    status=$?
    if [ $status -ne 0 ]; then
      return $status
    fi
    MOP_OFFSET=$(echo $result | jq -r .offset)
    url=$(echo $result | jq -r .url)
    echo "The result is uploaded to $url"
  }

  # List the commands in the current operation.
  function list() {
    local command=$($GH_MOP list $MOP_ISSUE)
    READLINE_LINE="$command"
    let READLINE_POINT+=${#command}
  }

  # Search for a command from the snippet issue.
  function utilities() {
    local command=$($GH_MOP utilities)
    READLINE_LINE="$command"
    let READLINE_POINT+=${#command}
  }

  # Show the help.
  function help() {
    echo "commands:"
    echo "  next: Move to the next step."
    echo "  prev: Move to the previous step."
    echo "  upload: Upload the results of executed commands."
    echo "shortcuts:"
    echo "  C-t: Insert the next command into the current line."
    echo "  C-j: List the commands in the current operation."
    echo "  C-o: Search for a command from the utilities issue."
    echo "environment variables:"
    echo "  MOP_REPO: The repository name."
    echo "  MOP_ISSUE: The issue number for the operation."
    echo "  MOP_STEP: The step number in the operation."
    echo "  MOP_OFFSET: The offset of the result in the issue comment."
    echo "  MOP_COMMAND_COUNT: The number of commands in the operation."
  }

  function finish() {
    mop-clear-env

    unset -f next
    unset -f prev
    unset -f insert
    unset -f upload
    unset -f list
    unset -f utilities
    unset -f help
    unset -f finish

    bind -r "\C-t"
    bind -r "\C-j"
    bind -r "\C-o"
  }

  bind -x '"\C-t":"insert"'
  bind -x '"\C-j":"list"'
  bind -x '"\C-o":"utilities"'
fi
