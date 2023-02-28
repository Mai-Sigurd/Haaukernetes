# Køre kali-vnc docker image for sig selv
`docker build . -t kali-vnc`
`docker run -d --name vnc -p 5901:5901 kali-vnc`

# Kubernetes setup
- Ideen er, at logon pod (exercise, logon.yaml) bliver exposed
"internt" i k8s netværket via logon_service.yaml
- vnc pod (vnc_deploy.yaml) udstilles også som en service i k8s netværket via vnc_service.yaml og helt ud til host via vnc_expose.yaml
- På den måde er exercise kun tilgængelig internt og kali/vnc er tilgængelig via en vnc client

## Køre setup
- Sørg for at minikube kører samt at docker images er bygget i minikube docker-env
- Sørg evt. for at der ikke kører andet i minikube
- Kør `kubectl apply -f` på alle .yaml filerne / brug de to scripts `start_vnc.sh` og `start_logon.sh` (eksisterer nu også for heartbleed, både stop og start)

## Forbinde til vnc 
- Kræver en vnc client (have no idea hvad man får til mac)
- minikube-ip:32320 (port specificeret i vnc_expose.yaml) kan indsættes i vnc client hvorefter "kali" er kodeordet
    - OBS: minikube-ip findes ved at køre `echo $(minikube ip)`

## Forbinde til exercise
- Det er navnet på pod samt port der afgør hvordan man finder den
- I dette tilfælde finder man den via url "logon:80" (portnummer specificeret i .yaml), i browseren i kali
