#!/bin/sh

# This client side check can easily be bypassed.
# git rebase and git commit --amend both bypass it.
# Passing env var LEFTHOOK=0 also bypasses.
# Ideally this client side check and a server side
# pre-receive hook would be used. I am not
# convinced it it worth the development effort.

# Running a script to check for invalid and
# improper prefixs during prefix analysis is
# sufficient.

DIRECTORY="gridlock"
VALID_PREFIXS='^(fix|add|update|remove|rewrite|refactor|optimize|plan'\
'|comment|style|test|tool|misc):|^w$|^wip$'
# git passes different args for different git hooks
# commit-msg hook only passes the commit message
# https://lefthook.dev/usage/features/git-args.html
# https://git-scm.com/docs/githooks
COMMIT_MSG=$1
FIRST_LINE=$(head -n1 "$COMMIT_MSG")
COMMIT_FILES=$(git diff --name-only --cached)

if [ ! -d "$DIRECTORY" ]; then
  echo "Did not find $DIRECTORY cannot validate git message prefix."
  echo "If the directory was renamed upate the commit-msg hook."
  exit 1
fi

VALIDATE_PREFIX="false"
for file_name in $COMMIT_FILES; do
  directory=$(echo "$file_name" | cut -d '/' -f1)
  if [[ "$directory" == "$DIRECTORY" ]]; then
    VALIDATE_PREFIX="true"
  fi
done

if [ "$VALIDATE_PREFIX" = "false" ]; then
  exit 0
fi

if ! [[ "$FIRST_LINE" =~ $VALID_PREFIXS ]]; then
  echo "invalid git commit message prefix."
  echo "valid prefixes: $VALID_PREFIXS"
  exit 1
fi
