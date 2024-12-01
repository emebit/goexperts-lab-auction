#GoExperts-Lab-Auction

### Executando o projeto

1. **Construir as Imagens**: Na raiz do projeto, execute:

   ```sh
   docker-compose build
   ```

2. **Executar o Docker Compose**: Na raiz do projeto (onde está o arquivo `docker-compose.yml`), execute:

   ```sh
   docker-compose up
   ```

3. **Acesse a aplicação:**

   A aplicação estará disponível em `http://localhost:8080`.

### Exemplos de Requisição

- **Para criar leilão:**

   ```sh
   curl -X POST http://localhost:8080/auction \
   -H "Content-Type: application/json" \
   -d '{
        "product_name": "",
        "category": "",
        "description": "",
        "condition": 1
       }'
   ```

- **Para buscar leilão por ID:**

   ```sh
   curl -X GET "http://localhost:8080/auction/{auctionId}" \
   -H "Content-Type: application/json"
   ```

- **Para buscar leilões por parâmetros de consulta:**

   ```sh
   curl -X GET "http://localhost:8080/auction?category=Movie&status=0&condition=1" \
   -H "Content-Type: application/json"
   ```
