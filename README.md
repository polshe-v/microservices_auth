# Auth Service

Auth Service includes 3 docker containers:
- auth - server processing authentication requests.
- postgres - permanent storage of data.
- migrator - performs migration in database using goose package.

## Deploy

Make sure docker network `service-net` is in place for microservices communication. If none exists, then create network:
```
# make docker-net
```

To deploy Auth Service:
```
# make docker-deploy ENV=<environment>
```
*ENV is used then as a config name. Possible ENV values are now `stage` and `prod` as these configs are now in the repository.*

To stop Auth Service:
```
# make docker-stop ENV=<environment>
```
