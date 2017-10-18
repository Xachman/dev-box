#!/bin/bash

source /etc/apache2/envvars
chmod 777 -R /projects
sudo service ssh start

"$@"