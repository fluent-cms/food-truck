export function getDistanceFromLatLonInKm(lat1: number, lon1: number, lat2: number, lon2: number): number {
    const R: number = 6371; // Radius of the Earth in kilometers
    const dLat: number = deg2rad(lat2 - lat1); // Convert latitude difference to radians
    const dLon: number = deg2rad(lon2 - lon1); // Convert longitude difference to radians
    const a: number =
        Math.sin(dLat / 2) * Math.sin(dLat / 2) +
        Math.cos(deg2rad(lat1)) * Math.cos(deg2rad(lat2)) *
        Math.sin(dLon / 2) * Math.sin(dLon / 2);
    const c: number = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    const distance: number = R * c; // Distance in kilometers
    return distance;
}
function deg2rad(deg: number): number {
    return deg * (Math.PI / 180);
}