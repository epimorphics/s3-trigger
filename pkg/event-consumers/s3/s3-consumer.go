package s3

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/epimorphics/s3-trigger/pkg/utils"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"time"
)

var (
	stopM     map[string](chan struct{})
	stoppedM  map[string](chan struct{})
	consumerM map[string]bool

	svc *s3.S3
)

type WatcherDataInstance struct {
	Bucket      string    `json:"bucket"`
	Key         string    `json:"key"`
	LastChecked time.Time `json:"lastChecked"`
}

type ObjectData struct {
	Data     string    `json:"data"`
	Bucket   string    `json:"bucket"`
	Key      string    `json:"key"`
	Modified time.Time `json:"modified"`
}

func init() {
	stopM = make(map[string](chan struct{}))
	stoppedM = make(map[string](chan struct{}))
	consumerM = make(map[string]bool)
	sess, err := session.NewSession(
		&aws.Config{Region: aws.String("eu-west-1")},
	)
	if err != nil {
		logrus.Fatalf("Couldn't create s3 service: %v", err)
	}
	// Create S3 service client
	svc = s3.New(sess)
}

func S3ConsumerProcess(bucket, prefix, funcName, ns string, pollFrequency int, clientset kubernetes.Interface, stopchan, stoppedchan chan struct{}) {
	defer close(stoppedchan)
	consumer := make(chan ObjectMetadata)
	closeWatcher := make(chan struct{})
	go Watcher(bucket, prefix, pollFrequency, consumer, closeWatcher)
	defer close(closeWatcher)
	defer close(consumer)
	for {
		select {
		case msg := <-consumer:
			go func() {
				jsonString, err := json.Marshal(msg)
				if err != nil {
					logrus.Errorf("Unable to marshal s3 message: %v", err)
					return
				}
				req, err := utils.GetHTTPReq(clientset, funcName, ns, "s3triggers.kubeless.io", "POST", string(jsonString))
				if err != nil {
					logrus.Errorf("Unable to elaborate request: %v", err)
				} else {
					err = utils.SendMessage(req)
					if err != nil {
						logrus.Errorf("Failed to send message to function: %v", err)
					} else {
						logrus.Infof("Message was sent to function %s successfully", funcName)
					}
				}
			}()
		case <-stopchan:
			return
		}
	}
}

func CreateS3Consumer(triggerObjName, funcName, ns, bucket, key string, pollFrequency int, clientset kubernetes.Interface) error {
	consumerID := generateUniqueConsumerGroupID(triggerObjName, funcName, ns, bucket, key)
	if !consumerM[consumerID] {
		logrus.Infof("Creating S3 consumer for the function %s associated with trigger %s", funcName, triggerObjName)
		stopM[consumerID] = make(chan struct{})
		stoppedM[consumerID] = make(chan struct{})
		go S3ConsumerProcess(bucket, key, funcName, ns, pollFrequency, clientset, stopM[consumerID], stoppedM[consumerID])
		consumerM[consumerID] = true
		logrus.Infof("Created S3 consumer for the function %s associated with trigger %s", funcName, triggerObjName)
	} else {
		logrus.Infof("Consumer for function %s associated with trigger %s already exists", funcName, triggerObjName)
	}
	return nil
}

func DeleteS3Consumer(triggerObjName, funcName, ns, bucket, key string) error {
	consumerID := generateUniqueConsumerGroupID(triggerObjName, funcName, ns, bucket, key)
	if consumerM[consumerID] {
		logrus.Infof("Stopping consumer for the function %s associated with trigger %s", funcName, triggerObjName)
		close(stopM[consumerID])
		<-stoppedM[consumerID]
		consumerM[consumerID] = false
		logrus.Infof("Stopped consumer for the function %s associated with trigger %s", funcName, triggerObjName)
	} else {
		logrus.Infof("Consumer for the function %s associated with trigger %s doesn't exist", funcName, triggerObjName)
	}
	return nil
}

func generateUniqueConsumerGroupID(triggerObjName, funcName, ns, bucket, key string) string {
	return ns + "_" + triggerObjName + "_" + funcName + "_" + bucket + "_" + key
}

/*
func since(input *s3.ListObjectsInput, sinceTime time.Time, data chan ObjectData) {
	resp, err := svc.ListObjects(input)
	if err != nil {
		log.Printf("Unable to list items in bucket %q, %v", input.Bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println(item.LastModified, sinceTime)
		if item.LastModified.After(sinceTime) {
			fmt.Println("Name:         ", *item.Key)
			fmt.Println("Last modified:", *item.LastModified)
			fmt.Println("Size:         ", *item.Size)
			fmt.Println("Storage class:", *item.StorageClass)
			fmt.Println("")
			result, err := svc.GetObject(&s3.GetObjectInput{
				Bucket: input.Bucket,
				Key:    item.Key,
			})
			if err != nil {
				log.Println(err)
			}
			file, err := ioutil.ReadAll(result.Body)
			if err != nil {
				log.Println(err)
			}
			fileData := ObjectData{
				Data:     string(file),
				Bucket:   *input.Bucket,
				Key:      *item.Key,
				Modified: *item.LastModified,
			}
			data <- fileData
			log.Println(fileData)
		}
	}
}*/
