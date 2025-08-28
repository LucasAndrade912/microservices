# Microsserviço Order - Validação de Estoque

## Como Usar

### 1. Iniciar o serviço
```bash
cd microservices/order
DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3306)/order" APPLICATION_PORT=3000 ENV=development go run ./cmd
```
> **Nota**: O `parseTime=true` é adicionado automaticamente pelo código, não é necessário incluir na URL.

### 2. Testar um pedido válido
```bash
grpcurl -d '{"customer_id": 123, "order_items": [{"product_code": "ABC123", "quantity": 2, "unit_price": 10.50}], "total_price": 21.00}' -plaintext localhost:3000 Order/Create
```

### 3. Testar produto inexistente
```bash
grpcurl -d '{"customer_id": 123, "order_items": [{"product_code": "INVALID", "quantity": 1, "unit_price": 10.00}], "total_price": 10.00}' -plaintext localhost:3000 Order/Create
```

### 4. Testar estoque insuficiente
```bash
grpcurl -d '{"customer_id": 123, "order_items": [{"product_code": "JKL111", "quantity": 50, "unit_price": 30.00}], "total_price": 1500.00}' -plaintext localhost:3000 Order/Create
```
> **Nota**: O produto JKL111 tem apenas 25 unidades em estoque, então pedindo 50 deve retornar erro.

## Funcionalidades

✅ **Validação de Produtos**: Verifica se produtos existem e estão ativos  
✅ **Controle de Estoque**: Atualiza estoque automaticamente  
✅ **Tratamento de Erros**: Retorna erros apropriados via gRPC  
✅ **Produtos de Exemplo**: 5 produtos criados automaticamente  

## Produtos Disponíveis

- ABC123 - Produto A (R$ 10,50, estoque: 100)
- XYZ789 - Produto B (R$ 20,00, estoque: 50)  
- DEF456 - Produto C (R$ 15,75, estoque: 75)
- GHI999 - Produto D (R$ 5,25, estoque: 200)
- JKL111 - Produto E (R$ 30,00, estoque: 25)
