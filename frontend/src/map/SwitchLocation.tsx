import {useEffect, useState} from "react";
import {Marker, Popup, useMap} from "react-leaflet";

export function SwitchLocation(){
    const [position, setPosition] = useState<L.LatLng | null>(null);
    const map = useMap();
    const [loading, setLoading] = useState(false)

    useEffect(() => {
            map.locate({setView: true, maxZoom: 16}).on('locationfound', function (e: L.LocationEvent) {
                setLoading(true)
                setPosition(e.latlng);
                map.flyTo(e.latlng, map.getZoom());
                setLoading(false)
            });
    }, [map]);

    return <>
        {loading&&<div>Loading</div>}
        {position === null ? null : (
            <Marker position={position}>
                <Popup>You are here</Popup>
            </Marker>)
        }
    </>
}