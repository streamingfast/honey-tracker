# honey-tracker

### Build command
```bash
go mod tidy
go install ./cmd/honey-tracker
```

### Run postgres from docker
```bash
docker run --rm --name postgres -p 5432:5432 -e POSTGRES_USER=<replace_with_your_user_of_choice> -e POSTGRES_PASSWORD=secureme -d postgres
```

### Create schema
```bash
docker ps -a # to see the id of your pod
docker exec -it <id_of_pod> bash
psql -d postgres -U <replace_with_your_user_of_choice> -W
CREATE SCHEMA IF NOT EXISTS hivemapper;
```

### Run honey-tracker
```bash
sftoken # to set SUBSTREAMS_API_TOKEN
honey-tracker mainnet.sol.streamingfast.io:443 https://github.com/streamingfast/substreams-hivemapper/releases/download/v0.1.0/hivemapper-v0.1.0.spkg map_outputs --db-host=localhost --db-port=5432 --db-user=eduard --db-password=secureme --db-name=postgres --output-module-type=proto:hivemapper.types.v1.Output
```
