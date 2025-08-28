# Microsserviços gRPC - Deploy com Docker

Este projeto implementa uma arquitetura de microsserviços usando gRPC, Go e MySQL, totalmente containerizada com Docker.

## 🏗️ Arquitetura

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Order     │    │   Payment   │    │  Shipping   │
│  Service    │    │   Service   │    │   Service   │
│  :3000      │    │   :3001     │    │   :3002     │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
               ┌─────────────┐
               │    MySQL    │
               │    :3306    │
               └─────────────┘
```

## 🚀 Quick Start

### 1. Pré-requisitos
- Docker
- Docker Compose
- Make (opcional, mas recomendado)

### 2. Executar com Docker Compose

```bash
# Clonar o repositório
cd microservices

# Iniciar todos os serviços
docker-compose up -d

# Ou usando make
make up
```

### 3. Verificar se os serviços estão rodando

```bash
# Ver logs
docker-compose logs -f

# Ou usando make
make logs
```

### 4. Testar os serviços

```bash
# Criar um pedido (Order Service)
grpcurl -d '{
  "customer_id": 123, 
  "order_items": [
    {"product_code": "ABC123", "quantity": 2, "unit_price": 10.50}
  ], 
  "total_price": 21.00
}' -plaintext localhost:3000 Order/Create

# Testar produto inexistente
grpcurl -d '{
  "customer_id": 123, 
  "order_items": [
    {"product_code": "INVALID", "quantity": 1, "unit_price": 10.00}
  ], 
  "total_price": 10.00
}' -plaintext localhost:3000 Order/Create
```

## 📋 Comandos Úteis

### Com Make
```bash
make help          # Ver todos os comandos
make build         # Construir imagens
make up            # Iniciar serviços
make down          # Parar serviços
make logs          # Ver logs
make clean         # Limpar volumes
make test          # Testar conectividade
```

### Com Docker Compose
```bash
docker-compose build                    # Construir imagens
docker-compose up -d                    # Iniciar em background
docker-compose down                     # Parar serviços
docker-compose logs -f [service]        # Ver logs
docker-compose restart [service]        # Reiniciar serviço
```

## 🗄️ Banco de Dados

### Conexão
- **Host**: localhost:3306
- **Usuário**: root
- **Senha**: minhasenha
- **Databases**: order, payment, shipping

### Produtos Pré-cadastrados
- ABC123 - Produto A (R$ 10,50, estoque: 100)
- XYZ789 - Produto B (R$ 20,00, estoque: 50)
- DEF456 - Produto C (R$ 15,75, estoque: 75)
- GHI999 - Produto D (R$ 5,25, estoque: 200)
- JKL111 - Produto E (R$ 30,00, estoque: 25)

## 🏃‍♂️ Executar Individual

### Order Service
```bash
cd order
docker build -t order-service .
docker run -p 3000:3000 \
  -e DATA_SOURCE_URL="root:minhasenha@tcp(mysql:3306)/order" \
  -e APPLICATION_PORT=3000 \
  -e PAYMENT_SERVICE_URL=payment:3001 \
  -e SHIPPING_SERVICE_URL=shipping:3002 \
  order-service
```

### Payment Service
```bash
cd payment
docker build -t payment-service .
docker run -p 3001:3001 \
  -e DATA_SOURCE_URL="root:minhasenha@tcp(mysql:3306)/payment" \
  -e APPLICATION_PORT=3001 \
  payment-service
```

### Shipping Service
```bash
cd shipping
docker build -t shipping-service .
docker run -p 3002:3002 \
  -e DATA_SOURCE_URL="root:minhasenha@tcp(mysql:3306)/shipping" \
  -e APPLICATION_PORT=3002 \
  shipping-service
```

## 🔧 Troubleshooting

### Logs dos serviços
```bash
docker-compose logs order
docker-compose logs payment
docker-compose logs shipping
docker-compose logs mysql
```

### Reiniciar um serviço específico
```bash
docker-compose restart order
```

### Limpar tudo e recomeçar
```bash
make clean
make build
make up
```

### Conectar ao MySQL
```bash
docker exec -it microservices-mysql mysql -u root -pminhasenha
```

## 📁 Estrutura do Projeto

```
microservices/
├── docker-compose.yml          # Orquestração dos serviços
├── Makefile                    # Comandos automatizados
├── init.sql                    # Script de inicialização do DB
├── order/
│   ├── Dockerfile
│   ├── .dockerignore
│   └── ...
├── payment/
│   ├── Dockerfile
│   ├── .dockerignore
│   └── ...
└── shipping/
    ├── Dockerfile
    ├── .dockerignore
    └── ...
```

## 🚨 Portas dos Serviços

- **Order Service**: 3000
- **Payment Service**: 3001  
- **Shipping Service**: 3002
- **MySQL**: 3306

## 📝 Notas

- Os serviços aguardam o MySQL estar saudável antes de iniciar
- Volumes persistentes para dados do MySQL
- Network isolada para comunicação entre containers
- Health checks configurados
- Restart automático em caso de falha
