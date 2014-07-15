# Routes

These are the available routes.

## GET /:imageid

The primary route end-users will interact with.
Used for redirecting to the imgur hosted image.

| Status | Reason |
---------|---------
302 | redirects to http://i.imgur.com/matchingId.jpg
404 | somelinkid has no url to redirect to
	

## POST /image

The route used for actually sending the image to the server.
This will require obtaining a secret string from `/secret`.
This route only places the image in the upload directory before
returning to the client. The image is then uploaded in the
background.

| Key | Value |
------|--------
id | the image id you want to associate with
secret | secret string obtained from /secret
image | the image file

| Status | Reason |
---------|---------
200 | success
??? |  invalid image (not an image/too large/etc)
??? | invalid secret

## `POST /secret`

| Key | Value |
------|--------
id | the image id you want to get the secret for

Route used to generate the secret string required to upload and
associate an id with an image.

| Status | Reason |
---------|---------
200 | success; also returns the requested secret
??? | requested id already in use
