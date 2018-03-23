import socket
import sys
import time
import random
import pika

from uuid import getnode as get_mac


connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
channel = connection.channel()


channel.queue_declare(queue='sensor_data')


## single measurement payload 
# 1. mac / unique id
# 2. sensortimestamp
# 3. latitude
# 4. longitude
# 5. sensortype (int)
# 6. data (float)

while True:
    mac = str(get_mac())
    sensortimestamp = str(time.time())
    latitude = str(389457.938457)
    longitude = str(549457.938564)
    sensortype = str(5)
    data = str(random.uniform(1.5, 1.9))    

    # Send data
    print >>sys.stderr, 'sending measurement' 
    message =  mac + "|" + sensortimestamp + "|" + latitude + "|" + longitude + "|" + sensortype + "|" + data
    channel.basic_publish(exchange='',routing_key='hello',body=message)    
    
    time.sleep(10)
    