# Escolha uma imagem base. Aqui, usamos a versão oficial do Go.
FROM golang:1.18

# Defina o diretório de trabalho dentro do container
WORKDIR /app

# Copie o arquivo 'go.mod' e, se você tiver, o arquivo 'go.sum' para o diretório de trabalho
COPY go.mod ./
COPY go.sum ./

# Baixe as dependências do módulo Go.
RUN go mod download

# Copie o restante dos arquivos do código fonte para o diretório de trabalho
COPY . .

# Compile sua aplicação. Substitua 'main.go' pelo caminho para o seu arquivo principal, se for diferente.
RUN go build -o /bling_limit

# A porta em que sua aplicação ouve. Certifique-se de que é a mesma porta que você define no seu código.
EXPOSE 5000

# Comando para executar sua aplicação
CMD [ "/bling_limit" ]
