# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

FROM alpine:3.7
WORKDIR /app
ADD . /app
EXPOSE 9002
ENV PORT 9002
ENV COUNTING_SERVICE_URL http://counting.service.consul:9001
CMD ["./dashboard-service"]
