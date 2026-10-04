# Stage 1: Redocly CLI + the dev watcher. docker-compose runs this stage
# on its own as open-api-bundler (target: bundler); it holds no spec, so
# a broken spec never stops it from building.
FROM node:22-alpine AS bundler
ENV REDOCLY_TELEMETRY=off HOME=/tmp
RUN npm install -g @redocly/cli@2.57.0
COPY scripts/watch-bundle.sh /usr/local/bin/watch-bundle
WORKDIR /open-api

# Stage 2: bundle the split spec (spec/) into the single file Scalar
# needs. A broken spec fails the deploy build here, on purpose.
FROM bundler AS bundle
COPY spec/ spec/
RUN mkdir public && redocly bundle spec/openapi.yaml -o public/openapi.yaml

# Stage 3: serve index.html + the bundled spec.
FROM nginx:alpine
COPY public/index.html /usr/share/nginx/html/
COPY --from=bundle /open-api/public/openapi.yaml /usr/share/nginx/html/
EXPOSE 80
