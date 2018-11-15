package etcdutil

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/coreos/etcd/embed"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

const (
	outputDir    = "../../test/output"
	etcdDir      = outputDir + "/default.etcd"
	etcdEndpoint = "http://localhost:2379"
)

var (
	etcd *embed.Etcd
	err  error
)

func TestEtcdutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Etcdutil Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	err = os.RemoveAll(outputDir)
	Expect(err).ShouldNot(HaveOccurred())

	etcd, err = startEmbeddedEtcd()
	Expect(err).ShouldNot(HaveOccurred())
	var data []byte
	return data
}, func(data []byte) {})

var _ = SynchronizedAfterSuite(func() {}, func() {
	etcd.Server.Stop()
	etcd.Close()
})

func startEmbeddedEtcd() (*embed.Etcd, error) {
	logger := logrus.New()
	logger.Infof("Starting embedded etcd")
	cfg := embed.NewConfig()
	cfg.Dir = etcdDir
	cfg.EnableV2 = false
	cfg.Debug = false
	cfg.GRPCKeepAliveTimeout = 0
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil, err
	}

	select {
	case <-e.Server.ReadyNotify():
		fmt.Printf("Embedded server is ready!\n")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		e.Close()
		return nil, fmt.Errorf("server took too long to start")
	}
	return e, nil
}
