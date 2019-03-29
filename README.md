# oauth2local

An oauth client for bridging local scripts.

## Generate GRPC server/client

```bash

protoc -I ipc/localauth/ ipc/localauth/locauth.proto --go_out=plugins=grpc:ipc/localauth

```
