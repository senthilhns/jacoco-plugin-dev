echo $GOL_IMAGE

docker build --network host --build-arg DWPJA2=$DWPJA2 -f GameOfLifeTestDockerfile -t $GOL_IMAGE .
docker tag $GOL_IMAGE senthilhns/$GOL_IMAGE:latest
docker push senthilhns/$GOL_IMAGE:latest

