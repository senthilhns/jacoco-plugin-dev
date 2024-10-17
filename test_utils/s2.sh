
GOL_IMAGE_S2=gol_t3_step2

docker build --network host -f GameOfLifeTestDockerfile-step2 -t $GOL_IMAGE_S2 .
docker tag $GOL_IMAGE_S2 senthilhns/$GOL_IMAGE_S2:latest
docker push senthilhns/$GOL_IMAGE_S2:latest

