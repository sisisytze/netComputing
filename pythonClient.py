import socket
import sys
import time

# Connect the socket to the port where the server is listening
server_address = ('localhost', 9001)

#s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
#s.connect(("8.8.8.8", 80))
#name = s.getsockname()[0]

while True:
    try:  
        # Create a TCP/IP socket
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)  
        sock.connect(server_address)
        # Send data
        message = 'message'
        print >>sys.stderr, 'sending "%s"' % message
        sock.sendall(message)        
    finally:
        sock.close()
        time.sleep(10)