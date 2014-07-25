var mysql = require('mysql');

function handleDisconnect(self, connection) {
	connection.on('error', function(err) {
		if (!err.fatal) {
			return;
		}

		if (err.code !== 'PROTOCOL_CONNECTION_LOST') {
			throw err;
		}
		console.log('Attempting to re-establish connection: ' + err.stack);
		self.connection = mysql.createConnection(connection.config);
		handleDisconnect(self, self.connection);
		self.connection.connect();
	});
}

function Storage(config){
	var connection = mysql.createConnection(config.mysql);
	handleDisconnect(this, connection);
	connection.connect();
	this.connection = connection;
}
exports.createStorage = function(config){
	return new Storage(config);
}

Storage.prototype.getLinkFromAlias = function(alias, callback){
	var query = 'SELECT get_link_from_alias(?) AS url'
	this.connection.query(query, alias, callback);
}
