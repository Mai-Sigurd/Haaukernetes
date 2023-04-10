docker build -f ./admin_logon/Dockerfile -t benis_logon ./admin_logon
docker build -f ./heartbleed/Dockerfile -t benis_heartbleed ./heartbleed 
docker build -f ./kali_dockerfile -t benis_kali-vnc . 
docker pull masipcat/wireguard:latest -t benis_wireguard
