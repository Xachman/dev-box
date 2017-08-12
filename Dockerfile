FROM node

RUN apt-get update; apt-get install -y curl php5-cli git; \
git clone https://github.com/krishnasrinivas/wetty; \
curl -sS https://getcomposer.org/installer | exiphp -- --install-dir=/usr/local/bin --filename=composer; \
npm install gulp -g; \
cd wetty; \
npm install; \
useradd -ms /bin/bash -p test adminuser

WORKDIR /app

CMD node /wetty/app.js -p 3000