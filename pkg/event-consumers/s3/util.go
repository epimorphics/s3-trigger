package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"time"
)

type ObjectMetadata struct {
	Bucket   string    `json:"bucket"`
	Key      string    `json:"key"`
	Modified time.Time `json:"modified"`
}

func Watcher(bucket string, key string, pollDuration int, etcdPolled func() *time.Time, polled chan time.Time, data chan ObjectMetadata, closechan chan struct{}) {
	input := &s3.ListObjectsInput{Bucket: aws.String(bucket), Prefix: aws.String(key)}
	etcdStart := etcdPolled()
	var lastPolled time.Time
	if etcdStart != nil {
		lastPolled = *etcdStart
	} else {
		lastPolled = time.Now()
	}
	for {
		select {
		case <-closechan:
			return
		default:
			metadata := since(input, lastPolled)
			lastPolled = time.Now()
			for _, object := range metadata {
				data <- object
			}
			polled <- lastPolled
			time.Sleep(time.Duration(pollDuration) * time.Second)
		}
	}
}

func since(input *s3.ListObjectsInput, sinceTime time.Time) []ObjectMetadata {
	resp, err := svc.ListObjects(input)
	if err != nil {
		logrus.Errorf("Unable to list items in bucket %v, %v", input.Bucket, err)
	}

	objects := make([]ObjectMetadata, 0)
	logrus.Infof("Polled s3 bucket %v, %v", input.Bucket, input.Prefix)
	for _, item := range resp.Contents {
		if item.LastModified.After(sinceTime) {
			logrus.Infof("Item Discovered")
			object := ObjectMetadata{
				Bucket:   *input.Bucket,
				Key:      *item.Key,
				Modified: *item.LastModified,
			}
			logrus.Infof("%v", object)
			objects = append(objects, object)
		}
	}
	return objects
}
