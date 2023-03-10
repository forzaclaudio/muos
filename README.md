# Media Uploader to Object Storage (MUOS)

This webservice receives a path to a file and uploads it to an object storage in the cloud.

To post an upload taks use the following:
```bash
curl -X POST X.X.X.X:1234/upload -d '{"win_filepath": "/path/to/a/file"}'  -H 'Content-Type: application/json'
```
