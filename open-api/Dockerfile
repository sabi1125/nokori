# Stage 1: bundle the split spec (spec/) into the single file Scalar
# needs. docker-compose also runs this stage on its own as the dev
# watcher (target: bundler).
FROM node:22-alpine AS bundler
ENV REDOCLY_TELEMETRY=off HOME=/tmp
RUN npm install -g @redocly/cli@2.57.0
COPY scripts/watch-bundle.sh /usr/local/bin/watch-bundle
WORKDIR /open-api
COPY spec/ spec/
RUN mkdir public && redocly bundle spec/openapi.yaml -o public/openapi.yaml

# Stage 2: serve index.html + the bundled spec.
FROM nginx:alpine
COPY public/index.html /usr/share/nginx/html/
COPY --from=bundler /open-api/public/openapi.yaml /usr/share/nginx/html/
EXPOSE 80
