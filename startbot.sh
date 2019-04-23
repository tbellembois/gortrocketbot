#!/bin/bash

export ROCKET_SERVERHOST="rocket.foo.com"
export ROCKET_SERVERSCHEME="https"
export ROCKET_USER="mybot"
export ROCKET_EMAIL="mybot@uca.fr"
export ROCKET_PASSWORD="mybotpassword"
export ROCKETP_HELLO_LANGUAGE="fr"

#
# add here the exports of the used plugins
#

/usr/local/gortrocketbot/gortrocketbot
