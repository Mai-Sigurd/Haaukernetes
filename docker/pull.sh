echo "Logging in to private imagerepo - this requires a digitalocean token in '~/do_token'"
docker login -u $(cat ~/do_token) -p $(cat ~/do_token) registry.digitalocean.com

echo "Pulling images for kali and challenges"
docker pull registry.digitalocean.com/haaukins-kubernetes-bsc/logon
docker tag registry.digitalocean.com/haaukins-kubernetes-bsc/logon logon 
docker pull registry.digitalocean.com/haaukins-kubernetes-bsc/heartbleed
docker tag registry.digitalocean.com/haaukins-kubernetes-bsc/heartbleed heartbleed 
docker pull registry.digitalocean.com/haaukins-kubernetes-bsc/kali-vnc
docker tag registry.digitalocean.com/haaukins-kubernetes-bsc/kali-vnc kali-vnc 

echo "Pulling public image for wireguard"
#docker login docker.io
docker pull masipcat/wireguard-go:latest

echo "Removing duplicated images (tagging)"
docker image rm registry.digitalocean.com/haaukins-kubernetes-bsc/logon
docker image rm registry.digitalocean.com/haaukins-kubernetes-bsc/heartbleed
docker image rm registry.digitalocean.com/haaukins-kubernetes-bsc/kali-vnc