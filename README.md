# product


connector settings go!
1. docker exec -it `name_container` psql -U admin -d product -c "SHOW wal_level;"
2. ALTER SYSTEM SET wal_level = logical;
   SELECT pg_reload_conf();
3. relog database
4. Invoke-RestMethod -Uri "http://localhost:8083/connectors" `
  -Method Post `
   -Headers @{"Content-Type"="application/json"} `
   -Body (Get-Content (Resolve-Path ".\connector.json") -Raw)
5. check it : Invoke-RestMethod -Uri "http://localhost:8083/connectors/product-connector/status"
6. check status : Invoke-RestMethod -Uri "http://localhost:8083/connectors/product-connector/status"



