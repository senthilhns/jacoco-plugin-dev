#!/bin/bash

# Print the JVM version installed
echo "Printing the installed JVM version:"
java -version

# Check if DWPJA2 is set, otherwise exit
if [ -z "$DWPJA2" ]; then
    echo "Error: DWPJA2 environment variable is not set."
    exit 1
fi

#find $DWPJA2

# List of required files
required_files=(
    "$DWPJA2/game-of-life/gameoflife-core/target/jacoco.exec"
    "$DWPJA2/game-of-life/gameoflife-web/target/jacoco.exec"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/main/com/wakaleo/gameoflife/domain/Universe.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/main/com/wakaleo/gameoflife/domain/Cell.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/main/com/wakaleo/gameoflife/domain/GridReader.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/main/com/wakaleo/gameoflife/domain/GridWriter.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/main/com/wakaleo/gameoflife/domain/Grid.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/domain/WhenYouCreateAGrid.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/domain/WhenYouCreateANewUniverse.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/domain/WhenYouPlayTheGameOfLife.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/domain/WhenYouCreateACell.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/domain/WhenYouReadAGridFromAString.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/domain/WhenYouPrintAGrid.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/hamcrest/MyMatchers.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/hamcrest/HasSizeMatcher.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/hamcrest/WhenIUseMyCustomHamcrestMatchers.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/test/categories/SlowTests.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/test/categories/IntegrationTests.class"
    "$DWPJA2/game-of-life/gameoflife-core/build/classes/test/com/wakaleo/gameoflife/test/categories/RegressionTests.class"
    "$DWPJA2/game-of-life/gameoflife-core/target/classes/com/wakaleo/gameoflife/domain/Universe.class"
    "$DWPJA2/game-of-life/gameoflife-core/target/classes/com/wakaleo/gameoflife/domain/Cell.class"
    "$DWPJA2/game-of-life/gameoflife-core/target/classes/com/wakaleo/gameoflife/domain/GridReader.class"
    "$DWPJA2/game-of-life/gameoflife-core/target/classes/com/wakaleo/gameoflife/domain/GridWriter.class"
    "$DWPJA2/game-of-life/gameoflife-core/target/classes/com/wakaleo/gameoflife/domain/Grid.class"
    "$DWPJA2/game-of-life/gameoflife-web/target/gameoflife/WEB-INF/classes/com/wakaleo/gameoflife/webtests/controllers/HomePageController.class"
    "$DWPJA2/game-of-life/gameoflife-web/target/gameoflife/WEB-INF/classes/com/wakaleo/gameoflife/webtests/controllers/GameController.class"
    "$DWPJA2/game-of-life/gameoflife-web/target/classes/com/wakaleo/gameoflife/webtests/controllers/HomePageController.class"
    "$DWPJA2/game-of-life/gameoflife-web/target/classes/com/wakaleo/gameoflife/webtests/controllers/GameController.class"
    "$DWPJA2/game-of-life/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/Universe.java"
    "$DWPJA2/game-of-life/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/Grid.java"
    "$DWPJA2/game-of-life/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/Cell.java"
    "$DWPJA2/game-of-life/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/GridReader.java"
    "$DWPJA2/game-of-life/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/GridWriter.java"
    "$DWPJA2/game-of-life/gameoflife-web/target/site/jacoco/com.wakaleo.gameoflife.webtests.controllers/HomePageController.java.html"
    "$DWPJA2/game-of-life/gameoflife-web/target/site/jacoco/com.wakaleo.gameoflife.webtests.controllers/GameController.java.html"
    "$DWPJA2/game-of-life/gameoflife-web/src/main/java/com/wakaleo/gameoflife/webtests/controllers/HomePageController.java"
    "$DWPJA2/game-of-life/gameoflife-web/src/main/java/com/wakaleo/gameoflife/webtests/controllers/GameController.java"
)

# Check if each file exists and print status
all_files_found=true

for file in "${required_files[@]}"
do
    if [ -f "$file" ]; then
        echo "Found: $file"
    else
        echo "Missing: $file"
        all_files_found=false
    fi
done

# Print success if all files are found
if [ "$all_files_found" = true ]; then
    echo "Success: All required files are present."
else
    echo "Error: Some required files are missing."
    exit 1
fi


echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"


DWP=/harness


echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "########## COPYING FILES #############"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"

cp -rf  $DWPJA2/game-of-life $DWP/

echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "########## listing all files #############"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"
echo "#######################"

#find $DWP

find $DWP -iname "*.exec*"
find $DWP -iname "*.class*"
find $DWP -iname "*.java*"

