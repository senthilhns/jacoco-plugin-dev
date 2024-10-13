
docker build --network host -t plugins/drone-coverage-report -f docker/Dockerfile .

Debug using
docker run -it --entrypoint /bin/sh plugins/drone-coverage-report

Test Locally
docker run -it -v /tmp:/tmp --entrypoint /bin/sh plugins/drone-coverage-report

