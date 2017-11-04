#!/bin/bash

source /etc/apache2/envvars
chmod 777 -R /projects
sudo service ssh start
ssh-keygen -t rsa -f ~/.ssh/id_rsa -q -P ""
ssh-keygen -y -f ~/.ssh/id_rsa > ~/.ssh/id_rsa.pub
cat ~/.ssh/id_rsa.pub > ~/.ssh/authorized_keys

"$@"