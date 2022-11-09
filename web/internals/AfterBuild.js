const fse = require('fs-extra');
const fs = require('fs');

const srcDir = `./build`;
const destDir = `../cmd/bot/app/web`;

fs.readdir(srcDir, (err, files) => {
	files.forEach(file => {
		console.log(file);
	});
});

// To copy a folder or file
fse.copySync(srcDir, destDir, { overwrite: true });