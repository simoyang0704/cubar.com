CREATE FUNCTION calc_distance (lat1 DECIMAL(10,6), lng1 DECIMAL(10,6), lat2 DECIMAL(10,6), lng2 DECIMAL(10,6))
RETURNS DECIMAL(10,6)
RETURN 6371 * acos(sin(radians(lat1))*sin(radians(lat2))+cos(radians(lat1))*cos(radians(lat2))*cos(radians(lng2)-radians(lng1)))