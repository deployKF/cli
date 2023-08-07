FROM alpine:3.18.2

ARG TARGETOS
ARG TARGETARCH

# args -- software versions
ARG DKF_CLI_OS=${TARGETOS}
ARG DKF_CLI_ARCH=${TARGETARCH}
ARG DKF_CLI_RELEASE_TAG=v0.1.0

# args -- uid/gid
ARG DEPLOYKF_USER=deploykf
ARG DEPLOYKF_GROUP=deploykf
ARG DEPLOYKF_UID=1001
ARG DEPLOYKF_GID=1001
ARG DEPLOYKF_HOME=/home/${DEPLOYKF_USER}

# install deploykf cli
RUN wget -q -O /tmp/deploykf "https://github.com/deploykf/cli/releases/download/${DKF_CLI_RELEASE_TAG}/deploykf-${DKF_CLI_OS}-${DKF_CLI_ARCH}" \
 && wget -q -O /tmp/deploykf.sha256 "https://github.com/deploykf/cli/releases/download/${DKF_CLI_RELEASE_TAG}/deploykf-${DKF_CLI_OS}-${DKF_CLI_ARCH}.sha256" \
 && echo "$(cat /tmp/deploykf.sha256 | awk '{ print $1; }')  /tmp/deploykf" | sha256sum -c - \
 && chmod +x /tmp/deploykf \
 && mv /tmp/deploykf /usr/local/bin/deploykf \
 && rm -rf /tmp/deploykf*

# create non-root 'deploykf' user/group
RUN addgroup -g ${DEPLOYKF_GID} "${DEPLOYKF_GROUP}" \
 && adduser -D -h "${DEPLOYKF_HOME}" -u ${DEPLOYKF_UID} -G ${DEPLOYKF_GROUP} "${DEPLOYKF_USER}"

USER ${DEPLOYKF_UID}:${DEPLOYKF_GID}
WORKDIR ${DEPLOYKF_HOME}

ENTRYPOINT ["/usr/local/bin/deploykf"]