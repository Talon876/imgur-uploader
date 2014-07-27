var express = require('express');
var logger = require('morgan');
var bodyParser = require('body-parser');
var busboy = require('connect-busboy');

var config = require('./config').config;
var routes = require('./routes');
var imageRoute = require('./routes/image');

var app = express();
app.use(logger('dev'));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded());
app.use(busboy({
	limits: {
		fileSize: config.maxFileSizeMb * 1024 * 1024
	}
}));
app.use('/image', imageRoute);

app.use(function(req, res, next) {
	var err = new Error('Not Found');
	err.status = 404;
	next(err);
});

if (app.get('env') === 'development') {
	app.use(function(err, req, res, next) {
		res.status(err.status || 500);
		res.send('internal server error: ' + 
			err.message + 
			'Details: ' + err);
	});
}

app.use(function(err, req, res, next) {
	res.status(err.status || 500);
	res.send('internal server error');
});

module.exports = app;

app.listen(config.port);

