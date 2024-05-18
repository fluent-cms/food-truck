# Food trucks 
This is a demo project showcase usage of Go, React, Redis, Docker.
Data is from San Francisco's food truck open dataset.

## Features
As a user, when you go to the websites' home page, you can see a list of marker, 
each marker represents a food truck.

![img.png](doc/images/home-page.png)

When you click a marker, you can see the food truck's applicant, location, and food items.

![img.png](doc/images/pop.png)

If you live in San Francisco, you want see the trucks near your location, 
you can click 'Your Location' button, the map will be switched to your location

![img.png](doc/images/your-location.png)

As a admin, you can search all trucks who are serving a type of food, e.g. taco

![img.png](doc/images/cli.png)

## Tech Stacks
### Backend
- Redis as in memory database
  - redis str, get truck(marshalled as json) by ID
  - redis geo, to search nearby trucks by latitude, longitude, and radius.
  - redis zset, to search a list of trucks by food items it served.
- Go  
- Iris Web Framework

### Frontend
- React
- react-leaflet for map related feature
- swr for state management

## Design pattern and best practice
- *hexagonal architecture*

![img_1.png](doc/images/hexagonal.png)

The core of backend is /backend/packages/services/facilitySvc.go, service layer doesn't depend on
storage layer, and doesn't depend on UI layer.   
Both Cli and Web can use facility service.

- *Dependency Injection*

facilitySvc is not hardcoded depending on rdb package(my own Redis lib), so if I want
change storage to mysql or mongodb later, I can implement the interface
and inject the implementation to service.

This also conform to Open/Close principle, the facilitySvc is open to extend functionality, 
but close to code change

- *Separation of Concern*

When do frontend coding, I also tried to apply this principle, for the truck map frontend.
I use 3 components to render map (Map, FacilityMaker, SwitchLocation), each component care about it's own job, improved readability.

- *Modular and DRY - Don't repeat your self*

I aimed to separate business logic from infrastructure code in our project. Taking the `facilitySvc` as an example:
each facility can have multiple food items, and each food item can be associated with multiple facilities.
This relationship pertains to business logic. In contrast, connecting to Redis and marshalling objects to JSON strings are common infrastructure tasks.

By wrapping Redis operations into a standalone package, instead of embedding this code within the facility service,
I made the codebase more modular and reusable. This separation improves maintainability and allows infrastructure code
to be reused across different services without duplication.

- *Template pattern*

There are a lot of boilerplate to start a web application, 
I put these code to /backend/packages/util/irisbase applying Template Pattern,
so the main file(/backend/cmds/web/main) looks clean and straightforward.

- *Error handling*  

Each function annotate error detail (e.g. which line throws the error, the cause of the error).
In develop mode, the API returns error detail to help frontend user to locate the issue.
In production mode, the API just return an 500 error to hide technical detail

- *Generic Programming to improve ability to reuse code*

For example, the parse function in /backend/packages/util/yaml.go demonstrates how to create a reusable utility for parsing 
YAML files into any specified type. By using generics, we can create a single, versatile function that works with 
any data structure, enhancing our ability to write clean, reusable, and maintainable code.
```
func parse[T any](t *T, filePath string) (err error) {
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	if err = yaml.NewDecoder(f).Decode(t); err != nil {
		return err
	}
	return nil
}
 ```


## Implementations
### Frontend
#### Code Structure
```
--frontend/
----src/
------map/
--------FacilityMarkers.tsx   # add markers to map
--------Map.tsx               # map container
--------SwitchLoaction.tsx    # switch to your location
------models/
------utils/
------config.ts               # global configs
----.env.development          # development enviroment virables
```
#### Api call and state management
All meaningful code resides in /frontend/src/map, code in models and utils is very simple. 
the useSWR hook combine api call and state management, one single line of code save the trouble of useEffect hook. 
```
    const {data: center} = useSWR(Config.APIHost + '/api/facilities/center', fetcher)
```
#### Environment variables
The frontend app might run in two mode 
##### Development Mode 
In development mode (pnpm dev), I start two web server, 
http://localhost:8080 as backend  http://localhost:5173 as frontend, so the api endpoint is http://localhost:8080/api/***.
  I put this dev environment settings to .env.development
```
VITE_REACT_APP_API_HOST='http://localhost:8080'
```
##### Production Mode
In production mode frontend and backend are served as single app(the distribution of frontend is copied to 
/backend/web, and served by backend web server). 
I can use relative path to call backend api /api/facilities. In production there won't be .env file, 
so the api host default to empty string''. 

