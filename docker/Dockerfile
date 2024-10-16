
FROM eclipse-temurin:8-jre-alpine AS alpine

RUN apk add -U --no-cache ca-certificates curl

ENV HARNESS_PLUGIN_DEPS_PATH="/opt/harness/plugins-deps"
ENV JACOCO_VERSION="0.8.12"
ENV JACOCO_PATH="${HARNESS_PLUGIN_DEPS_PATH}/jacoco/${JACOCO_VERSION}"

RUN mkdir -p ${HARNESS_PLUGIN_DEPS_PATH}
RUN mkdir -p ${JACOCO_PATH}

RUN curl -L -o ${JACOCO_PATH}/org.jacoco.cli-${JACOCO_VERSION}.jar \
    https://repo1.maven.org/maven2/org/jacoco/org.jacoco.cli/0.8.12/org.jacoco.cli-0.8.12-nodeps.jar

RUN ln -snf ${JACOCO_PATH}/org.jacoco.cli-${JACOCO_VERSION}.jar \
    ${JACOCO_PATH}/jacoco.jar

FROM eclipse-temurin:21-jre-alpine

ENV GODEBUG netdns=go
ENV HARNESS_PLUGIN_DEPS_PATH="/opt/harness/plugins-deps"
ENV JACOCO_PATH="${HARNESS_PLUGIN_DEPS_PATH}/jacoco/${JACOCO_VERSION}"

COPY --from=alpine ${JACOCO_PATH} ${JACOCO_PATH}
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD release/linux/amd64/plugin /bin/

ENTRYPOINT ["/bin/plugin"]
