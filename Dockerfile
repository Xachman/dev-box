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
sed -i 's/<VirtualHost \*:80>/<VirtualHost \*:8000>/g' /etc/apache2/sites-available/000-default.conf && \
sed -i "s/Listen 80/Listen 8000/g" /etc/apache2/ports.conf

RUN curl -fsSL get.docker.com -o /opt/get-docker.sh && sudo sh /opt/get-docker.sh
RUN usermod -aG docker user
VOLUME /projects

COPY ./entrypoint.sh /usr/local/bin
COPY ./bash/start_container.sh /usr/local/bin/start_container
RUN chmod +x /usr/local/bin/entrypoint.sh;
RUN chmod +x /usr/local/bin/start_container;

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

WORKDIR /dashboard

USER user

CMD cat /home/user/.ssh/id_rsa & /app/src/main & npm start