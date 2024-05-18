#!/bin/sh

for f in data/*; do \
  HOURS_PER_SHIFT=80 FLEET_SIZE=100 go run ./cmd/... $f
    if [ $? -ne 0 ]; then
      break;
    fi;
done 2>&1 \
| grep '"done!"' \
| jq '.shift_cost' \
| awk '{s+=$1} END {printf "%.0f\n", s}'
