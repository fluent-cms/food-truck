import React, {useState} from "react";
import {MapContainer, Marker, Popup, TileLayer, useMap} from 'react-leaflet'
import useSWR from 'swr'
import {fetcher} from "../utils/fetcher";
import {Config} from "../config";
import {FacilityMarkers} from "./FacilityMarkers";
import {Button} from "primereact/button";
import {SwitchLocation} from "./SwitchLocation";

export default function Map() {
    const {data: center} = useSWR(Config.APIHost + '/api/facilities/center', fetcher)
    const [your, setYour] = useState(false)

    return (
        <div>
            {!your && <Button onClick={()=> setYour(true)} label="Your Location" style={{ backgroundColor: 'var(--primary-color)', color: 'var(--primary-color-text)'}}/>}
            {center && <MapContainer center={[center.Lat, center.Lon]} zoom={18} scrollWheelZoom={false}>
                <TileLayer
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
                <FacilityMarkers/>
                {your && <SwitchLocation/>}
                </MapContainer>
            }
        </div> );
}