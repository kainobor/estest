## Description

Create a REST API Search Service

Let’s say our company has e-commerce platform with selling fashion products. Each product
has title, brand, price and stock. The platform should have search functionality.
Your task is to implement microservice to serve search requests.

The Service has to support:
- GET method to perform search queries.
E.g. https://example.com/products?q=black shoes
- Authentication
- Api should have versioning
- Pagination and sorting
- Filtering. E.g. https://example.com/products?q=black shoes&filter=brand:brand_name
Requirements:
- Language: Go
- Storage: ElasticSearch
- Service should be dockerized

## How to start

Start elasticsearch docker:

`docker-compose up`
`docker exec -ti es sh`
`сd /migrations`
`curl -H 'Content-Type: application/json' -XPUT 'localhost:9200/user' --data-binary @create_user`
`curl -H 'Content-Type: application/x-ndjson' -XPOST 'localhost:9200/product/product/_bulk?pretty' --data-binary @insert_product`
`curl -H 'Content-Type: application/x-ndjson' -XPOST 'localhost:9200/user/_doc/_bulk?pretty' --data-binary @insert_user`

Check that:

`curl -XGET "localhost:9200/product/product/_search?pretty" -d ''`