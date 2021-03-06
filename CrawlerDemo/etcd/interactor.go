package etcd

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/etcdserver/api/v3client"
	"github.com/pantskun/commonutils/pathutils"
	"github.com/pantskun/golearn/CrawlerDemo/xlogutil"
)

const timeoutSecond = 1000.0

type Interactor interface {
	Close()
	Get(key string) (string, error)
	Put(key string, value string) error
	Del(key string) error
	Lock() (context.CancelFunc, error)
	Unlock() (context.CancelFunc, error)
	TxnSync(key string) bool
}

type interactor struct {
	e *embedetcd
	c *clientv3.Client
	s *concurrency.Session
	m *concurrency.Mutex
}

var _ Interactor = (*interactor)(nil)

type InteractorError struct {
	msg string
}

func (err *InteractorError) Error() string {
	return err.msg
}

// NewInteractorWithEmbed
// 获取与嵌入etcd交互的Interactor
func NewInteractorWithEmbed() (Interactor, error) {
	e, err := newEmbedetcd()
	if e == nil {
		return nil, err
	}

	c := v3client.New(e.etcd.Server)

	// new seesion, new mutex
	// ctx, cancel := context.WithTimeout(context.TODO(), timeoutSecond*time.Second)
	// defer cancel()

	s, ce := concurrency.NewSession(c /*concurrency.WithContext(ctx),*/, concurrency.WithTTL(timeoutSecond))
	if ce != nil {
		return nil, ce
	}

	m := concurrency.NewMutex(s, "/crawler-lock/")

	return &interactor{e, c, s, m}, nil
}

// NewInteractor
// 获取与etcd交互的Interactor
func NewInteractor() (Interactor, error) {
	configPath := pathutils.GetModulePath("CrawlerDemo") + "/configs/etcdConfig.json"
	config := GetClientConfig(configPath)

	c, err := clientv3.New(config)
	if err != nil {
	}

	// new seesion, new mutex
	// ctx, cancel := context.WithTimeout(context.TODO(), timeoutSecond*time.Second)
	// defer cancel()

	s, ce := concurrency.NewSession(c /*concurrency.WithContext(ctx),*/, concurrency.WithTTL(timeoutSecond))
	if ce != nil {
		return nil, ce
	}

	m := concurrency.NewMutex(s, "/crawler-lock/")

	return &interactor{nil, c, s, m}, nil
}

func (i *interactor) Close() {
	if i.e != nil {
		i.e.close()
	}

	i.c.Close()
	i.s.Close()
}

func (i *interactor) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeoutSecond*time.Second)
	defer cancel()

	rsp, ge := i.c.Get(ctx, key)
	if ge != nil {
		return "", ge
	}

	if len(rsp.Kvs) == 0 {
		return "", nil
	}

	return string(rsp.Kvs[0].Value), nil
}

func (i *interactor) Put(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), timeoutSecond*time.Second)
	defer cancel()

	_, err := i.c.Put(ctx, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (i *interactor) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), timeoutSecond*time.Second)
	defer cancel()

	_, err := i.c.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func (i *interactor) Lock() (context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeoutSecond*time.Second)

	err := i.m.Lock(ctx)
	if err != nil {
		return cancel, err
	}

	return cancel, nil
}

func (i *interactor) Unlock() (context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeoutSecond*time.Second)

	err := i.m.Unlock(ctx)
	if err != nil {
		return cancel, err
	}

	return cancel, nil
}

func (i *interactor) TxnSync(key string) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), timeoutSecond*time.Second)
	defer cancel()

	txn := i.c.Txn(ctx)
	txn = txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
	txn = txn.Then(clientv3.OpPut(key, "1"))
	resp, err := txn.Commit()

	if err != nil {
		xlogutil.Error(err)
		return false
	}

	return resp.Succeeded
}
