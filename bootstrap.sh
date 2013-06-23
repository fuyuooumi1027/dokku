DOKKU_REPO=${DOKKU_REPO:-"https://github.com/progrium/dokku.git"}
DOKKU_STACK=${DOKKU_STACK:-"https://s3.amazonaws.com/progrium-dokku/progrium_buildstep.tgz"}
DOCKER_PKG=${DOCKER_PKG:-"https://launchpad.net/~dotcloud/+archive/lxc-docker/+files/lxc-docker_0.4.2-1_amd64.deb"}
set -e

apt-get update
DEBIAN_FRONTEND=noninteractive apt-get install -y linux-image-extra-`uname -r`
apt-get install -y git nginx make curl dnsutils

wget -qO- "$DOCKER_PKG" > /tmp/lxc-docker_0.4.2-1_amd64.deb
dpkg --force-depends -i /tmp/lxc-docker_0.4.2-1_amd64.deb && apt-get install -f -y
rm /tmp/lxc-docker_0.4.2-1_amd64.deb

sed -i.bak 's/docker -d/docker -d -r=true/' /etc/init/docker.conf # Docker should restart containers 

cd ~ && git clone ${DOKKU_REPO}
cd dokku && make install
if [[ $DOKKU_STACK ]]; then
  curl "$DOKKU_STACK" | gunzip -cd | docker import - progrium/buildstep
else
  cd buildstep && make build
fi

[[ -f /etc/nginx/nginx.conf ]] && sed -i 's/# server_names_hash_bucket_size/server_names_hash_bucket_size/' /etc/nginx/nginx.conf

/etc/init.d/nginx start
start nginx-reloader

[[ $(dig +short $HOSTNAME) ]] && HOSTFILE=DOMAIN || HOSTFILE=HOSTNAME
echo $HOSTNAME > /home/git/$HOSTFILE

echo
echo "Be sure to upload a public key for your user:"
echo "  cat ~/.ssh/id_rsa.pub | ssh root@$HOSTNAME \"gitreceive upload-key progrium\""
