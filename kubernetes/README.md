SSH Keys secret is generated outside the kubernetes lifecycle.

```bash
kubectl create secret generic bot-ssh-key --from-file=id_rsa=/path/to/.ssh/id_rsa
```
Using the above command. The name of it is `bot-ssh-key`.