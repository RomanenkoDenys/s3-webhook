#!/bin/bash
# Require aws cli
# if file added - copy it to backup bucket
# if file removed - remove it from backup bucket
# Variables
ENDPOINT="https://hb.bizmrg.com"
AWSCLI=`which aws`" --endpoint-url ${ENDPOINT} s3"
BACKUP_BUCKET="myfiles-backup"

SOURCE_BUCKET="${1}"
SOURCE_FILE="${2}"
ACTION="${3}"

SOURCE="s3://${SOURCE_BUCKET}/${SOURCE_FILE}"
TARGET="s3://${BACKUP_BUCKET}/${SOURCE_FILE}"

case ${ACTION} in
    "copy")
    COMMAND="cp"
    ;;

    "delete")
    COMMAND="rm"
    SOURCE=""
    ;;

    *)
    echo "Usage: ${0} sourcebucket sourcefile copy/delete"
    exit 1
    ;;
esac

# Run aws cli command

${AWSCLI} ${COMMAND} ${SOURCE} ${TARGET}
