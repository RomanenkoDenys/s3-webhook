#!/bin/bash
# Require aws cli
# if file added - copy it to backup bucket
# if file removed - remove it from backup bucket
# Variables
ENDPOINT_MCS="https://hb.bizmrg.com"
AWSCLI_MCS=`which aws`" --endpoint-url ${ENDPOINT_MCS} --profile mcs s3"
AWSCLI_AWS=`which aws`" --profile aws s3"
BACKUP_BUCKET="myfiles-backup"

SOURCE_BUCKET="${1}"
SOURCE_FILE="${2}"
ACTION="${3}"

SOURCE="s3://${SOURCE_BUCKET}/${SOURCE_FILE}"
TARGET="s3://${BACKUP_BUCKET}/${SOURCE_FILE}"
TEMP="/tmp/${SOURCE_BUCKET}/${SOURCE_FILE}"

case ${ACTION} in
    "copy")
    ${AWSCLI_MCS} cp "${SOURCE}" "${TEMP}"
    ${AWSCLI_AWS} cp "${TEMP}" "${TARGET}"
    rm ${TEMP}
    ;;

    "delete")
    ${AWSCLI_AWS} rm ${TARGET}
    ;;

    *)
    echo "Usage: ${0} sourcebucket sourcefile copy/delete"
    exit 1
    ;;
esac
