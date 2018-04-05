/*lock database*/
SELECT s1.*
FROM sensor AS s1 INNER JOIN sensor_synched AS ss ON s1.uuid = ss.sensor_uuid

TRUNCATE TABLE sensor_synched;
/*free database*/

/*lock database*/
SELECT s1.*
FROM measurement AS m1 INNER JOIN measurement_synched AS sm ON m1.uuid = sm.measurement_uuid

TRUNCATE TABLE sensor_synched;
/*free database*/
