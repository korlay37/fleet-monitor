
# Coding Challenge: Fleet Monitor
### Clone repo
```bash
    git clone <repo_url>
```
## Run locally
### Run to install dependencies
```bash
    go mod download
```
### Set up environment variables
The program uses DEVICES_FILE and PORT env variables, if not found will default to "devices.csv" and 6733.

### Build(Creates executable)
```bash
    go build .\cmd\api\main.go
```
Run executable file.

### Run directly
API will run on localhost and the PORT sent in env or 6733.
```bash
    go run .\cmd\api\main.go

```
## Run container (must have docker installed)
#### Requires .env file with variables DEVICES_FILE and PORT
Docker will build the "fleet-monitor" image and run a "fleet-monitor" container
```bash
    docker compose up --build -d
```



## Solution description

### Architecture design decisions:
- I used Gin for the web framework because it seems to be simple, fast and with good performance according to online resources. It seems mature and well maintained. Has JSON built in ssupport and supports middlewares that would allow to easily add more features like authentication if needed.
- I used Zap for logging because according to my quick investigation, Go's log package is not as good as Zap, also, Zap supports logging levels, and easy to log JSONs for structured logging (to be hooked to Grafana, Datadog, Logz.io or others).
- Use of Air package for development, I found the lack of live reloading in go and gin to be difficult to work with, therefore I googled and decided to go with the first reasonable live reloading package I found and not waste too much time looking for one.

### Questions

#### How long did you spend working on the problem? What did you find to be the most difficult part?

Around 8 hours, which include "Golang research/learn" like using the docs, and skimming through a Udemy course and youtube videos to learn/understand Golang concepts.

The most difficult part was figuring out that I had to use a Mutex, something I did not foresee and I am still getting up to speed with advanced usage of pointers, go routines, channels and mutex.

#### How would you modify your data model or code to account for more kinds of metrics?

Expand the DeviceData model to add other device information that could be important like storage or memory usage, error information and metrics, firmware or updates state or versions, and add secure communication and authentication for the devices.

#### Discuss your solution’s runtime complexity

Up to this first iteration I think looping through the uploadTimes might be the slowest part, I am also concerned about the locking with mutex but need further research to improve it, if I have time available I will look into optimizing more the code. 