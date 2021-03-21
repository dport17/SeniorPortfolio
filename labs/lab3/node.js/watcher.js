const fs = require('fs');
fs.watch('target.txt', function(){
	console.log("File 'target.txt' just changed!");
});
console.log("Lab3-3.a.ii by Devin Porter. Watching target.txt for changes.");
