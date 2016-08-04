# abc (Azure Blob Command)

This Tool can upload files to a azure storage account and sets the content type of this file. It also can create, delete and list container.


# Installation

    go get github.com/schreibe72/abc

# Usage

## save azure storage account credentials
You can provided the access key to your azure storage account per command line or you can store it in an config file.

```
abc save credentials -h
Usage of subcommand:
  -a string
    	a Azure Storage Account
  -k string
    	a Azure Key
  -v	Verbose info
```
**Example:**
```
abc save credentials -k kXRpCKCzMak5+RbxqNixfDGVnBgmx8ywiMtIHIiLo+GkedOTfUvzMOy4HJlSrxEQgURzTx654uoAzYmjTazvrQ== -a teststorageaccount0
```

## upload blob
this subcommand uploads files or data from a pipe to a Block Blob. If the Block Blob gets bigger than 195GB, this tool split the blob in parts and create a *blobname*-bundle.json file.
If the data comes per pipe and you know that it is bigger than 195GB you should use the -big switch. Than every file gets an index, also the first file. Remember it is not possible to rename blobs.

```
abc upload -h
Usage of subcommand:
  -a string
    	a Azure Storage Account
  -big
    	spilt file which are bigger than 195GB in part Blockblobs
  -c string
    	a Azure Container
  -cacheControl string
    	CacheControl for the uploaded file
  -contentEncoding string
    	ContentEncoding for the uploaded file
  -contentLanguage string
    	ContentLanguage for the uploaded file
  -contentType string
    	Contenttype for the uploaded file
  -f string
    	Filename to upload
  -k string
    	a Azure Key
  -n string
    	The Blob File Name
  -pipe
    	incoming Pipe
  -v	Verbose info
  -worker int
    	download Worker Count (default 10)
```

**Upload a big File per pipe:**
```
tar cz /var/lib/mysql | abc upload -big -pipe -n backup-db.tgz -c backup-20160807
```

**Upload a file(with automatic contentType by extension):**
```
abc upload -c pictures -n cat.jpg -cacheControl "cache-control: private, max-age=60, no-cache"
```

## download a blob

this subcommand let you download a blob to a file or put the output to stdout. If the blobfile is a bundle and there is  *blobname*-bundle.json file all parts will be downloaded and chained in to one file or output to stdout.

```
abc download -h
Usage of subcommand:
  -a string
    	a Azure Storage Account
  -c string
    	a Azure Container
  -f string
    	Filename to download
  -k string
    	a Azure Key
  -n string
    	The Blob File Name
  -pipe
    	outgoing Pipe
  -v	Verbose info
  -worker int
    	download Worker Count (default 10)
```

**download a blob to stdout**
```
abc download -c backup -n mysql.tar.gz -pipe | tar xz

```

**download a blob to file**
```
abc download -c backup -n mysql.tar.gz
```

# container operations

**list all containers**
```
abc container list
```
**list all containers starts with prefix**
```
abc container list -cp backup
```
**show content of a container (all blobs)**
```
abc container show -c pictures
```
**show content of container with blob prefix**
```
abc container show -c pictures -bp /test/
```
**delete container**
```
abc container delete -c pictures
```
**create container**
```
abc container create -c pictures
```
**delete blob**
```
abc blob delete -c pictures -n /test/cat.jpg
```

-----
# License

This project is published under [Apache 2.0 License](LICENSE).
