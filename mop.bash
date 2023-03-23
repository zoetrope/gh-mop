#!/bin/bash

# Fetch open issues from the repository.
function mop-issues() {
  echo "issues"
}

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

  export MOP_REPO="mop"
  export MOP_ISSUE=$1
  export MOP_STEP=0

  script -q -f -a /tmp/mop/${MOP_ISSUE}/typescript
}

if [ -n "$MOP_ISSUE" ]; then
  export PS1="[\${MOP_REPO} OP#\${MOP_ISSUE}:STEP\${MOP_STEP}]$ "
  # Move to the next step in the operation.
  function next() {
    MOP_STEP=$(($MOP_STEP + 1))
  }

  # Move to the previous step in the operation.
  function prev() {
    MOP_STEP=$(($MOP_STEP - 1))
  }

  # Insert a command into the current line.
  function insert() {
    local command=$(go run main.go next $MOP_ISSUE $MOP_STEP)
    READLINE_LINE="$command"
    let READLINE_POINT+=${#command}
  }

  # Show the current operation in Markdown.
  function show() {
    echo "show"
  }

  # Upload the results of executed commands.
  # The result is uploaded to the issue's comment.
  function upload() {
    echo "upload"
    cat /tmp/mop.log | ansi2txt | col -b
  }

  # List the commands in the current operation.
  function list() {
    echo "list"
  }

  # Search for a command from the snippet issue.
  function snippet() {
    echo "snippet"
  }

  # Show the help.
  function help() {
    echo "commands:"
    echo "  next: Move to the next step."
    echo "  prev: Move to the previous step."
    echo "  show: Show the current operation in Markdown."
    echo "  upload: Upload the results of executed commands."
    echo "shortcuts:"
    echo "  C-t: Insert a command into the current line."
    echo "  C-j: List the commands in the current operation."
    echo "  C-o: Search for a command from the snippet issue."
    echo "environment variables:"
    echo "  MOP_REPO: The repository name."
    echo "  MOP_ISSUE: The issue number for the operation."
    echo "  MOP_STEP: The step number in the operation."
  }

  bind -x '"\C-t":"insert"'
  bind -x '"\C-j":"list"'
  bind -x '"\C-o":"snippet"'
fi
