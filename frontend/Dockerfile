FROM node:24-alpine AS build

WORKDIR /app

ARG VITE_API_URL=/api
ENV VITE_API_URL=${VITE_API_URL}

COPY package.json package-lock.json ./
RUN npm ci

COPY index.html vite.config.ts tsconfig*.json ./
COPY public ./public
COPY src ./src

RUN npm run build

FROM nginx:alpine AS runtime

ENV API_UPSTREAM=http://backend:8080
ENV PORT=80

RUN set -eux; \
  mkdir -p /etc/nginx/templates; \
  printf '%s\n' \
    'server {' \
    '  listen ${PORT};' \
    '  server_name _;' \
    '  root /usr/share/nginx/html;' \
    '  index index.html;' \
    '' \
    '  location /api/ {' \
    '    resolver ${NGINX_LOCAL_RESOLVERS} valid=30s;' \
    '    set $api_upstream ${API_UPSTREAM};' \
    '    proxy_pass $api_upstream;' \
    '    proxy_http_version 1.1;' \
    '    proxy_set_header Host $host;' \
    '    proxy_set_header X-Real-IP $remote_addr;' \
    '    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;' \
    '    proxy_set_header X-Forwarded-Proto $scheme;' \
    '  }' \
    '' \
    '  location / {' \
    '    try_files $uri $uri/ /index.html;' \
    '  }' \
    '}' \
    > /etc/nginx/templates/default.conf.template

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
