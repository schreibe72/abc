# abc (Azure Blob Command)

This Tool can upload files to a azure storage account and sets the content type of this file. It also can create, delete and list container.


# Installation

    go get github.com/schreibe72/abc

# Usage

## save azure storage account credentials
You can provided the access key to your azure storage account per command line or you can store it in an config file.

```
abc save credentials -h
Save Key and Account

Usage:
  abc save credentials [flags]

Flags:
  -f, --filename string   Filename to store Config (default "/Users/manfred/.azure/abc-config.json")

Global Flags:
  -a, --account string   Azure Storage Account
  -k, --key string       Azure Storage Account Key
  -v, --verbose          Verbose Output
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
uploads a file to a selected container and stets the contentType

Usage:
  abc upload [flags]

Flags:
  -b, --big                      spilt file which are bigger than 195GB in part Blockblobs
  -n, --blobname string          The Blob File Name (required for pipe)
  -C, --cacheControl string      CacheControl for the uploaded file
  -c, --container string         a Azure Container (required)
  -E, --contentEncoding string   ContentEncoding for the uploaded file
  -L, --contentLanguage string   ContentLanguage for the uploaded file
  -T, --contentType string       Contenttype for the uploaded file
  -f, --filename string          Filename to upload (required if no pipe)
  -p, --pipe                     incoming Pipe
  -w, --worker int               download Worker Count (default 10)

Global Flags:
  -a, --account string   Azure Storage Account
  -k, --key string       Azure Storage Account Key
  -v, --verbose          Verbose Output
```

**Upload a big File per pipe:**
```
tar cz /var/lib/mysql | abc upload -b -p -n backup-db.tgz -c backup-20160807
```

**Upload a file(with automatic contentType by extension):**
```
abc upload -c pictures -n cat.jpg -C "cache-control: private, max-age=60, no-cache"
```

## download a blob

this subcommand let you download a blob to a file or put the output to stdout. If the blobfile is a bundle and there is  *blobname*-bundle.json file all parts will be downloaded and chained in to one file or output to stdout.

```
abc download -h
Downloads a blob file selected in you storage account and selected container

Usage:
  abc download [flags]

Flags:
  -n, --blobname string    The Blob File Name (required)
  -c, --container string   a Azure Container (required)
  -f, --filename string    Filename to download
  -p, --pipe               outgoing Pipe
  -w, --worker int         download Worker Count (default 10)

Global Flags:
  -a, --account string   Azure Storage Account
  -k, --key string       Azure Storage Account Key
  -v, --verbose          Verbose Output
```

**download a blob to stdout**
```
abc download -c backup -n mysql.tar.gz -p | tar xz

```

**download a blob to file**
```
abc download -c backup -n mysql.tar.gz
```

# container operations

## container list operation
```
abc container list -h
Here you can list all containers in your storage account. You can also list
       all container by a certain prefix.

Usage:
 abc container list [flags]

Flags:
 -p, --prefix string   a Azure Container Prefix

Global Flags:
 -a, --account string   Azure Storage Account
 -k, --key string       Azure Storage Account Key
 -v, --verbose          Verbose Output
```

**list all containers**
```
abc container list
```
**list all containers starts with prefix**
```
abc container list -p backup
```

## container show operation
```
abc container show -h
show the containing blobs in your storage account

Usage:
  abc container show [flags]

Flags:
  -c, --container string   a Azure Container (required)
  -p, --prefix string      a Azure Blob Prefix

Global Flags:
  -a, --account string   Azure Storage Account
  -k, --key string       Azure Storage Account Key
  -v, --verbose          Verbose Output
```
**show content of a container (all blobs)**
```
abc container show -c pictures
```
**show content of container with blob prefix**
```
abc container show -c pictures -p /test/
```
## container delete operation
```
abc container delete -h
Here you can delete a container in your storage account

Usage:
  abc container delete [flags]

Flags:
  -c, --container string   a Azure Container (required)

Global Flags:
  -a, --account string   Azure Storage Account
  -k, --key string       Azure Storage Account Key
  -v, --verbose          Verbose Output
```
**delete container**
```
abc container delete -c pictures
```
## container create operation
```
abc container create -h
Here you can create a container in your storage account

Usage:
  abc container create [flags]

Flags:
  -c, --container string   a Azure Container (required)

Global Flags:
  -a, --account string   Azure Storage Account
  -k, --key string       Azure Storage Account Key
  -v, --verbose          Verbose Output
```
**create container**
```
abc container create -c pictures
```
## blob delete operation
```
abc blob delete -h
delete selected blob files in your storage account

Usage:
  abc blob delete [flags]

Flags:
  -n, --blobname string    The Blob File Name (required)
  -c, --container string   a Azure Container (required)

Global Flags:
  -a, --account string   Azure Storage Account
  -k, --key string       Azure Storage Account Key
  -v, --verbose          Verbose Output
```
**delete blob**
```
abc blob delete -c pictures -n /test/cat.jpg
```

# Bash Complition
```
abc save bash -h
Save Bash Complition

Usage:
  abc save bash [flags]

Flags:
  -f, --filename string   Filename to store Bash Complitition

Global Flags:
  -a, --account string   Azure Storage Account
  -k, --key string       Azure Storage Account Key
  -v, --verbose          Verbose Output
```
**Install bash complition**

```
abc save bash -f ~/abc.completion.sh
echo "source ~/abc.completion.sh" >> .bash_profile
```

On Mac OS X there is an  *_get_comp_words_by_ref command not found* error. This can be fixed with:
```
brew install bash-completion
```

bash completion works only with bash version 4. If there is an *flaghash["${flagname}"]: bad array subscript* error you need to Install
bash version 4. For mac:
```
brew install bash
```
You also have to configure your terminal programm to use this shell.

-----
# License

This project is published under [Apache 2.0 License](LICENSE).
