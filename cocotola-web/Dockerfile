FROM node:14.21-slim as builder

WORKDIR /app
COPY package.json package-lock.json .npmrc ./
RUN npm install

COPY . .
RUN npm run-script build

# Application image
FROM nginx:1.21.4

COPY --from=builder /app/dist /usr/share/nginx/html

COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/default.conf /etc/nginx/conf.d/default.conf

RUN groupadd -r docker && useradd -r -g docker docker && \
    chown -R docker: /var/cache/nginx \
    && chmod -R g+w /var/cache/nginx

USER docker

EXPOSE 8080

CMD ["/usr/sbin/nginx", "-g", "daemon off;"]
