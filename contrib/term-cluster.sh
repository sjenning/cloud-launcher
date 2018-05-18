#!/usr/bin/env bash

## Script to mark nodes for termination on AWS
## usage: term-cluster.sh [filter on name] [region (default: us-east-1)]
##
## Dependencies: jq and awless
## 

set -eou pipefail

die() {
  echo $@
  exit 1
}

FILTER=${1:-}
if [[ $FILTER == "" ]]; then
  die "invalid filter"
fi

REGION=${2:-us-east-1}
declare -A instances=()

# Query instances
awsdata=$(awless list instances --filter name="$FILTER" --aws-region=$REGION --format json --no-sync)
filtered=$(echo $awsdata | jq -r '.[] | "\(.ID) \(.Name) "')
while read id name; do
  echo "Instance $id - $name"
  instances[$id]=$name
done <<< "$filtered"

# Prompt for deletion
read -p "Terminate these instances? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  for key in "${!instances[@]}"; do 
    name=${instances[$key]}
    awless create tag resource=$key key=Name value="$name-terminate" --aws-region=$REGION -f --no-sync
  done
fi
