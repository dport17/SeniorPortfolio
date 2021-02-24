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
		printf("Usage: %s <servername> <path>\n", argv[0]);
		exit(1);
	}
	char *servername = argv[1];
	char *path = argv[2];
   printf("TCP CLient program by Devin Porter\n");
   printf("Servername = %s, path = %s\n", argv[1], argv[2]);

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
   int addr_lookup = getaddrinfo(servername, "http", &hints, &serveraddr);

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
		printf("Connected to the server %s at  %s\n",servername, path);
	}
	freeaddrinfo(serveraddr);

	int BUFFERSIZE = 1024; //define the size of the buffer
	char buffer[BUFFERSIZE]; //define the buffer
	bzero(buffer,BUFFERSIZE);
	sprintf(buffer, "GET %s HTTP/1.1\r\nHost: %s\r\n\r\n",path,servername);
	int byte_sent = send(sockfd,buffer, strlen(buffer), 0);

	bzero(buffer,BUFFERSIZE);
	int byte_received = recv(sockfd, buffer, BUFFERSIZE, 0);
	if(byte_received <0){
		perror("Error in reading");
		exit(4);
	}
	printf("Received from server: %s", buffer);
	close(sockfd);
	return 0;
}
