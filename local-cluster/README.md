## Iniciando o ambiente

Execute o seguinte comando para subir o cluster

```bash
$ kind create cluster --config kind.yaml
```

Para instalar o ingress controller no cluster execute o comando abaixo:

```bash
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
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

Teste a aplicação acessando a seguinte URL:

```bash
$ curl http://localhost/world
```