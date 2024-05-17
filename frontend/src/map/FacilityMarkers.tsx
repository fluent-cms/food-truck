import {Marker, Popup, useMap, useMapEvents} from "react-leaflet";
import {getDistanceFromLatLonInKm} from "../utils/getDistance";
import useSWR from "swr";
import {Config} from "../config";
import {fetcher} from "../utils/fetcher";
import {useState} from "react";

export  function FacilityMarkers(){
    const [location, setLocation] = useState<any>({lat:0, lng:0, radius:0})
    const getLocation = ()=> {
        const {lat, lng} = map.getCenter()
        const se = map.getBounds().getSouthEast()
        const radius = getDistanceFromLatLonInKm(lat, lng, se.lat, se.lng)
        setLocation({ lat,lng,radius})
    }
    const map = useMapEvents({
        zoomend:(e) => {
            getLocation()
        },
        moveend:(e) => {
            getLocation()
        }
    })
    const {data : facilities} = useSWR<Facility[]>(Config.APIHost + `/api/facilities?lat=${location.lat}&lon=${location.lng}&radius=${location.radius}`, fetcher)
    return facilities &&<>
        {facilities.map(item => {
           return <Marker key={item.locationID} position={[item.latitude, item.longitude]} >
               <Popup>
                   {item.applicant}<br/>{item.locationDescription} <br/> {item.foodItems}
               </Popup>
           </Marker>
        })}
    </>
}
