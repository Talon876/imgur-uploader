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


exports.addlink = function(req, res){
	res.form.complete(function(err, fields, files){
		if(err){
			next(err);
		} else{
			ins = fs.createReadStream(files.photo.path);
			ous = fs.createWriteStream('/tmp');
			util.pump(ins, ous, function(err){
				if(err){
					next(err);
				} else{
					res.redirect('/image');
				}
			});
			console.log('\nUploaded %s to %s', files.photo.filename, files.photo.path);
			res.send('Uploaded ' + files.photo.filename + ' to ' + files.photo.path);
		}

	});

}

exports.addlinkForm = function(req, res){
  res.send('<form method="post" enctype="multipart/form-data">'
  + '<p>Data: <input type="filename" name="filename" /></p>'
  + '<p>file: <input type="file" name="image" /></p>'
  + '<p><input type="submit" value="Upload" /></p>'
  + '</form>');
}
