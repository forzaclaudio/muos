# Media Uploader to Object Storage (MUOS)

This webservice receives a path to a file and uploads it to an object storage in the cloud.

## To post an upload taks use the following:
```bash
curl -X POST X.X.X.X:1234/upload -d '{"win_filepath": "/path/to/a/file"}'  -H 'Content-Type: application/json'
```

## To list the ports currently forwarded between WSL and Windows
```bash
netsh interface portproxy show v4tov4
```

## To forward ports between WSL and Windows
```bash
netsh interface portproxy set v4tov4 listenport=8080 listenaddress=0.0.0.0 connectport=8080 connectaddress=$(wsl hostname -I)
```
