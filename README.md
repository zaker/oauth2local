# oauth2local

An oauth client providing authenticated tokens to local processes.

## Generate GRPC server/client

```bash

protoc -I ipc/localauth/ ipc/localauth/locauth.proto --go_out=plugins=grpc:ipc/localauth

```
