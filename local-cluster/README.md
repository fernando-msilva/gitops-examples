## Iniciando o ambiente

Antes de tudo adicione esses dois registros no seu arquivo `/etc/hosts`

```
127.0.0.1 argocd.local
127.0.0.1 hello.local
```

Execute o seguinte comando para subir o cluster

```bash
$ kind create cluster --config kind.yaml
```

Para instalar o ingress controller no cluster execute o comando abaixo:

```bash
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

O ingress levará um certo tempo para ser totalmente configurado. Para validar se o controller foi devidamente criado execute o seguinte comando:

```bash
$ kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=1500s
```

Crie a imagem da api de exemplo com o comando abaixo:

```bash
$ docker build -t hello ./application
```

Para carregar a imagem local dentro do cluster você pode executar o seguinte comando:

```bash
$ kind load docker-image --name gitops-local-lab hello
```

Faça o deploy da aplicação com o seguinte comando:

```bash
$ kubectl apply -f deployments
```

Caso o ingress ainda não esteja disponível você verá a seguinte mensagem de erro:

```
Error from server (InternalError): error when creating "deployments/hello-app-deployment.yaml": Internal error occurred: failed calling webhook "validate.nginx.ingress.kubernetes.io": failed to call webhook: Post "https://ingress-nginx-controller-admission.ingress-nginx.svc:443/networking/v1/ingresses?timeout=10s": dial tcp 10.96.233.212:443: connect: connection refused
```

Teste a aplicação acessando a seguinte URL:

```bash
$ curl http://hello.local/world
{"hello":"world"}%

$ curl http://hello.local/test
{"hello":"test"}%

$ curl http://hello.local/there
{"hello":"there"}%
```

## Instalando o ArgoCD

Crie um namespace chamado argocd e execute o deploy do ArgoCD com os seguintes comandos

```bash
$ kubectl create namespace argocd

$ helm install --namespace argocd --values ./argocd/values.yaml argocd ./argocd
```

Tente acessar o painel web do ArgoCD através do endereço `http://argocd.local`. É possível que o retorno seja

O usuário padrão da instalação é `admin`. Para recuperar a senha padrão execute o seguinte comando:

```bash
$ kubectl get secret argocd-initial-admin-secret -n argocd -o jsonpath="{.data.password}" | base64 -d
```

Instale o cli do ArgoCD com os comandosa abaixo:

```bash
$ curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64

$ sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd

$ rm argocd-linux-amd64
```

Autentique o ArgoCD via cli com o comando abaixo:

```bash
$ argocd login argocd.local
```

Agora execute o comando abaixo para recuperar o contexto do cluster e adicina-lo na configuração do cli do argocd:

```bash
$ kubectl config current-context
O_NOME_DO_SEU_CONTEXT

$ argocd cluster add O_NOME_DO_SEU_CONTEXT
WARNING: This will create a service account `argocd-manager` on the cluster referenced by context `O_NOME_DO_SEU_CONTEXT` with full cluster level privileges. Do you want to continue [y/N]? y
INFO[0010] ServiceAccount "argocd-manager" created in namespace "kube-system" 
INFO[0010] ClusterRole "argocd-manager-role" created    
INFO[0010] ClusterRoleBinding "argocd-manager-role-binding" created 
INFO[0015] Created bearer token secret for ServiceAccount "argocd-manager" 
WARN[0016] Failed to invoke grpc call. Use flag --grpc-web in grpc calls. To avoid this warning message, use flag --grpc-web. 
FATA[0016] rpc error: code = Unknown desc = Get "https://127.0.0.1:36803/version?timeout=32s": dial tcp 127.0.0.1:36803: connect: connection refused 
```

Para contornar esse erro execute o comando `kubectl get -n default endpoints`. A saída será algo parecido com isso:

```bash
NAME         ENDPOINTS         AGE
kubernetes   172.18.0.2:6443   103m
```

Agora copie o ip e porta que foi mostrado com a execução do comando anterior e altere somente o valor de endereço do server no seu arquivo `.kube/config`, como no exemplo abaixo onde o ip antigo foi comentado e o novo endereço foi configurado:

```yaml
apiVersion: v1
clusters:
- cluster:
    #server: https://127.0.0.1:32919
    server: https://172.18.0.2:6443
  name: kind-kind
```


Após essa modificação execute novamente o comando para adicionar o cluster ao ArgoCD

```bash
argocd cluster add O_NOME_DO_SEU_CONTEXT
```