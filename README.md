# sympton-checker

```shell
curl "http://localhost:8081/symptoms"
```

```shell
curl -X POST --location "http://localhost:8081/symptoms" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer whatever" \
    -d  "[\"HP:0000256\", \"HP:0001249\"]"
```