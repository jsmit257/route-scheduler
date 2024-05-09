## Route Scheduler

### Requirements
This project can run completely from [dockerhub][dockerhub.io/), or with no docker at all, so not all the following are necessary.

- golang development environment (for building)
- python3 runtime (for local (non-docker) testing)
- docker/docker-compose
- GNU/Posix compatible shell tools (make/grep/etc...)

### Using
The environment variable FLEET_SIZE on the command line or in `docker-compose.yml` determines how many vehicles are servicing the problem manifest.

For the full suite of tests in docker, run the fillowing command from the project root:
```sh
docker-compose up [--build] --remove-orphans run-docker
# or (note that this forces a build of the docker image), just:
make run-docker
```
An image has been pushed to [docker-hub](dockerhub.io), so `--build` is not technically necessary; 

Sample files are copied into the image's `/data` directory, so a line like this will run an individual problem sample:
```sh
docker-compose run [--build] run-docker /rs /data/problem1.txt
```
Assume everyone here knows what the problem filenames look like and how many there are. 

Also, since it can read from `stdin` *instead of* a command-line filename, and noisily ignores unparseable lines (aka, header rows), you could also combine problems:
```sh
docker-compose run [--build] run-docker /bin/bash -c 'cat /data/problem{1,3,6}.txt | /rs'
```
<!-- or even, from the project root:
```sh
cat ./data/problem{1,3,6}.txt | docker-compose exec run-docker /rs
``` -->

The above solutions all require `docker`/`docker-compose`, and only require the golang environment if `--build` is specified. This option requires golang and python3, but no `docker*`
```sh
[FLEET_SIZE=nn] make run-local
# shorthand for:
# $ python3 ./bin/evaluateShared.py --cmd "go run ./cmd/..." --problemDir data
```
