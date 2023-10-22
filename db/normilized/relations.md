# Таблицы
## Product - Таблица с полями объявления

|  Название       |  Тип           |
|:----------------|:---------------|
| $\underline{\text{id}}$             | bigserial      |
| saler_id        | bigint         |
| title           | char[256]      |
| description     | text           |
| price           | bigint         |
| category_id     | bigint         |
| creation_date   | timestamp     |
| views           | int |
| in_favourites   | int |
| available_count | int |
| city            | char[256]      |
| delivery        | boolean        |
| safe_dial       | boolean        |  
### ФЗ
{ id } $\rightarrow$ saler_id, category_id, title, description, price, creation_date, views, in_favourites, available_count, city, delivery, safe_dial
### Нормальные формы
1. все типы атрибутов атомарны $\Rightarrow$ отношение находится в 1НФ.
2. отношение находится в 1НФ, ключ является не составным $\Rightarrow$ не имеем зависимости неключевых атрибутов от части ключа $\Rightarrow$ отношение находится в 2НФ.
3. отношение находится в 2НФ и мы не имеем транзитивных ФЗ неключевых атрибутов от ключевых $\Rightarrow$ отношение находится в 3НФ.
4. отношение находится в 3НФ и мы не имеем составного ключа $\Rightarrow$ отношение находится в НФБК.
### Ограничения
id primary key
saler_id foreign key
## User

|Название|Тип|
|---|---|
|$\underline{\text{id}}$|bigserial|
|email|char[256]|
|phone|char[18]|
|name|char[256]|
|pass|char[256]|
|birthday|timestamp|
### ФЗ
{ id } $\rightarrow$ email, phone, name, pass, birthday
### Нормальные формы
1. все типы атрибутов атомарны $\Rightarrow$ отношение находится в 1НФ.
2. отношение находится в 1НФ, ключ является не составным $\Rightarrow$ не имеем зависимости неключевых атрибутов от части ключа $\Rightarrow$ отношение находится в 2НФ.
3. отношение находится в 2НФ и мы не имеем транзитивных ФЗ неключевых атрибутов от ключевых $\Rightarrow$ отношение находится в 3НФ.
4. отношение находится в 3НФ и мы не имеем составного ключа $\Rightarrow$ отношение находится в НФБК.

## Order
|Название|Тип|
|---|---|
|$\underline{\text{id}}$|bigserial|
|owner_id|bigint|
|product_id|bigint|
|count|smallint|
|status|smallint|
|create_date|timestamp|
|update_date|timestamp|
|close_date|timestamp|
### ФЗ
{ id } $\rightarrow$ owner_id, product_id, count, status, creation_data, update_data, closed_data
### Нормальные формы
1. все типы атрибутов атомарны $\Rightarrow$ отношение находится в 1НФ.
2. отношение находится в 1НФ, ключ является не составным $\Rightarrow$ не имеем зависимости неключевых атрибутов от части ключа $\Rightarrow$ отношение находится в 2НФ.
3. отношение находится в 2НФ и мы не имеем транзитивных ФЗ неключевых атрибутов от ключевых $\Rightarrow$ отношение находится в 3НФ.
4. отношение находится в 3НФ и мы не имеем составного ключа $\Rightarrow$ отношение находится в НФБК.
## Image
|Название|Тип|
|---|---|
|$\underline{\text{id}}$|bigserial|
|url|char[256]|
|product_id|bigint|
### ФЗ
{ id } $\rightarrow$ url, product_id
### Нормальные формы
1. все типы атрибутов атомарны $\Rightarrow$ отношение находится в 1НФ.
2. отношение находится в 1НФ, ключ является не составным $\Rightarrow$ не имеем зависимости неключевых атрибутов от части ключа $\Rightarrow$ отношение находится в 2НФ.
3. отношение находится в 2НФ и мы не имеем транзитивных ФЗ неключевых атрибутов от ключевых $\Rightarrow$ отношение находится в 3НФ.
4. отношение находится в 3НФ и мы не имеем составного ключа $\Rightarrow$ отношение находится в НФБК.
## Category
|Название|Тип|
|---|---|
|$\underline{\text{id}}$|bigserial|
|name|char[256]|
|parent_id|bigint|
### ФЗ
{ id } $\rightarrow$ name, parent_id
### Нормальные формы
1. все типы атрибутов атомарны $\Rightarrow$ отношение находится в 1НФ.
2. отношение находится в 1НФ, ключ является не составным $\Rightarrow$ не имеем зависимости неключевых атрибутов от части ключа $\Rightarrow$ отношение находится в 2НФ.
3. отношение находится в 2НФ и мы не имеем транзитивных ФЗ неключевых атрибутов от ключевых $\Rightarrow$ отношение находится в 3НФ.
4. отношение находится в 3НФ и мы не имеем составного ключа $\Rightarrow$ отношение находится в НФБК.
## Favourite
|Название|Тип|
|---|---|
|$\underline{\text{id}}$|bigserial|
|product_id|bigint|
|owner_id|bigint|
### ФЗ
{ id } $\rightarrow$ product_id, owner_id
### Нормальные формы
1. все типы атрибутов атомарны $\Rightarrow$ отношение находится в 1НФ.
2. отношение находится в 1НФ, ключ является не составным $\Rightarrow$ не имеем зависимости неключевых атрибутов от части ключа $\Rightarrow$ отношение находится в 2НФ.
3. отношение находится в 2НФ и мы не имеем транзитивных ФЗ неключевых атрибутов от ключевых $\Rightarrow$ отношение находится в 3НФ.
4. отношение находится в 3НФ и мы не имеем составного ключа $\Rightarrow$ отношение находится в НФБК.
***
# ERDiagram
Код **mermaid** диаграммы представлен ниже. ![Визуализация](https://www.mermaidchart.com/raw/bb85db47-1f82-454d-920a-3b15c2000041?version=v0.1&theme=dark&format=svg)
```postgresql
---
title: ERDiagram "Работяги"
---
erDiagram
    User ||--o{ Favourite : add_to
    User ||--o{ Product : create
    User ||--o{ Order : places
    User {
        bigserial id
        char[256] email
        char[256] name
        char[256] pass
        char[18] phone
        timestamp birthday
    }

    Order {
        bigserial id
        bigint owner_id
        bigint product_id
        smallint count
        smallint status
        timestamp creation_date
        timestamp update_data
        timestamp closed_date
    }

    Favourite }|--|| Product : contains
    Favourite {
        bigserial id
        bigint product_id
        bigint owner_id
    }

    Product ||--|{ Image : contains
    Product |o--|| Category : associated_with
    Product {
        bigserial id
        bigint saler_id
        char[256] title
        text description
        bigint price
        bigint category_id
        timestamp creation_date
        int views
        int in_favourites
        int available_count
        char[256] city
        boolean delivery
        boolean safe_dial
    }

    Image {
        bigserial id
        bigint product_id
        char[256] url
    }

    Category {
        bigserial id
        bigint parent_id
        char[256] name
    }
```
