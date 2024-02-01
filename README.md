## Presentation
### Home view
[![homepage2.png](https://i.postimg.cc/KzSHS5NB/homepage2.png)](https://postimg.cc/sBmc57bg)

Presented with a stream of photos, uploaded by the own users followed. On each photo, the user logged can put a like, leave a comment (and later delete), or view all comments. At the top of the page, the user can:
* Search other users with the searchBar
* Upload a photo
* Visit own profile

### User view
![alt text](https://i.postimg.cc/fLjrf9Ns/Profilo-marta.png)

Presented with a stream of photos, uploaded by the user visited. On each photo, the user logged can do the default operations. At the top of the page, the user can:
* Follow the user
* Ban the user

### Profile view
[![profilepage.png](https://i.postimg.cc/vHf7CdKq/profilepage.png)](https://postimg.cc/cr188P0w)

Presented with a stream of photos, uploaded by the user logged. On each photo, the user logged can do the default operations and also remove each comment or remove the photo. At the top of the page, the user can:
* Perform logout
* Change username

## Project structure

* `cmd/` contains all executables; Go programs here should only do "executable-stuff", like reading options from the CLI/env, etc.
	* `cmd/healthcheck` is an example of a daemon for checking the health of servers daemons; useful when the hypervisor is not providing HTTP readiness/liveness probes (e.g., Docker engine)
	* `cmd/webapi` contains an example of a web API server daemon
* `demo/` contains a demo config file
* `doc/` contains the documentation (usually, for APIs, this means an OpenAPI file)
* `service/` has all packages for implementing project-specific functionalities
	* `service/api` contains an example of an API server
	* `service/globaltime` contains a wrapper package for `time.Time` (useful in unit testing)
* `vendor/` is managed by Go, and contains a copy of all dependencies
* `webui/` is an example of a web frontend in Vue.js; it includes:
	* Bootstrap JavaScript framework
	* a customized version of "Bootstrap dashboard" template
	* feather icons as SVG
	* Go code for release embedding

Other project files include:
* `open-npm.sh` starts a new (temporary) container using `node:lts` image for safe web frontend development (you don't want to use `npm` in your system, do you?)

## Node/NPM vendoring

This repository contains the `webui/node_modules` directory with all dependencies for Vue.JS. You should commit the content of that directory and both `package.json` and `package-lock.json`.

## How to build

If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:

```shell
go build ./cmd/webapi/
```

If you're using the WebUI and you want to embed it into the final executable:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run build-embed
exit
# (outside the NPM container)
go build -tags webui ./cmd/webapi/
```

## How to run (in development mode)

You can launch the backend only using:

```shell
go run ./cmd/webapi/
```

If you want to launch the WebUI, open a new tab and launch:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run dev
```

## For running the preview, open a new tab and launch (after running the backend):

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run build-prod
npm run preview
```
## Deployment
### How to build the images 
Backend
```
$ docker build -t wasa-photos-backend:latest -f Dockerfile.backend .
```
Frontend 
```
$ docker build -t wasa-photos-frontend:latest -f Dockerfile.frontend .
```
### How to run the container images
Backend
```
$ docker run -it --rm -p 3000:3000 wasa-photos-backend:latest
```
Frontend
```
$ docker run -it --rm -p 8080:80 wasa-photos-frontend:latest
```
### License

See [LICENSE](LICENSE).
