
1) Build the Game of Life project by running the following command in the terminal, this will install
the necessary jacoco reports jacoco.html and jacoco.xml in the target directory:
sh ./build_game_of_life.sh

2) Test and check whether Game of Life image is working correctly
docker run --rm -e DWPJA2=$DWPJA2 game-of-life-jacoco-image

docker run -it --rm -e DWPJA2=$DWPJA2 game-of-life-jacoco-image /bin/bash

docker run -it --name game_of_life_container --rm -e DWPJA2=$DWPJA2 game-of-life-jacoco-image /bin/sh
