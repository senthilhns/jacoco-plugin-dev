# Dev testing for JDK Docker Image

```bash
docker build --network host -t plugins/drone-coverage-report -f docker/Dockerfile.DevTest.Amd64 .
```

# Production JRE support Docker Image

```bash
docker build --network host -t plugins/drone-coverage-report -f docker/Dockerfile .
```

# Debug and Test Locally

```bash
docker run -it -v /tmp:/tmp --entrypoint /bin/sh plugins/drone-coverage-report
```

# JDK setup for testing

```bash
sudo apt install openjdk-8-jdk -y
sudo update-alternatives --config java
sudo update-alternatives --config javac

export MAVEN_OPTS="-Xmx1024m -XX:MaxPermSize=512m"
mvn clean verify jacoco:report
```

# Test case repo

[https://github.com/syamv/game-of-life](https://github.com/syamv/game-of-life)

# Generate report with class and source files

```bash
java -jar jacoco.jar \
    report ./gameoflife-core/target/jacoco.exec ./gameoflife-web/target/jacoco.exec \
    --classfiles ./gameoflife-core/target/classes \
    --sourcefiles ./gameoflife-core/src/main/java \
    --html ./gameoflife-core/target/site/jacoco_html \
    --xml ./gameoflife-core/target/site/jacoco.xml
```

# Generate report with only class files and no source files

```bash
java -jar jacoco.jar \
    report ./gameoflife-core/target/jacoco.exec ./gameoflife-web/target/jacoco.exec \
    --classfiles ./gameoflife-core/target/classes \
    --html /tmp/jacoco_html \
    --xml /tmp/jacoco.xml
```

# Test case for quick check

```bash
rm -rvf /opt/hns/test-resources/game-of-life-master/jacoco-workspace--* && \
BUILD_ROOT_PATH=/opt/hns/test-resources/game-of-life-master go test -count=1 -run ^TestSourcePathWithIncludeAndExclude$
```
