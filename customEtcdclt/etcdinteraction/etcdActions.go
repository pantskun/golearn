package etcdinteraction

import (
	"context"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// etcd action type.
const (
	EtcdActGet    = 1
	EtcdActPut    = 2
	EtcdActDelete = 3
)

const (
	ONE = 1
	TWO = 2
)

const timeoutSecond = 5.0

// EtcdActionInterface etcd action interface.
type EtcdActionInterface interface {
	Equal(EtcdActionInterface) bool
	Exec(client *clientv3.Client) ([]string, error)
}

// EtcdActionGet etcd get action info.
type EtcdActionGet struct {
	ActionType int
	Key        string
	RangeEnd   string
}

// EtcdActionPut etcd put action info.
type EtcdActionPut struct {
	ActionType int
	Key        string
	Value      string
}

// EtcdActionDelete etcd delete action info.
type EtcdActionDelete struct {
	ActionType int
	Key        string
	RangeEnd   string
}

// NewGetAction new get action.
// args: key, rangeEnd
func NewGetAction(args ...string) EtcdActionInterface {
	key := ""
	rangeEnd := ""

	switch len(args) {
	case ONE:
		key = args[0]
	case TWO:
		key = args[0]
		rangeEnd = args[1]
	}

	return &EtcdActionGet{
		ActionType: EtcdActGet,
		Key:        key,
		RangeEnd:   rangeEnd,
	}
}

// NewPutAction new put action.
// args: key, value
func NewPutAction(args ...string) EtcdActionInterface {
	key := ""
	value := ""

	switch len(args) {
	case ONE:
		key = args[0]
	case TWO:
		key = args[0]
		value = args[1]
	}

	return &EtcdActionPut{
		ActionType: EtcdActPut,
		Key:        key,
		Value:      value,
	}
}

// NewDeleteAction new delete action.
// args: key, rangeEnd
func NewDeleteAction(args ...string) EtcdActionInterface {
	key := ""
	rangeEnd := ""

	switch len(args) {
	case ONE:
		key = args[0]
	case TWO:
		key = args[0]
		rangeEnd = args[1]
	}

	return &EtcdActionDelete{
		ActionType: EtcdActDelete,
		Key:        key,
		RangeEnd:   rangeEnd,
	}
}

// EtcdError etcd error.
type EtcdError struct {
	Msg string
}

func (e EtcdError) Error() string {
	return e.Msg
}

type EtcdActionNilError struct{}

func (e EtcdActionNilError) Error() string {
	return "action is nil"
}

// Equal get action.
func (action *EtcdActionGet) Equal(b EtcdActionInterface) bool {
	v, ok := b.(*EtcdActionGet)
	if !ok || action.Key != v.Key || action.RangeEnd != v.RangeEnd {
		return false
	}

	return true
}

// Equal delete action.
func (action *EtcdActionDelete) Equal(b EtcdActionInterface) bool {
	v, ok := b.(*EtcdActionDelete)
	if !ok || action.Key != v.Key || action.RangeEnd != v.RangeEnd {
		return false
	}

	return true
}

// Equal put action.
func (action *EtcdActionPut) Equal(b EtcdActionInterface) bool {
	v, ok := b.(*EtcdActionPut)
	if !ok || action.Key != v.Key || action.Value != v.Value {
		return false
	}

	return true
}

// Exec execute etcd get action.
func (action EtcdActionGet) Exec(client *clientv3.Client) ([]string, error) {
	if client == nil {
		return nil, EtcdError{"can not connect to etcd"}
	}

	if action.Key == "" {
		return nil, EtcdError{"get command needs one argument as key and an optional argument as range_end"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	var (
		getResp *clientv3.GetResponse
		err     error
	)

	if getResp, err = client.Get(ctx, action.Key); err != nil {
		return nil, err
	}

	result := make([]string, getResp.Count*TWO)
	for i, elem := range getResp.Kvs {
		result[i*2] = string(elem.Key)
		result[i*2+1] = string(elem.Value)
	}

	return result, nil
}

// Exec execute etcd put action.
func (action EtcdActionPut) Exec(client *clientv3.Client) ([]string, error) {
	if client == nil {
		return nil, EtcdError{"can not connect to etcd"}
	}

	if action.Key == "" || action.Value == "" {
		return nil, EtcdError{"put command needs 2 arguments"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := client.Put(ctx, action.Key, action.Value); err != nil {
		return nil, err
	}

	return []string{"OK"}, nil
}

// Exec execute etcd delete action.
func (action EtcdActionDelete) Exec(client *clientv3.Client) ([]string, error) {
	if client == nil {
		return nil, EtcdError{"can not connect to etcd"}
	}

	if action.Key == "" {
		return nil, EtcdError{"del command needs one argument as key and an optional argument as range_end"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	var (
		getResp *clientv3.DeleteResponse
		err     error
	)

	if getResp, err = client.Delete(ctx, action.Key); err != nil {
		return nil, err
	}

	return []string{strconv.FormatInt(getResp.Deleted, 10)}, nil
}

// GetEtcdClient return etcd client.
func GetEtcdClient(config clientv3.Config) *clientv3.Client {
	client, _ := clientv3.New(config)
	return client
}

func ExecuteAction(action EtcdActionInterface, client *clientv3.Client) (string, error) {
	if action == nil {
		err := EtcdActionNilError{}
		return "", err
	}

	if msgs, err := action.Exec(client); err != nil {
		return "", err
	} else {
		result := ""
		for _, msg := range msgs {
			result = result + msg + "\n"
		}
		return result, nil
	}
}
