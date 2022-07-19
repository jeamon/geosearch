# geosearch

This is a go-based project to demonstrate a clean design and implementation of a jobs search platform with google maps display features. It provides both a RestFul API backend and a common Web backend to allow browsing like on a real website. 

## Get-Started

* Clone the repository

```
$ git clone https://github.com/jeamon/geosearch.git
$ cd search-website
```

* Replace the Google Maps API Key by a valid one insde the config/config.yaml file.

```
 maps_api_key: "<put-your-google-maps-api-key-here>"
```

* Build the api server

```
$ docker build -t api-server -f Dockerfile.api .
```

* Run the api-server container

```
$ docker run -d -p 8095:8095 --name api-server --rm api-server:latest
```

* Check the api-server container status

```
$ docker ps | grep api-server
```

* Build the web server

```
$ docker build -t web-server -f Dockerfile.web .
```

* Run the web-server container

```
$ docker run -d -p 8090:8090 --name web-server --rm web-server:latest
```

* Check the web-server container status

```
$ docker ps | grep web-server
```

* View the container logs

```
$ docker logs <container-id>
```