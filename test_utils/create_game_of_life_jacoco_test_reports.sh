#!/bin/bash

set -e

if [ -z "$DWPJA2" ]; then
    echo "Error: The DWPJA2 environment variable is not set."
    exit 1
fi

echo "Creating the workspace directory: $DWPJA2"
mkdir -p $DWPJA2

echo "Changing to the workspace directory: $DWPJA2"
cd $DWPJA2

echo "Updating the package manager and installing Git, OpenJDK 8, and Curl"
apk update && \
apk add --no-cache git openjdk8 curl

echo "Setting JAVA_HOME environment variable"
export JAVA_HOME=/usr/lib/jvm/java-1.8-openjdk

echo "Installing Maven"
MAVEN_VERSION=3.8.6
MAVEN_DOWNLOAD_URL=https://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz

# Download and install Maven
curl -fSL $MAVEN_DOWNLOAD_URL -o /tmp/apache-maven.tar.gz
if [ $? -ne 0 ]; then
    echo "Error downloading Maven. Please check the URL."
    exit 1
fi

tar -xzf /tmp/apache-maven.tar.gz -C /opt
ln -s /opt/apache-maven-$MAVEN_VERSION/bin/mvn /usr/bin/mvn

echo "Setting Maven options"
export MAVEN_OPTS="-Xmx1024m -XX:MaxPermSize=512m"

echo "Cloning the repository into $DWPJA2/game-of-life"
git clone https://github.com/syamv/game-of-life.git $DWPJA2/game-of-life

echo "Changing to the repository directory: $DWPJA2/game-of-life"
cd $DWPJA2/game-of-life

echo "Running Maven clean, verify, and JaCoCo report"
mvn clean verify jacoco:report

echo "Maven build and JaCoCo report generation complete!"
