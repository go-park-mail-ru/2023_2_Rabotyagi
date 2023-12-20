#!/bin/sh
COVERAGE=$(go test --race -coverpkg=./... -coverprofile=cover.out ./... &&
  cat cover.out | grep -v "mocks" | grep -v "easyjson" | grep -v ".pb" > pure_cover.out &&
  go tool cover --func pure_cover.out)
COVERAGE_PURE=$(echo $COVERAGE | grep -oP '\(statements\)(\s+)(\d+\.\d+)' | grep -oP '(\d+\.\d+)')

rm cover.out;
rm pure_cover.out;

if [ $(echo "${COVERAGE_PURE} < 60.00" | bc) -eq 1 ]
then
  echo "Coverage is $COVERAGE_PURE less than 60%";
  exit 1;
else
  echo "Coverage is $COVERAGE_PURE more than 60%";
fi
