FROM debian:jessie

RUN apt-get update; apt-get install -y sudo openssh-server curl php5-cli git apache2 php5 build-essential && \
mkdir /var/run/sshd && \
sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd && \
echo "%sudo ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers && \
useradd -u 1000 -G users,sudo -d /home/user --shell /bin/bash -m user && \
echo "secret\nsecret" | passwd user && \
curl -sL https://deb.nodesource.com/setup_6.x | bash - && \
curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer; \
useradd -ms /bin/bash -p test adminuser; \
rm -rfv /var/www/html/*; \
git clone https://github.com/Codiad/Codiad /var/www/html/; \
touch /var/www/html/config.php; \
chown www-data:www-data -R /var/www/html/ && \
apt-get install -y nodejs && \
apt-get clean && \
apt-get -y autoremove && \
rm -rf /var/lib/apt/lists/* && \
git clone https://github.com/krishnasrinivas/wetty /opt/wetty && \
cd /opt/wetty && npm install && \
sed -i 's/<VirtualHost \*:80>/<VirtualHost \*:8000>/g' /etc/apache2/sites-available/000-default.conf

USER user

COPY ./entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]

WORKDIR /app

CMD node /opt/wetty -p 3000; /usr/sbin/apache2 -D FOREGROUND