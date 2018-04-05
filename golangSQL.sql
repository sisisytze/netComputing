   SELECT m2.value, s2.latitude, s2.longtitude
	 FROM measurement AS m2 JOIN sensor AS s2 ON m3.sensor_uuid = s2.uuid,  
	 (
	     SELECT m1.uuid, MAX(m1.at) AS moment
			 FROM measurement AS m1 INNER JOIN sensor AS s1 ON m1.sensor_uuid = s1.uuid
			 WHERE
			   m1.at > "time - 1h" AND
				 m1.at < "time" AND
				 s1.longtitude > "long - zoom" AND
				 s1.longtitude < "long + zoom" AND
				 s1.latitude > "lat - zoom" AND
				 s1.latitude > "lat + zoom"
			GROUP BY m1.sensor_uuid
	 ) AS m3
 	 WHERE m3.uuid = m2.uuid;
