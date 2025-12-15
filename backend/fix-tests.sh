#!/bin/bash

echo "Fixing test type mismatches..."

# Fix clickhouse_test.go
find . -name "clickhouse_test.go" -exec sed -i.bak \
  -e 's/Tags:.*map\[string\]interface{}{}/Tags: map[string]string{}/g' \
  -e 's/DurationMs:.*uint32(\([^)]*\))/DurationMs: int64(\1)/g' \
  -e 's/TotalTokens:.*uint32(\([^)]*\))/TotalTokens: int(\1)/g' {} \;

# Fix trace_service_test.go  
find . -name "trace_service_test.go" -exec sed -i.bak \
  -e 's/service\.calculateCost([^,]*, [^,]*, tt\.promptTokens, tt\.completionTokens)/service.calculateCost(\1, \2, int(tt.promptTokens), int(tt.completionTokens))/g' {} \;

# Fix test_repository.go
find . -name "test_repository.go" -exec sed -i.bak \
  -e 's/Tags:.*map\[string\]interface{}{}/Tags: map[string]string{}/g' \
  -e 's/StartTime:.*startTime\.Format[^,]*/StartTime: startTime/g' \
  -e 's/EndTime:.*endTime\.Format[^,]*/EndTime: endTime/g' {} \;

echo "Done! Test files fixed."