##### config.js
All environment reading code are put into config.ts , ensure single source of truth.
```
export const Config = {
    APIHost : import.meta.env.VITE_REACT_APP_API_HOST || '',
}
```
#### Map Related Features
I tried google map API first, but it's not totally free, I don't want checkin API Key to repo. And I want 
people can easily play with this app, so I followed this link https://medium.com/@ujjwaltiwari2/a-guide-to-using-openstreetmap-with-react-70932389b8b1 
to use react-leaflet
### Backend
#### Code structure
```
--cmds
----cli/          # entrance of cli
----web/          # entrance of web
--packages
----controllers/  # endpoint of APIs
----models/      
----services/     # implement business logic of food trucks
----utils/        # infrastrcutres
--web/            # put frontend distribution here
```
#### Seed Data
In packages/services/facilitySvc Seed() function, it read configs/data.csv, and parse it as Facility array,
then populate the data to redis.

### Endpoints
- */api/facilities/center*  Get the center of all trucks
- */api/facilities?lat=&lon=&radius=* Get the facilities near the center with in the radius 
### Cli 
- share Facility Service with web, provides function of search facility by food items

## Installation
If you don't have go, node, pnpm installed on you local machine, you can simply use docker compose to start the app.
### Docker 
in the root directory of the project, run
```shell
docker-compose up
```
When you see messages similar to below, then the app is up.
```shell
food-truck-backend-1  | Now listening on:
food-truck-backend-1  | > Network:  http://172.21.0.3:8080
food-truck-backend-1  | > Local:    http://localhost:8080
food-truck-backend-1  | Application started. Press CTRL+C to shut down.
```
### Web
Use your browser, go to http://localhost:8080 to see the food truck app.
### Cli
to run cli
```shell
# use docker ps to check docker container name
⚡➜ ~ docker ps
CONTAINER ID   IMAGE                COMMAND                  CREATED         STATUS         PORTS                    NAMES
0f6039ec90d6   food-truck-backend   "./main"                 5 minutes ago   Up 5 minutes   0.0.0.0:8080->8080/tcp   food-truck-backend-1
40f20e7696e4   redis:latest         "docker-entrypoint.s…"   5 minutes ago   Up 5 minutes   0.0.0.0:6379->6379/tcp   food-truck-redis-1

# start a shell session
⚡➜ ~ docker exec -it food-truck-backend-1 /bin/bash

# in the shell session,  run ./food-cli
root@0f6039ec90d6:/go/src/app# ./food-cli
Load Config from  ./configs/cli.yaml

# when you saw the 'Enter Food Item to search facility:', input a food item you want to find
Enter Food Item to search facility: breakfast
Munch A Bunch MISSION ST: 14TH ST to 15TH ST (1800 - 1899)
Munch A Bunch BRYANT ST: ALAMEDA ST intersection
Munch A Bunch FULTON ST: FRANKLIN ST to GOUGH ST (300 - 399)
Munch A Bunch LARKIN ST: FERN ST to BUSH ST (1127 - 1199)
Munch A Bunch 12TH ST: ISIS ST to BERNICE ST (332 - 365)
Munch A Bunch 07TH ST: CLEVELAND ST to HARRISON ST (314 - 399)
Munch A Bunch PARNASSUS AVE: HILLWAY AVE to 03RD AVE (400 - 599)
Munch A Bunch 17TH ST: SAN BRUNO AVE to UTAH ST (2200 - 2299)
```

## Development
### Spin up a redis server
```shell
docker run --name food-redis -d -p 6379:6379 redis
```
### Start backend
go to /backend, 
```shell
go run backend/cmds/web/main.go
```
if you got the following error, it means backend can not connect to a host 'redis'
```
panic: facilitySvc.go:58 failed to cache facilities, dial tcp: lookup redis: no such host
```
you can add 'redis' to you development machine's /etc/hosts file
```
127.0.0.1       localhost redis
```
or you can modify /backend/configs/web.yaml file, change the following line
```
  addr: redis:6379
```
to 
```
  addr: localhost:6379
```
### Start CLI
```
go run backend/cmds/cli/main.go
```
Cli's config file is at /backend/configs/cli.yaml

### Start Frontend
you need to install node, pnpm, then go to frontend/, run
```
pnpm install
pnpm dev
```