# Microsserviços gRPC - Deploy com Docker

Este projeto implementa uma arquitetura de microsserviços usando gRPC, Go e MySQL, com deploy completo via Docker.

## 🏗️ Arquitetura

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Order Service │    │ Payment Service │    │Shipping Service │
│     Port 3000   │◄──►│     Port 3001   │    │     Port 3002   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                       │                       │
          └───────────────────────┼───────────────────────┘
                                  │
                    ┌─────────────────┐
                    │  MySQL Database │
                    │     Port 3306   │
                    └─────────────────┘
```

## 🚀 Deploy Rápido

### Pré-requisitos

- Docker
- Docker Compose

### 1. Clonar o repositório

```bash
git clone <seu-repositorio>
cd pratica-grpc/microservices
```

### 2. Iniciar todos os serviços

```bash
docker-compose up -d
```

### 3. Verificar status dos serviços

```bash
docker-compose ps
```

### 4. Verificar logs

```bash
# Logs de todos os serviços
docker-compose logs -f

# Logs de um serviço específico
docker-compose logs -f order-service
```

## 🧪 Testando os Serviços

### Testar Order Service

```bash
# Pedido válido
grpcurl -d '{"customer_id": 123, "order_items": [{"product_code": "ABC123", "quantity": 2, "unit_price": 10.50}], "total_price": 21.00}' -plaintext localhost:3000 Order/Create

# Produto inexistente
grpcurl -d '{"customer_id": 123, "order_items": [{"product_code": "INVALID", "quantity": 1, "unit_price": 10.00}], "total_price": 10.00}' -plaintext localhost:3000 Order/Create

# Estoque insuficiente
grpcurl -d '{"customer_id": 123, "order_items": [{"product_code": "JKL111", "quantity": 50, "unit_price": 30.00}], "total_price": 1500.00}' -plaintext localhost:3000 Order/Create
```

## 📦 Serviços Disponíveis

| Serviço  | Porta | Descrição                   |
| -------- | ----- | --------------------------- |
| Order    | 3000  | Gerenciamento de pedidos    |
| Payment  | 3001  | Processamento de pagamentos |
| Shipping | 3002  | Cálculo de frete e envio    |
| MySQL    | 3306  | Banco de dados              |

## 🗄️ Banco de Dados

O projeto cria automaticamente 3 databases:

- `order` - Pedidos e produtos
- `payment` - Pagamentos
- `shipping` - Envios

### Produtos de exemplo (Order Service)

- ABC123 - Produto A (R$ 10,50, estoque: 100)
- XYZ789 - Produto B (R$ 20,00, estoque: 50)
- DEF456 - Produto C (R$ 15,75, estoque: 75)
- GHI999 - Produto D (R$ 5,25, estoque: 200)
- JKL111 - Produto E (R$ 30,00, estoque: 25)

## 📝 Variáveis de Ambiente

| Variável         | Descrição            | Padrão                                 |
| ---------------- | -------------------- | -------------------------------------- |
| DATA_SOURCE_URL  | URL de conexão MySQL | `root:minhasenha@tcp(mysql:3306)/[db]` |
| APPLICATION_PORT | Porta do serviço     | `3000/3001/3002`                       |
| ENV              | Ambiente de execução | `production`                           |

## 🎯 Features Implementadas

✅ **Microsserviços**: Arquitetura distribuída com gRPC  
✅ **Validação de Estoque**: Controle automático de produtos  
✅ **Containerização**: Deploy completo com Docker  
✅ **Orquestração**: Docker Compose para múltiplos serviços  
✅ **Health Checks**: Monitoramento de saúde dos serviços  
✅ **Persistência**: Volumes MySQL para dados persistentes  
✅ **Network**: Rede isolada para comunicação entre serviços
