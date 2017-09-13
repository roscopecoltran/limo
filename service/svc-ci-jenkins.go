package service

/*
import (
	"io/ioutil"

	"github.com/bndr/gojenkins"
	"github.com/sirupsen/logrus"
)

// Jenkins holds the information about one or more Jenkins servers that the bot
// should send or retrieve information from.
type Jenkins struct {
	URL            string             `json:"jenkinsURL"`
	User           string             `json:"jenkinsUser"`
	Token          string             `json:"jenkinsToken"`
	CACertFilePath string             `json:"caCertFilePath"`
	Client         *gojenkins.Jenkins `json:"-"`
	Jobs           []JenkinsJob       `json:"-"`
}

// JenkinsJob holds the information about a job in jenkins.
type JenkinsJob struct{}

func newJenkinsInstance() (*Jenkins, error) {
	return nil, nil
}

// Start creates a new Jenkins instance.
func (j *Jenkins) Start() error {
	client := gojenkins.CreateJenkins(j.URL, j.User, j.Token)

	if j.CACertFilePath != "" {
		caCert, err := ioutil.ReadFile(j.CACertFilePath)
		if err != nil {
			return err
		}
		if len(caCert) == 0 {
			logrus.Warnf("Specified CA Certificate file (%s) is empty. Using unencrypted connection", j.CACertFilePath)
		} else {
			client.Requester.CACert = caCert
		}
	}
	instance, err := client.Init()

	if err != nil {
		return err
	}

	logrus.Debugf("This is the jenkins instance: %#v", instance)
	return nil
}

// Stop will gracefully stop a Jenkins instance
func (j *Jenkins) Stop() error { return nil }

// StartJob starts a new JenkinsJob
func (j *JenkinsJob) StartJob() error { return nil }

// CancelJob will gracefully cancel a job in Jenkins.
func (j *JenkinsJob) CancelJob() error { return nil }
*/