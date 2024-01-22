FROM node:21 as tailwind-builder

WORKDIR /tmp/tailwind

RUN npm install -g tailwindcss

COPY ./tailwind.config.js ./
COPY ./src/templates ./templates
COPY ./style.input.css ./

RUN npx tailwindcss -i style.input.css -o style.css

###
FROM nginx:1.25

COPY ./static/ /usr/share/nginx/html/
COPY --from=tailwind-builder /tmp/tailwind/style.css /usr/share/nginx/html/style.css
COPY ./default.nginx.conf /etc/nginx/templates/default.conf.template
