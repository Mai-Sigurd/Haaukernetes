# Connecting via Wireguard

You must have a public/private key pair for Wireguard. 
This can be generated via the following command (on Unix-based systems, as per the official Wireguard documentation: https://www.wireguard.com/quickstart/)

``wg genkey | tee privatekey | wg pubkey > publickey``

You must also have Wireguard installed on your system.

You will provide your public key to Haaukernetes and in turn recieve a configuration file that you must save and insert your private key into.

The configuration file will look like the following, though with the PLACEHOLDER values substituted for real values.
You must insert your private key in place for the PRIVATEKEY value.

    [Interface]
    Address = 10.33.0.2/32
    PrivateKey = PRIVATEKEY
    DNS = 10.96.0.10
    
    [Peer]
    PublicKey = PLACEHOLDER
    
    Endpoint = PLACEHOLDER
    AllowedIPs = 10.96.0.0/12
    PersistentKeepalive = 25

Now you can activate the connection by executing the command (on Unix-based systems)

``sudo wg-quick up ./PATH-TO-CONFIGURATION-FILE``

Or consult the official website for instructions on how to do it on your specific platform.