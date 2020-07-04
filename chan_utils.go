package chan_utils

import (
	"context"
	"fmt"
)

type Message interface {
	GetMessage(data interface{})
}

func NewListener(ctx context.Context, ch chan Message, onMessage func(Message), onError func(error)) (func(), context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	return func() {
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

func GetMessage(ctx context.Context, ch chan Message) (Message, error) {
	var ok bool
	var err error
	var msg Message

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
