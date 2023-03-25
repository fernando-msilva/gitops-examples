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
  --timeout=800s
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
```

## Instalando o ArgoCD

Crie um namespace chamado argocd e execute o deploy do ArgoCD com os seguintes comandos

```bash
$ kubectl create namespace argocd

$ kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

Instale o cli do ArgoCD com os comandosa abaixo:

```bash
$ curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64

$ sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd

$ rm argocd-linux-amd64
```

Crie o ingress para acessar a interface web do ArgoCD:

```bash
$ kubectl apply -f argocd
```