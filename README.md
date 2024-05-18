## Route Scheduler
This particular implementation of the requirements is essentially a configurable number of travelling salesman each choosing their next nearest-neighbor from a shared hat, until they're all gone, or no driver has time in their (default 12 hour) shift to deliver the remaining pickups. It's not fast. More importantly, it's not well-optimized: it takes 20 drivers 12 hours to deliver 50 loads at about 50% efficiency, where efficiency is the quotient `{pickup_to_dropoff}/({here_to_pickup+pickup_to_dropoff})`. And because of the randomness of goroutines, even that doesn't always work.

A proper solution would use a map/reduce strategy over the massive number of possible permutations in a set of 200 start/end pairs (is it `40000!`). The 12-hour time constraint lets us truncate a lot of those paths into something manageable. Haven't quite had time to sort all that out, yet, but it would yield the most accurate results. Hopefully fast enough to be useful.

It is possible to run the complete suite successfully in a 12 hour shift, but it requires 190 vehicles.

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
- FLEET_SIZE: how many vehicles are servicing this shift (default: 12)
- HOURS_PER_SHIFT: self-explanatory (default: 12)
- LOG_LEVEL: case-insensitive (trace, debug, warn, ...) (default: Debug)
- ORIGIN_X: x-coordinate of the origin (default: 0)
- ORIGIN_Y: y-coordinate of the origin (default: 0)

So, any of the following commands can be preceeded on the command-line with 0 or more of the these variables (copy/paste the needful ones), and we won't clutter each command with these optional configs:
```sh
[FLEET_SIZE=nn] [HOURS_PER_SHIFT=nn] [LOG_LEVEL=...] [ORIGIN_X=nn] [ORIGIN_Y=nn] cmd ...
```
That being said, it's our habit to always include these two in most commands:
```sh
FLEET_SIZE=20 HOURS_PER_SHIFT=12 make test
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
