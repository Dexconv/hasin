## How to run

you can simply run this project using the provided docker compose:

```bash
mkdir postgres/data/ store/files/
docker compose up --build
```

this will set up the database and run both services. then you can proceed with using the services as neede, the postman collection is provided to simplify the usage of the api. the requests for both services are provided but you only need to intract with the "retrieve" service. this service handles user management and access controll to the "store" service api. simply log in or register with the retrieve service and use the jwt token provided to authenticate your call to the upload and download endpoints.

1. from hasin collection, in the retrieve folder use the register, or login endpoint to recieve the jwt token, edit the json body as needed.
1. save the new toke to the value of "jwtToken" variable, in the collection variables tab.
1. navigate to the "upload file" endpoint in the retrieve folder, edit the form data as neede, enter tags seprated with "," and select a file smaller than the 1MB limit, upload how many files you want.
1. navigate to the "download file" endpoint in the retrieve folder, edit the form data as neede, choose selection mode from "name" and "tags", and enter the tags field value accordingly (single name, single tag, multiple tags).

you will recieve a zip file, containing all the files matching the query. here's how the files are chosen:

1. if the full name matches the provided name (in name mode)
1. if it has all the tags provided (in tags mode)
1. the first file uploaded, if non of the above matches any files

## project specifications

- docker compose, custom network and services
- postgres, seprate db for each service, persistent data
- set 1M body limit on upload request (file upload limit)
- used structured log (log/slog) with added source info
- takes tags from form data and filename from uploaded file
- extracts MimeType of the file + extention from the body of uploaded file
- serves one or more files, if any exists, in a zip archive
- simple token based authentication using jwt
- proxies api calls from retrieve to store service as requested, with added authentication
- files are encrypted before storing in filesystem with a random name
- decrypt files and return them with origional name
- the stored files overall size is checked to ensure storage doesn't go over the set limit (default 50MB)
- the config file is located in each services path "config/file/config.json" and will be auto generated using default values if it doesn't exist
