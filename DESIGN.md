Required Functionality/Routes

	* `GET /:imageid`

		The primary link end-users will interact with.
		Used for redirecting to the imgur hosted image.

		Example: http://snip.snippingtoolpluspl.us/someimageid.jpg
		-> 302 redirects to http://i.imgur.com/matchingId.jpg
		-> 404 somelinkid has no url to redirect to
		

	* `POST /image`

		{
			id:"someimageid",
			image:"<multi-part-file>",
			secret:"secret required to set somelinkid"
		}

		The route used for actually sending the image to the server.
		This will require obtaining a secret string from `/secret`.
		This route only places the image in the upload directory before
		returning to the client. The image is then uploaded in the
		background.

		-> 200 ok
		-> ??? invalid image (not an image/too large/etc)
		-> ??? invalid secret
	
	* `POST /secret`

		{
			id:"someimageid"
		}

		Route used to generate the secret string required to upload and
		associate an id with an image.

		-> 200 
			{
				secret: "<the secret for the requested id>"
			}
		-> ??? requested id already in use


