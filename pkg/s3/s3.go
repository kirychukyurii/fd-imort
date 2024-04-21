package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/webitel/wlog"

	"github.com/kirychukyurii/fd-import/config"
)

// MaxListKeys is the maximum number of keys to be listed in the ListObjects method of the Bucket type.
const MaxListKeys = 10000

// Bucket represents a container for storing objects.
//
// name is the name of the bucket.
//
// log is a logger used for logging.
//
// cli is an S3 client for interacting with the AWS S3 service.
//
// objectQueue is a queue of objects.
//
// errorsCh is a channel used for sending and receiving errors.
type Bucket struct {
	name string
	log  *wlog.Logger
	cli  *s3.Client

	objectQueue *objectQueue
	errorsCh    chan error
}

// objectQueue is a type used to represent a queue of objects.
//
// items is a channel used for sending and receiving objects of type *types.Object.
//
// counter is an unsigned 64-bit integer used to keep track of the number of objects in the queue.
type objectQueue struct {
	items   chan string
	counter uint64
}

func New(log *wlog.Logger, cfg *config.S3) *Bucket {
	cli := awshttp.NewBuildableClient()
	cli.WithTransportOptions(func(transport *http.Transport) {
		transport.MaxIdleConns = 1000
		transport.IdleConnTimeout = 90 * time.Second
	})

	awsc := aws.Config{
		Region:           strings.ToLower(cfg.Region),
		Credentials:      credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		RetryMode:        aws.RetryModeStandard,
		RetryMaxAttempts: 3,
		HTTPClient:       cli,
		DefaultsMode:     aws.DefaultsModeCrossRegion,
	}

	return &Bucket{
		name:     cfg.Bucket,
		log:      log,
		cli:      s3.NewFromConfig(awsc),
		errorsCh: make(chan error),
		objectQueue: &objectQueue{
			items:   make(chan string, MaxListKeys),
			counter: 0,
		},
	}
}

func (b *Bucket) EnqueueObjectPool(obj string) {
	atomic.AddUint64(&b.objectQueue.counter, 1)
	b.objectQueue.items <- obj
}

func (b *Bucket) DequeueObjectPool() {
	atomic.AddUint64(&b.objectQueue.counter, ^uint64(0))
}

func (b *Bucket) ObjectPool() chan string {
	return b.objectQueue.items
}

// ListObjects lists the objects in a bucket.
func (b *Bucket) ListObjects(ctx context.Context, key string, lastKey string) error {
	defer close(b.objectQueue.items)
	req := &s3.ListObjectsV2Input{
		Bucket: aws.String(b.name),
		Prefix: aws.String(key),
	}

	if lastKey != "" {
		req.StartAfter = aws.String(lastKey)
	}

	p := s3.NewListObjectsV2Paginator(b.cli, req, func(o *s3.ListObjectsV2PaginatorOptions) {
		if v := int32(MaxListKeys); v != 0 {
			o.Limit = v
		}
	})

	// Iterate through the S3 object pages, printing each object returned.
	var i int
	for p.HasMorePages() {
		i++

		/*if i == 2 {
			return nil
		}*/

		page, err := p.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("get page %d: %v", i, err)
		}

		b.log.Debug("fetched page", wlog.Int("page", i), wlog.Int("len", len(page.Contents)))
		for _, obj := range page.Contents {
			b.EnqueueObjectPool(*obj.Key)
		}
	}

	return nil
}

func (b *Bucket) HeadObject(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	req := &s3.HeadObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(key),
	}

	object, err := b.cli.HeadObject(ctx, req)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// ReadObject gets an object from a bucket and stores it in a local file.
func (b *Bucket) ReadObject(ctx context.Context, key string) ([]byte, error) {
	var i int

start:
	req := &s3.GetObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(key),
	}

	result, err := b.cli.GetObject(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get object: %v", err)
	}

	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		i++
		b.log.Error("failed to read object body, repeat", wlog.Int("len", len(body)), wlog.Err(err),
			wlog.String("key", key), wlog.Int("attempt", i))

		// FIXME
		goto start
		// return nil, fmt.Errorf("read all body: %v", err)
	}

	return body, nil
}

func (b *Bucket) DownloadObject(ctx context.Context, key, filepath string) error {
	object, err := b.ReadObject(ctx, key)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}
	defer file.Close()
	_, err = file.Write(object)
	if err != nil {
		return fmt.Errorf("write: %v", err)
	}

	return nil
}
