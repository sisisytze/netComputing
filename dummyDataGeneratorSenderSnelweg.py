## dummyDataGenerator for Environmental Statistics

import random
import socket
import sys
import time
import pika

ringwegLocs = [
            [53.200061, 6.564646],
            [53.197052, 6.563187],
            [53.193863, 6.563616],
            [53.190751, 6.565376],
            [53.187562, 6.567393],
            [53.183524, 6.574817],
            [53.179280, 6.580396],
            [53.176321, 6.583958],
            [53.198505, 6.546545],
            [53.196293, 6.539636],
            [53.195701, 6.532340],
            [53.196369, 6.526804],
            [53.196985, 6.521268]
        ]

def randomMAC():
    return random.getrandbits(32)

def genData(locList, averageWeight, spread):
    for location in locList:
        message =  str(randomMAC()) + ";" + str(time.time()) + ";" + str(location[0]) + ";" + str(location[1]) + ";" + "CO2" + ";" + str(round(random.uniform(averageWeight - spread, averageWeight + spread),2))
        
        connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
        channel = connection.channel()
        channel.queue_declare(queue='sensor_data')
        channel.basic_publish(exchange='',routing_key='sensor_data',body=message)    
        connection.close() 


genData(ringwegLocs, 50, 5)

# {location: new google.maps.LatLng(53.2263633, 6.5444437), weight: 6.0},   