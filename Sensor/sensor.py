import socket
import sys
import time
import random
import pika

from uuid import getnode as get_mac

## single measurement payload 
# 1. mac / unique id
# 2. sensortimestamp
# 3. latitude
# 4. longitude
# 5. sensortype string
# 6. data (float)

while True:
    mac = str(get_mac())
    sensortimestamp = str(time.time())
    latitude = str(53.240493)
    longitude = str(6.536319)
    sensortype = 'grade'
    data = str(random.uniform(8.5, 9.9))    

    # Send data
    message =  mac + ";" + sensortimestamp + ";" + latitude + ";" + longitude + ";" + sensortype + ";" + data
    print('sending message: ' + message)
    
    connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
    channel = connection.channel()
    channel.queue_declare(queue='sensor_data')
    channel.basic_publish(exchange='',routing_key='sensor_data',body=message)    
    connection.close()    
    time.sleep(10)   