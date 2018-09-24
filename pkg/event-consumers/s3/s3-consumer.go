package s3

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	v1beta1 "github.com/epimorphics/s3-trigger/pkg/apis/kubeless/v1beta1"
	"github.com/epimorphics/s3-trigger/pkg/client/clientset/versioned"
	"github.com/epimorphics/s3-trigger/pkg/utils"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

var (
	stopM     map[string](chan struct{})
	stoppedM  map[string](chan struct{})
	consumerM map[string]bool

	svc          *s3.S3
	PolledFormat string
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
	PolledFormat = "2006-01-02T15:04:05.000Z"
}

func UpdateS3TriggerTime(s3interface versioned.Interface, ns, triggerObjName string, timestamp time.Time) {
	triggerObj, err := s3interface.KubelessV1beta1().S3Triggers(ns).Get(triggerObjName, metav1.GetOptions{})
	if err != nil {
		logrus.Fatalf("Could not get trigger %s from kubernetes API, status not updated: %s", triggerObjName, err)
		return
	}
	logrus.Infof("Updating trigger %s status", triggerObjName)
	triggerClone := triggerObj.DeepCopy()
	triggerClone.Status = v1beta1.S3TriggerStatus{LastPolled: timestamp.UTC().Format(PolledFormat)}
	err = utils.UpdateKafkaTriggerCustomResource(s3interface, triggerClone)
	if err != nil {
		logrus.Fatalf("could not update trigger %s: %s", triggerObjName, err)
		return
	}
	logrus.Infof("Updated trigger %s status", triggerObjName)
}

func GetS3TriggerTime(s3interface versioned.Interface, ns, triggerObjName string) *time.Time {
	triggerObj, err := s3interface.KubelessV1beta1().S3Triggers(ns).Get(triggerObjName, metav1.GetOptions{})
	if err != nil {
		logrus.Fatalf("Could not get trigger %s metadata from kubernetes API: %s", triggerObjName, err)
		return nil
	}
	time, err := time.Parse(PolledFormat, triggerObj.Status.LastPolled)
	if err != nil {
		logrus.Errorf("Could not parse last poll time for trigger %s from kubernetes API: %s", triggerObjName, err)
		return nil
	}
	return &time
}

func S3ConsumerProcess(bucket, prefix, triggerObjName, funcName, ns string, pollFrequency int, clientset kubernetes.Interface, s3interface versioned.Interface, stopchan, stoppedchan chan struct{}) {
	defer close(stoppedchan)
	consumer := make(chan ObjectMetadata)
	pollUpdate := make(chan time.Time)
	closeWatcher := make(chan struct{})
	go Watcher(bucket, prefix, pollFrequency, func() *time.Time { return GetS3TriggerTime(s3interface, ns, triggerObjName) }, pollUpdate, consumer, closeWatcher)
	defer close(closeWatcher)
	defer close(consumer)
	defer close(pollUpdate)
	for {
		select {
		case msg := <-pollUpdate:
			go func() {
				UpdateS3TriggerTime(s3interface, ns, triggerObjName, msg)
			}()
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

func CreateS3Consumer(triggerObjName, funcName, ns, bucket, key string, pollFrequency int, clientset kubernetes.Interface, s3interface versioned.Interface) error {
	consumerID := generateUniqueConsumerGroupID(triggerObjName, funcName, ns, bucket, key)
	if !consumerM[consumerID] {
		logrus.Infof("Creating S3 consumer for the function %s associated with trigger %s", funcName, triggerObjName)
		stopM[consumerID] = make(chan struct{})
		stoppedM[consumerID] = make(chan struct{})
		go S3ConsumerProcess(bucket, key, triggerObjName, funcName, ns, pollFrequency, clientset, s3interface, stopM[consumerID], stoppedM[consumerID])
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
