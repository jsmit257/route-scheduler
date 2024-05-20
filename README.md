## Route Scheduler
Still a travelling salesman. This time around, we 
- try every delivery as a possible starting point, then
- calculate a nearest neighbor path through all the remaining points, and
- split the shortest path into segments that can all be delivered in one 12-hour shift

In a perfect world, the result of the nearest neighbor is just a limit to set when using a permutational approach so paths can be truncated when they're already longer than a known-shortest path, thus saving a lot of extra work. This actually works and is deployed, but it's only practical for small data sets: for the 10 records in `problem1.txt` this algorithm cuts the mean cost by about 1300, but it takes much longer. For large sets - anything greater than 10 deliveries - the current implementation is either killed by the test timeout (default: 10m), or it overflows the max number of goroutines. Remember that `50!` is about `3*10^64` combinations, so a naive approach would overflow 64-bit stackframe-IDs by a lot.

_tl;dr_
To run this, just:
```sh
make run-docker
```
Or, to run it with more noise (default is info, skipping trace & debug messages):
```sh
LOG_LEVEL=trace make run-docker
```

### Requirements
This project can run completely in [docker](https://hub.docker.com/repository/docker/jsmit257/rs/general), or with no docker at all, so not all the following are necessary.
- golang development environment (for building/non-docker testing)
- python3 runtime (for local (non-docker) testing)
- docker/docker-compose
- GNU/Posix compatible shell tools (make/grep/etc...)
- `jq`: very optional, only used by [bin/shift-cost.sh](./bin/shift-cost.sh), but it's also awesome!

### Using

#### Environment variables
These are passed through either `make` or `docker-compose.yml` on the command-line. All are optional with sane defaults.
- HOURS_PER_SHIFT: self-explanatory (default: 12)
- LOG_LEVEL: case-insensitive (trace, debug, warn, ...) (default: Debug)
- ORIGIN_X: x-coordinate of the origin (default: 0)
- ORIGIN_Y: y-coordinate of the origin (default: 0)

So, any of the following commands can be preceeded on the command-line with 0 or more of the these variables (copy/paste the needful ones), and we won't clutter each command with these optional configs:
```sh
[HOURS_PER_SHIFT=nn] [LOG_LEVEL=...] [ORIGIN_X=nn] [ORIGIN_Y=nn] cmd ...
```
That being said, it's our habit to always include these two in most commands:
```sh
HOURS_PER_SHIFT=12 make test
```

#### Docker
For the full suite of tests in docker using the `eveluateShared.py` script, run the fillowing command from the project root:
```sh
make run-docker
```
An image has been pushed to [dockerhub](https://hub.docker.com/repository/docker/jsmit257/rs/general), so `--build` is not technically necessary unless you're developing new stuff. 

Sample files are copied into the image's `/data` directory, so a line like this will run an individual problem sample:
```sh
docker-compose run [--build] run-docker /rs /data/problem1.txt
```
Assume everyone here knows what the problem filenames look like and how many there are. 

Also, since it can read from `stdin` *instead of* a command-line filename, and noisily ignores unparseable lines (aka, header rows), you could combine problems:
```sh
docker-compose run [--build] run-docker /bin/bash -c 'cat /data/problem{1,3,6}.txt | /rs'
```

#### Native
The above solutions all require `docker`/`docker-compose`, and only require the golang environment if `--build` is specified. These options require golang and/or python3, but no docker suite
```sh
make run-local
# shorthand for:
# $ python3 ./bin/evaluateShared.py --cmd "go run ./cmd/..." --problemDir data
```

If you're only looking to test a file w/o having it evaluated or without building a new docker image, just:
```sh
go run ./cmd/... ./data/problem1.txt # or $ cat ./data/problem[1..5].txt | go run ./cmd/... 
```
