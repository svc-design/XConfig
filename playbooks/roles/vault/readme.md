# Init Vault Server

kubectl exec -t -i vault-0 -n vault -- sh 
vault operator init -key-shares=5 -key-threshold=3
