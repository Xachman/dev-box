#!/bin/bash

run_unison="unison --auto /mnt/droplet /unison"

if [ -n "$TARGET" ]; then
    sshfs -o allow_other,default_permissions $TARGET /mnt/droplet
    $run_unison
    while sleep 1; do $run_unison; done
else
    echo "No TARGET variable"
    exit 1
fi