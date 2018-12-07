#!/bin/bash

if [ -z "$1" ]
  then
    echo "Must provide local CSI Driver endpoint"
    exit 1
fi

endpoint=$1
echo "Endpoint: $endpoint"

printf "\n\n###########################\n"
printf "    RUNNING SANITY TESTS\n"
printf "###########################\n\n\n"

./csi-sanity -csi.endpoint=$endpoint

printf "\n\n###########################\n"
printf "     RUNNING E2E TESTS\n"
printf "###########################\n\n\n"

if ["$DRIVER_MANIFEST" = ""]
  then
    printf "\nMust set DRIVER_MANIFEST environment varibale to run e2e-tests\n\n"
    exit 1
fi

(cd ../../csi-e2e/ && GOCACHE=off go test -v ./test/e2e)

printf "\n\n###########################\n"
printf "     E2E TESTS COMPLETE\n"
printf "###########################\n\n\n"
