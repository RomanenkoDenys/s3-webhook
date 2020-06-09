#!/bin/bash
# Require aws cli
# if file added - copy it to backup bucket
# if file removed - remove it from backup bucket
# Variables
BACKUP_BUCKET="s3://myfiles-backup"

ACTION="${1}"
SOURCEFILE="${2}"

if [ -z ${SOURCEFILE} ]; then
    echo "Usage: ${0} action sourcefile"
    exit 1
fi


echo "Success"
