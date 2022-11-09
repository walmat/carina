#!/bin/sh

(cd ../web && yarn install)

FILE=../web/.env
if [ ! -f "$FILE" ]; then
    echo "SKIP_PREFLIGHT_CHECK=true\nBROWSER=none" > $FILE
    echo "Env file not found. One was created for you."
fi

echo "Bootstrapping complete."