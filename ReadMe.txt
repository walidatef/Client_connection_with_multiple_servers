One request for save data in multiple servers

A connection between a client and a number of servers using [ net/http ] and [MySQL] by Go language. 

Idea:
When the client sends a message to one of the servers, it is stored in the server's database [MySQL], then the server sends it to the rest of the servers, and it is also stored in their own database.

In the end, the data on all servers is identical.

Suppose we have three servers:

if Client send data to server 1
	then server 1 send data to server 2 and server 3

if Client send data to server 2
	then server 1 send data to server 1 and server 3

if Client send data to server 3
	then server 1 send data to server 1 and server 2