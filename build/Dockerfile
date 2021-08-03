FROM registry.access.redhat.com/ubi8/ubi-minimal

ENV TZ="Europe/Zurich" \
  LANG="en_US.UTF-8" \
  WEBHOOK=/usr/local/bin/webhook \
  UID=1001:1001

COPY bin/webhook ${WEBHOOK}

USER ${UID}

CMD ["${WEBHOOK}"]
