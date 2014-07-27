var express = require('express');
var fs = require('fs');
var path = require('path');
var router = express.Router();
var config = require('../config').config;

router.post('/', function(req, res) {
	var filestream;
	req.pipe(req.busboy);
	req.busboy.on('file', function(field, file, filename) {
		var fileLocation = path.join(config.tempDir, filename);
		console.log('saving ' + fileLocation)
		filestream = fs.createWriteStream(fileLocation);
		file.pipe(filestream);
		filestream.on('close', function() {
			res.send('done');
		});
	});
});

module.exports = router;

