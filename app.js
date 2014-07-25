var express = require('express');
var http = require('http');
var path = require('path');
var config = require('./config').config;
var routes = require('./routes');

var app = express();
app.configure(function(){
	app.set('port', process.env.PORT || config.port);
	app.use(express.favicon());
	app.use(express.logger('dev'));
	app.use(express.bodyParser());
	app.use(express.methodOverride());
	app.use(app.router);
	app.use(express.static(path.join(__dirname, 'public')));
});

app.configure('development', function(){
	app.use(express.errorHandler());
});

app.get('/', routes.root);
app.get('/:alias', routes.redirect);

http.createServer(app).listen(app.get('port'), function(){
	console.log('Imgur redirect server listening on port '
		+ app.get('port'));
});

