FROM debian:jessie

RUN apt-get update; apt-get install -y curl php5-cli git apache2 php5; \
curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer; \
useradd -ms /bin/bash -p test adminuser; \
rm -rfv /var/www/html/*; \
git clone https://github.com/Codiad/Codiad /var/www/html/; \
touch /var/www/html/config.php; \
chown www-data:www-data -R /var/www/html/;

COPY ./entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]

WORKDIR /app

CMD ["/usr/sbin/apache2", "-D", "FOREGROUND"]