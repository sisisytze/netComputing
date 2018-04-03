import socket
import sys
import time
import random
from uuid import getnode as get_mac

# Connect the socket to the port where the server is listening
server_address = ('localhost', 9001)
routing_address = ('localhost', 9002)

#s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
#s.connect(("8.8.8.8", 80))
#name = s.getsockname()[0]

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
    
    try:  
        # Create a TCP/IP socket
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)  
        sock.connect(server_address)
        # Send data
        print >>sys.stderr, 'sending measurement' 
        sock.sendall( mac + "|" + sensortimestamp + "|" + latitude + "|" + longitude + "|" + sensortype + "|" + data)        
    finally:
        sock.close()
        time.sleep(10)