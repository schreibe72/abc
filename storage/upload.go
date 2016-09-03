package storage

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	as "github.com/Azure/azure-sdk-for-go/storage"
)

func (a *StorageAttributes) saveBlob(reader io.Reader, container string, name string, contentSetting ContentSetting) (bundleItem, error) {

	if err := validateBlobName(container, name); err != nil {
		return bundleItem{}, err
	}

	blocklist := make([]as.Block, 0, 500)
	var err error
	var l int
	var wg sync.WaitGroup
	var item bundleItem
	c := make(chan *uploadBlock)
	freeList := make(chan *uploadBlock, (a.WorkerCount + 1))

	item.EOF = false
	item.BlobName = name

	hash := md5.New()
	for i := 0; i < a.WorkerCount; i++ {
		wg.Add(1)
		go uploadWorker(a, container, name, c, freeList, &wg)
	}

	for {
		var u *uploadBlock
		select {
		case u = <-freeList:

		default:

			buf := make([]byte, as.MaxBlobBlockSize)
			u = new(uploadBlock)
			u.buf = buf
		}
		l, err = io.ReadFull(reader, u.buf)
		u.buf = u.buf[:l]

		item.BlobSize += uint64(l)

		hash.Write(u.buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = nil
			item.EOF = true
		} else if err != nil {
			return item, err
		}
		u.id = fmt.Sprintf("%032d", len(blocklist))
		status := as.BlockStatus(as.BlockStatusUncommitted)
		block := as.Block{ID: u.id,
			Status: status}
		blocklist = append(blocklist, block)
		c <- u
		if item.EOF || len(blocklist) >= MaxBlobBockCount {
			break
		}
	}
	close(c)
	wg.Wait()
	if a.Verbose {
		a.Logger.Print("\n")
	}
	err = a.blobService.PutBlockList(container, name, blocklist)
	if err != nil {
		return item, err
	}
	b64hash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	cs := as.BlobHeaders(contentSetting)
	cs.ContentMD5 = b64hash

	item.BlobMD5 = fmt.Sprintf("%x", hash.Sum(nil))

	err = a.blobService.SetBlobProperties(container, name, cs)
	return item, err
}

func (a *StorageAttributes) SaveBlob(reader io.Reader, container string, name string, big bool, contentSetting ContentSetting) error {

	if err := validateBlobName(container, name); err != nil {
		return err
	}

	var item bundleItem
	var err error
	item.EOF = false
	var b bundle
	b.Bundle = make([]bundleItem, 0, 1000)
	b.FileName = name
	b.ContentType = contentSetting.ContentType
	blobname := name
	for !item.EOF {
		if big {
			blobname = fmt.Sprintf("%s-%03d", name, len(b.Bundle))
		}
		if a.Verbose {
			a.Logger.Printf("Save Blob in new BlobFile: %s\n", blobname)
		}
		item, err = a.saveBlob(reader, container, blobname, contentSetting)
		if err != nil {
			return err
		}
		b.Bundle = append(b.Bundle, item)
		if !item.EOF {
			big = true
		}
	}
	if big {
		err := a.storeBundleFile(container, name, b)
		return err
	}
	return nil
}

func (a *StorageAttributes) storeBundleFile(container string, name string, b bundle) error {

	if err := validateBlobName(container, name); err != nil {
		return err
	}

	name = fmt.Sprintf("%s-bundle.json", name)
	buf, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	if a.Verbose {
		a.Logger.Printf("Create BundleFile: %s\n", name)
	}
	hash := md5.New()
	blocklist := make([]as.Block, 0, 10)
	id := fmt.Sprintf("%032d", 1)
	status := as.BlockStatus(as.BlockStatusUncommitted)
	block := as.Block{ID: id,
		Status: status}
	blocklist = append(blocklist, block)
	hash.Write(buf)
	err = a.blobService.PutBlock(container, name, id, buf)
	if err != nil {
		return err
	}
	err = a.blobService.PutBlockList(container, name, blocklist)
	if err != nil {
		return err
	}
	var cs as.BlobHeaders
	cs.ContentMD5 = base64.StdEncoding.EncodeToString(hash.Sum(nil))
	cs.ContentType = "application/json"
	err = a.blobService.SetBlobProperties(container, name, cs)
	return err
}

func uploadWorker(a *StorageAttributes, container string, name string, c <-chan *uploadBlock, freeList chan<- *uploadBlock, wg *sync.WaitGroup) {
	defer wg.Done()
	trys := 3
	for u := range c {
		for {
			trys--
			err := a.blobService.PutBlock(container, name, u.id, u.buf)
			if err == nil {
				trys = 3
				break
			} else if trys < 0 {
				panic(fmt.Sprintf("Can't upload Block %s! Panic!", u.id))
			} else {
				// Wait 5 sec before retry
				time.Sleep(5 * time.Second)
			}
		}
		freeList <- u
	}
}
