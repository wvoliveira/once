# Once

> **Compartilhamento de arquivos seguro, efêmero e auto-hospedado.**
> Visualize uma vez, e o arquivo desaparece para sempre.

Um sistema de *Burn After Reading* (queimar após ler) construído em **Golang**. Projetado para ser extremamente simples de implantar (binário único), mas robusto o suficiente para produção graças à persistência em SQLite e armazenamento local.


## Destaques

* **Zero dependências externas:** Não precisa de Redis, S3 ou Docker Compose complexo. Apenas o binário.
* **Persistência real:** Utiliza **SQLite** para manter o estado dos arquivos. Se o servidor reiniciar, nada é perdido.
* **Atomicidade garantida:** Usa transações SQL (`DELETE ... RETURNING`) para garantir que **apenas uma pessoa** consiga visualizar o arquivo, eliminando condições de corrida (*race conditions*).
* **Coletor de lixo automático:** Um processo em background limpa arquivos expirados do disco automaticamente.
* **Driver pure Go:** Utiliza `modernc.org/sqlite`, permitindo compilação cruzada (CGO-free) fácil para qualquer sistema operacional.


## Arquitetura

O sistema adota uma abordagem minimalista de Clean Architecture:

1.  **Storage (Disco):** Os arquivos cifrados/brutos são salvos na pasta `./uploads`.
2.  **Metadata (SQLite):** Dados como ID, Nome, MimeType e Expiração ficam no `app.db`.
3.  **Atomic Lock:** A leitura utiliza uma instrução atômica de banco de dados para buscar e deletar o registro simultaneamente. Se dois usuários tentarem baixar ao mesmo tempo, o banco de dados garante que apenas um receba o arquivo.


## Como Rodar

### Pré-requisitos
* Go 1.22+

### Instalação

1.  Clone o repositório:
    ```bash
    git clone git@github.com:wvoliveira/once.git
    cd once
    ```

2.  Instale as dependências:
    ```bash
    go mod tidy
    ```

3.  Execute o projeto:
    ```bash
    go run cmd/once/main.go
    ```
    *Isso criará automaticamente a pasta `./uploads` e o arquivo `app.db`.*

### Build

Como usamos o driver *Pure Go* do SQLite, o build é simples e estático:

```bash
# Exemplo para Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .
```


## Uso (exemplos de Integração)

Nota: O código atual implementa a lógica Core. Abaixo exemplos de como seria o uso caso exponha via HTTP.

1. Upload de Arquivo
Envia um arquivo que expira em 1 hora ou após a primeira visualização.

```bash
# Exemplo de chamada interna
sys.Upload("id-unico-123", "contrato.pdf", conteudoBytes, 1 * time.Hour)
```

2. Visualização Única
Tenta recuperar o arquivo. Se funcionar, o arquivo é deletado instantaneamente do sistema.

```bash
# Exemplo de chamada interna
data, err := sys.ViewOneTime("id-unico-123")
if err != nil {
    fmt.Println("Arquivo não encontrado ou já lido!")
}
```


## Configuração
As configurações são definidas por constantes no main.go (podem ser refatoradas para Variáveis de Ambiente):

|Variável|Padrão|Descrição|
|--|--|--|
|ONCE_DATABASE_URI | file:app.db?_journal_mode=WAL&_cache_size=2000&_foreign_keys=on&_busy_timeout=5000&_synchronous=NORMAL | Onde os arquivos físicos serão salvos.|
|ONCE_STORAGE_FOLDER | uploads | Pasta onde os arquivos serão armazenados.|
|ONCE_STORAGE_CLEANER_INTERVAL | 60 segundos | Intervalo para verificar se há arquivos para serem deletados (expirados).|
|ONCE_HTTP_SERVER_PORT | 8080 | A porta do web server.|
|ONCE_HTTP_SERVER_READ_TIMEOUT | 3 segundos | Tempo para ler a requisição.|
|ONCE_HTTP_SERVER_WRITE_TIMEOUT | 3 segundos | Tempo para responder a requisição.|
|ONCE_HTTP_SERVER_GRACEFUL_SHUTDOWN | 10 segundos | Tempo para o graceful shutdown do web server.|


## Por que SQLite e Disco Local?

Para um MVP, a complexidade de infraestrutura é o maior inimigo.

Custo: Elimina custos de AWS S3 e instâncias gerenciadas de Redis.

Latência: Ler do disco SSD local e do SQLite (em processo) é ordens de grandeza mais rápido que requisições de rede para serviços de nuvem.

Simplicidade: O backup do seu sistema inteiro consiste em copiar a pasta ./uploads e o arquivo app.db.


## Licença

Este projeto está sob a licença MIT. Sinta-se livre para usar e modificar.