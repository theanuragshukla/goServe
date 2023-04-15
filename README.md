# goServe
#### GoServe is a platform independent CLI tool written using GO, which acts as a File Server.
--- 
# Features
- Access files over HTTP
- Expose local Files over LAN
- Options to change default behaviour
  - change host (default: localhost)
  - change port (default: 8080)
  - change directory (default: current directory)
  - expose (default: false)
# Flags
## -host / -l
This flag accepts a string value. This will override the default host for the file server.
```js
goServe -l 0.0.0.0
```
```js
goServe -host 0.0.0.0
```
## -port / -p
This flag accepts an Integer value, This will override the default port for the file server.
```js
goServe -p 8000
```
```js
goServe -port 8000
```
## -directory / -d
This flag accepts a string value, which will used as default directory for the file server. It only accepts the relative path for the directory.
```js
goServe -d ../../home
```
```js
goServe -directory ../../home
```
## -expose / -e
This flag accepts a string value, This flag if set to `true` exposes the file server to LAN.
```js
goServe -e
```
```js
goServe -expose
```
## -help / -h
This flag will show all the options available for this tool.
```js
goServe -h
```
```js
goServe -help
```
