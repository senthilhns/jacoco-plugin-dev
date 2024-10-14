
Dev testing for JDK Docker Image
================================
docker build --network host -t plugins/drone-coverage-report -f docker/Dockerfile.DevTest.Amd64 .

Production JRE support Docker Image
===================================
docker build --network host -t plugins/drone-coverage-report -f docker/Dockerfile .


Debug and Test Locally
======================
docker run -it -v /tmp:/tmp --entrypoint /bin/sh plugins/drone-coverage-report

JDK setup for testing
=====================

sudo apt install openjdk-8-jdk -y
sudo update-alternatives --config java
sudo update-alternatives --config javac

export MAVEN_OPTS="-Xmx1024m -XX:MaxPermSize=512m"
mvn clean verify jacoco:report

Test case repo
==============
https://github.com/syamv/game-of-life


java -jar jacoco.jar \
    report   ./gameoflife-core/target/jacoco.exec   ./gameoflife-web/target/jacoco.exec   \
    --classfiles ./gameoflife-core/target/classes   \
    --sourcefiles ./gameoflife-core/src/main/java   \
    --html ./gameoflife-core/target/site/jacoco_html   \
    --xml ./gameoflife-core/target/site/jacoco.xml

