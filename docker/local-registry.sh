#use 'curl -k localhost:5000/v2/_catalog' to list all images in repo

# this is "async", but dont know if it is faster than synchronous...
#er det muligt at wrappe hele det indre loop i () & 
#det er lidt et problem at terminalen "hænger" når det er slut, men det virker
docker run -d -p 5000:5000 --restart=always --name registrytest registry:2

docker tag kali localhost:5000/kali
docker push localhost:5000/kali
docker image rm localhost:5000/kali

for i in {1..5}
do
    (docker tag logon localhost:5000/logon$i;
    docker push localhost:5000/logon$i; 
    docker image rm localhost:5000/logon$i) &

    (docker tag heartbleed localhost:5000/heartbleed$i;
    docker push localhost:5000/heartbleed$i;
    docker image rm localhost:5000/heartbleed$i) &

    (docker tag for-fun-and-profit localhost:5000/for-fun-and-profit$i;
    docker push localhost:5000/for-fun-and-profit$i;
    docker image rm localhost:5000/for-fun-and-profit$i) &
    
    (docker tag hide-and-seek localhost:5000/hide-and-seek$i;
    docker push localhost:5000/hide-and-seek$i;
    docker image rm localhost:5000/hide-and-seek$i) &
    
    (docker tag program-behaviour localhost:5000/program-behaviour$i;
    docker push localhost:5000/program-behaviour$i;
    docker image rm localhost:5000/program-behaviour$i) &
    
    (docker tag reverseapk localhost:5000/reverseapk$i;
    docker push localhost:5000/reverseapk$i;
    docker image rm localhost:5000/reverseapk$i) &
done