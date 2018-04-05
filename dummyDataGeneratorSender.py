## dummyDataGenerator for Environmental Statistics

import random
import socket
import sys
import time
import pika

ringwegLocs = [
            [53.203246, 6.564907],
            [53.205686, 6.575570],
            [53.210628, 6.586484],
            [53.213608, 6.597814],
            [53.217179, 6.609401],
            [53.220467, 6.615795],
            [53.224166, 6.613649],
            [53.229355, 6.610130],
            [53.236860, 6.593625],
            [53.241534, 6.585857],
            [53.246310, 6.581737],
            [53.248903, 6.576501],
            [53.246745, 6.572167],
            [53.240735, 6.568562],
            [53.238346, 6.561138],
            [53.237575, 6.548564],
            [53.235956, 6.539037],
            [53.234491, 6.528265],
            [53.230483, 6.531441],
            [53.226757, 6.534488],
            [53.222306, 6.538189],
            [53.213137, 6.541998],
            [53.207585, 6.548779],
            [53.202623, 6.552298],
            [53.234868, 6.603167],
            [53.217870, 6.539329],
            [53.202519, 6.559937],
            [53.208011, 6.580399],
            [53.212970, 6.592201],
            [53.215359, 6.604131],
            [53.204229, 6.569756]
        ]
        
parkLocs = [
            [53.205860, 6.542687],
            [53.203571, 6.540241],
            [53.199637, 6.537838],
            [53.218449, 6.532525],
            [53.224975, 6.516904],
            [53.241645, 6.543769],
            [53.240437, 6.551451],
            [53.235017, 6.570892],
            [53.228697, 6.582393],
            [53.228131, 6.548704],
            [53.226409, 6.558274],
            [53.224301, 6.555227],
            [53.221371, 6.554583],
            [53.232252, 6.593819],
            [53.192011, 6.540175],
            [53.195422, 6.547616],
            [53.226750, 6.586754],
            [53.221456, 6.553623]
        ]

otherLocs = [
            [53.229756, 6.546128],
            [53.210355, 6.562007],
            [53.232671, 6.557447],
            [53.234489, 6.580317],
            [53.226191, 6.598170],
            [53.220229, 6.590316],
            [53.214806, 6.557443],
            [53.213931, 6.573579],
            [53.227369, 6.570789],
            [53.222538, 6.566540],
            [53.221972, 6.578428]
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

genData(ringwegLocs, 60, 5)
genData(parkLocs, 5, 5)
genData(otherLocs, 30, 15)

# {location: new google.maps.LatLng(53.2263633, 6.5444437), weight: 6.0},   