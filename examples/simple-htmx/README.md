## Usage

```
# Start backend
cd $ROOT_FOLDER
go run . \
  --addr localhost:8080 \
  --allowed-root-account-handle zerokspot.com
  --allowed-origin http://localhost:9980

# In a separate terminal, start caddy:
caddy file-server . \
  --listen localhost:9980
```
