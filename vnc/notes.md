Planen 
- der laves en deployment/pod med kali-vnc image
- der laves en deployment/pod med logon
- der laves en service der forbinder de to jf. https://www.tutorialworks.com/kubernetes-pod-communication/ -> altså: en service "rundt om" logon bør give vnc pod mulighed for at tilgå den...
- der laves en ekstern/nodeport service som exposer kali-vnc pod
LETS GOOOO
-> gør det i yaml til at starte med...

DET FUCKING VIRKER !!! 


minikube dockerenv: eval $(minikube -p minikube docker-env)

