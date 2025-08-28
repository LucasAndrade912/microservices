-- Script para popular a tabela de produtos
-- Execute este script no banco 'order' para adicionar produtos de exemplo

USE `order`;

-- Inserir produtos de exemplo se não existirem
INSERT INTO products (product_code, name, unit_price, stock, active, created_at, updated_at) 
VALUES 
    ('ABC123', 'Produto A - Eletrônico', 10.50, 100, 1, NOW(), NOW()),
    ('XYZ789', 'Produto B - Livro', 20.00, 50, 1, NOW(), NOW()),
    ('DEF456', 'Produto C - Roupa', 15.75, 75, 1, NOW(), NOW()),
    ('GHI999', 'Produto D - Acessório', 5.25, 200, 1, NOW(), NOW()),
    ('JKL111', 'Produto E - Decoração', 30.00, 25, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE 
    name = VALUES(name),
    unit_price = VALUES(unit_price),
    updated_at = NOW();

-- Verificar produtos inseridos
SELECT * FROM products WHERE active = 1;
