# S3Snip

A quick and dirty tool for taking screenshots, uploading them to S3 and putting a URL on the clipboard.
Written partly as a way to learn Go, so you know, might break in weird ways. It aims to be a private/personal replacement for puush.me.

Requires a dir `.s3snip` in the users home directory and a `conf.json` within.

`conf.json` should contain:

```
{
	"awsRegion" : "REGION",
	"awsSecretKey" : "SECRET",
	"awsAccessKey" : "KEY",
	"awsBucket" : "BUCKET"
}
```

Your S3 bucket should have the following policy if you want screenshots to automatically be public and you probably do:

```
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "AddPerm",
			"Effect": "Allow",
			"Principal": "*",
			"Action": "s3:GetObject",
			"Resource": [
				"arn:aws:s3:::BUCKETNAME/*",
			]
		}
	]
}
```

Tested only on OSX using `go1.6.2 darwin/amd64`.

Automator can be used to create a service/application for S3Snip and a keyboard shortcut can be bound to trigger it.
