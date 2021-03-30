var net = require('net');
var mode = "pub";
var user = "all";

if(process.argv.length != 4){
	console.log("Usage: node %s <host> <port>", process.argv[1]);
	process.exit(1);	
}

var host=process.argv[2];
var port=process.argv[3];

if(host.length >253 || port.length >5 ){
	console.log("Invalid host or port. Try again!\nUsage: node %s <port>", process.argv[1]);
	process.exit(1);	
}

var client = new net.Socket();
console.log("Simple telnet.js developed by Phu Phung, SecAD");
console.log("Connecting to: %s:%s", host, port);

client.connect(port,host, connected);

function connected(){
	console.log("Connected to: %s:%s", client.remoteAddress, client.remotePort);
	console.log("You are currently in public messaging mode.");
	console.log("Type private, direct, dm, or pm to enter private messaging.");
	console.log("Type .exit to quit.");
	setTimeout(() => { runLogin() }, 1000);
	
}

function runLogin(){
	keyboard.question("Username: ",function(user){
		keyboard.question("Password: ", function(pass){
			var creds = { username : `${user}`, password : `${pass}` }
			var data = JSON.stringify(creds);
			client.write("login  "+data);
		});
	});
}

function setUser(){
	return new Promise(resolve => keyboard.question("Who would you like to DM?", ans =>{
		resolve(ans);
	}))
}

async function switchToPriv(){
	console.log("Switching to private mode... ");
	user = await setUser();
	mode = "pr";
	console.log("Switched to private chat with "+user+". If you'd like to switch to public chat, type 'Public'.");
}


client.on("data",data=>{
	if(data.toString()=="LF"){
		console.log("The login was not successful, please try again.");
		runLogin();
	}else{
		console.log("New message: "+data);
	}
	console.log(data.toString());
});

client.on("error",function(err){
	console.log("error");
	process.exit(2);
});

client.on("close",function(data){
	console.log("Connection disconnected");
	process.exit(3);
});

const keyboard = require('readline').createInterface({
	input: process.stdin,
	output: process.stdout
});

keyboard.on('line',(input)=>{
	console.log(`You typed: ${input}`);
	if(input==".exit"){
		client.destroy();
		console.log("Disconnected!");
		keyboard.close();
	}else{

		if((input == "direct" || input == "private" || input == "dm" || input == "pm")&&mode!="pr"){
			switchToPriv();
			return;
		}
		else if(input == "Public" || input == "public"){
			mode = "pub";
			user = "all";
			console.log("Switched to public mode.");
			return;
		}
		var msg = {mode : mode, user: user, msg: input}

		client.write(JSON.stringify(msg));
	}
});