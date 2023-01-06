cd ..
rm -f coverage.out
rm -f coverage_without_event.out
go test -v -coverprofile=coverage.out ./tests -coverpkg=./... | grep FAIL
cat coverage.out | grep -v collector.go | grep -v requester.go | grep -v /log/ >> coverage_without_event.out
go tool cover -func=coverage_without_event.out | grep total