package chan_utils

import (
	"context"
	"fmt"
)

func NewListener(ctx context.Context) (func(chan interface{}, func(interface{}), func(error)), context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	return func(ch chan interface{}, onMessage func(interface{}), onError func(error)) {
		defer func() {
			if e := recover(); e != nil {
				onError(e.(error))
			}
		}()

		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					onError(fmt.Errorf("channel was closed"))
					return
				}

				onMessage(msg)
			case <-ctx.Done():
				return
			}
		}
	}, cancel
}

func GetMessage(ctx context.Context, ch chan interface{}) (interface{}, error) {
	var ok bool
	var err error
	var msg interface{}
	
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	select {
	case msg, ok = <-ch:
		if !ok {
			err = fmt.Errorf("channel was closed")
		}
	case <-ctx.Done():
		break
	}

	return msg, err
}