# Food trucks 

## Tech Stacks
### Backend
- Redis
- Iris
### Frontend
- PrimeReact

## Bootstrap
- run redis on docker
- build frontend
- build backend
- start backend

## Features and implementations
### Search nearby trucks
#### Frontend
- Frontend(React) using React GMap Component https://www.primefaces.org/primereact-v8/gmap/
- Frontend get a list of Facilities by passing location(longitude, latitude) and radim to backend
- Give a link for each facility to go to Facility's detail page
### Backend
- Backend(go) using Redis's GeoLocation function to query nearby Facilities
- Seed cvs to redis when app launch 
### Search Facility by Food Items(Cli)
- Using Redis's zset to index Facilities by food items, leave score to 0, so the result is ordered by facilities Name
### Facility introduction Page
