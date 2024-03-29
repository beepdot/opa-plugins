# Copyright 2018 The OPA Authors. All rights reserved.
# Use of this source code is governed by an Apache2
# license that can be found in the LICENSE file.

ARG BASE

FROM ${BASE}

# Any non-zero number will do, and unfortunately a named user will not, as k8s
# pod securityContext runAsNonRoot can't resolve the user ID:
# https://github.com/kubernetes/kubernetes/issues/40958. Make root (uid 0) when
# not specified.
ARG USER=0

USER ${USER}

WORKDIR /app

COPY opa /app

ENTRYPOINT ["./opa"]

CMD ["run"]
