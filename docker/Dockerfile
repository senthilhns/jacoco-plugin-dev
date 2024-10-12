
FROM eclipse-temurin:21-jre-alpine AS alpine
RUN apk add -U --no-cache ca-certificates curl


RUN mkdir -p /opt/harness/plugins-deps


# Plugin specific steps add here
RUN mkdir -p /opt/harness/plugins-deps/jacoco/0.8.13


RUN curl -L -o /opt/harness/plugins-deps/jacoco/0.8.13/org.jacoco.cli-0.8.13.jar \
    https://repo1.maven.org/maven2/org/jacoco/org.jacoco.cli/0.8.13/org.jacoco.cli-0.8.13.jar

RUN ln -snf /opt/harness/plugins-deps/jacoco/0.8.13/org.jacoco.cli-0.8.13.jar \
    /opt/harness/plugins-deps/jacoco/0.8.13/jacoco.jar

FROM eclipse-temurin:21-jre-alpine

ENV GODEBUG netdns=go

COPY --from=alpine /opt/harness/plugins-deps/jacoco/0.8.13 /opt/harness/plugins-deps/jacoco/0.8.13

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD release/linux/amd64/plugin /bin/

ENTRYPOINT ["/bin/plugin"]

