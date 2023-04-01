#!/bin/bash

GH_MOP="go run main.go --config=config.json"
# GH_MOP="gh mop"

# Start an operation.
# $1: The issue number for the operation.
function mop-start() {
  if [ -z $1 ]; then
    echo "Error: The issue number is required."
    return 1
  fi

  if [ -n "$MOP_ISSUE" ]; then
    echo "Error: The operation is already started."
    return 1
  fi

  export MOP_REPO=$(cat config.json | jq -r .repository)
  export MOP_DATADIR=$(cat config.json | jq -r .datadir)
  export MOP_ISSUE=$1
  export MOP_STEP=0

  $GH_MOP operation $MOP_ISSUE
  script -q -f -a ${MOP_DATADIR}/${MOP_REPO}/${MOP_ISSUE}/typescript.txt
}

if [ -n "$MOP_ISSUE" ]; then
  export PS1="[\${MOP_REPO}#\${MOP_ISSUE}:Step\${MOP_STEP}]$ "
  # Move to the next step in the operation.
  function next() {
    MOP_STEP=$(($MOP_STEP + 1))
  }

  # Move to the previous step in the operation.
  function prev() {
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
    $GH_MOP upload $MOP_ISSUE ${MOP_DATADIR}/${MOP_REPO}/${MOP_ISSUE}/typescript.txt
  }

  # List the commands in the current operation.
  function list() {
    echo "list"
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
  }

  bind -x '"\C-t":"insert"'
  bind -x '"\C-j":"list"'
  bind -x '"\C-o":"utilities"'
fi
