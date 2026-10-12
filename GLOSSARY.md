# go-db

Connection management for the datastores an ooqls service talks to (Postgres, Redis, Valkey, Elasticsearch). Turns registry configuration into live, reconnecting client handles.

## Language

**Dial Config**:
The resolved, dial-ready shape of a datastore connection: host, port, credentials, target database, and TLS. Produced once from a registry entry and consumed by every datastore adapter (redis, valkey, pgx, sqlx, elasticsearch) to open its connection.
_Avoid_: Options, connection options, registry options
