# Auth Service

Auth Service includes 3 docker containers:
- auth - server processing authentication requests.
- postgres - permanent storage of data.
- migrator - performs migration in database using goose package.

## Deploy

To deploy Auth Service:
```
# make docker-deploy ENV=<environment>
```
*ENV is used then as a config name. Possible ENV values are now `stage` and `prod`.*

To stop Auth Service:
```
# make docker-stop ENV=<environment>
```
