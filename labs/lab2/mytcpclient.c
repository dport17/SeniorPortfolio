/* include libraries */
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <unistd.h>
#include <netdb.h>
#include <string.h> 

int main (int argc, char *argv[])
{

	if(argc!=3){
		printf("Usage: %s <servername> <port>\n", argv[0]);
		exit(1);
	}
	char *servername = argv[1];
	char *port = argv[2];
   printf("TCP CLient program by Devin Porter\n");
   printf("Servername = %s, port = %s\n", argv[1], argv[2]);

   int sockfd = socket(AF_INET, SOCK_STREAM, 0);
   if(sockfd<0){
   	perror("ERROR Opening socket");
   	exit(sockfd);
   }

   printf("A socket is opened!\n");
   struct addrinfo hints, *serveraddr;
   memset(&hints, 0, sizeof hints);
   hints.ai_family = AF_INET;
   hints.ai_socktype = SOCK_STREAM;
   int addr_lookup = getaddrinfo(servername, port, &hints, &serveraddr);

   if (addr_lookup != 0) {
   	fprintf(stderr, "getaddrinfo: %s\n",gai_strerror(addr_lookup));
   	exit(3);
   }
   int connected = connect(sockfd, serveraddr->ai_addr,serveraddr->ai_addrlen);
	if(connected < 0){
		perror("Cannot connect to the server\n");
		exit(3);
	}
	else{
		printf("Connected to the server %s at port %s\n",servername, port);
	}
	freeaddrinfo(serveraddr);
	close(sockfd);
}
