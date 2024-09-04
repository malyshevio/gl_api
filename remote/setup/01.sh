#!/bin/bash
set -eu

TIMEZONE=Asia/Yekaterinburg

USERNAME=gl

read -p "Enter password for gl_api_dsn: " DB_PASSWORD

export LC_ALL=en_US.UTF-8

# ================================================================================== #
# SCRIPT
# ================================================================================== #

add-apt-repository --yes universe

apt update
apt --yes -o Dpkg::Options::='--force-confnew' upgrade

timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"

passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"

rsync --archive --chown=${USERNAME}:${USERNAME} /home/atomskjd/.ssh /home/${USERNAME}

ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

apt --yes install fail2ban

curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate /usr/local/bin/migrate

apt --yes install postgresql

sudo -i -u postgresql psql -c "CREATE DATABASE gl_api"
sudo -i -u postgresql psql -d gl_api -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgresql psql -d gl_api -c "CREATE ROLE gl_api WITH LOGIN PASSWORD '${DB_PASSWORD}'"

echo "GL_API_DSN='postgres://gl_api:${DB_PASSWORD}@localhost/gl_api'" >> /etc/environment

apt --yes install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list

apt update
apt --yes install caddy

echo "COMPLETED!!!"

reboot