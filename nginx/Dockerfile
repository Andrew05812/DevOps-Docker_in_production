FROM nginx:1.27-alpine

COPY conf.d/default.conf /etc/nginx/conf.d/default.conf

COPY --from=notes-frontend /usr/share/nginx/html /usr/share/nginx/html
