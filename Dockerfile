FROM golang:1.13.5

ARG USER_ID=1000
ARG GROUP_ID=1000

RUN apt-get update && apt-get install -y wget fonts-noto
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add -
RUN echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list
RUN apt-get update && apt-get install -y google-chrome-stable

ENV WORK_DIR /go/src/github.com/builtinnya/web2ppm
RUN mkdir -p "${WORK_DIR}"

ENV CONTAINER_GROUP nya
ENV CONTAINER_USER nya
RUN groupadd -g "${GROUP_ID}" "${CONTAINER_GROUP}" && \
  useradd -l -m -u "${USER_ID}" -g "${CONTAINER_GROUP}" "${CONTAINER_USER}" && \
  chown -R "${CONTAINER_USER}:${CONTAINER_GROUP}" "${WORK_DIR}"
USER "${CONTAINER_USER}"

WORKDIR ${WORK_DIR}

CMD tail -f /dev/null
