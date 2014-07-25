var storage = require('../storage');
var config = require('../config').config;

exports.root = function(req, res){
	res.writeHead(200, {'Content-Type': 'text/plain'});
	res.end('ok');
};

var db = storage.createStorage(config);
exports.redirect = function(req, res){
	var alias = req.params.alias;
	db.getLinkFromAlias(alias, function(err,rows){
		if(err || rows.length == 0) {
			linkNotFound(res, alias);
			return;
		} else {
			var link = rows[0];
			linkRedirect(res, link);
			return;
		}
	});
}

function linkNotFound(res, alias) {
	res.writeHead(404, {'Content-Type': 'text/plain'});
	res.end(alias + ' -> ???');
}

function linkRedirect(res, link) {
	console.log('-> ' + link.url);
	res.redirect(link.url);
}

