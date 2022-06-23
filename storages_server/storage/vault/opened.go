package vault

import "sync"

type OpenedCallback = func()

type Opened struct {
    sync.Mutex

    values      map[string]int
    callbacks   map[string]OpenedCallback
}

func NewOpened() *Opened {
    return &Opened{
        values      : make(map[string]int),
        callbacks   : make(map[string]OpenedCallback),
    }
}


func (o *Opened) Open(h string) {
    o.Lock()
    defer o.Unlock()

    _, ok := o.values[h]
    if !ok {
        o.values[h] = 1
    } else {
        o.values[h] += 1
    }
}

func (o *Opened) Close(h string) {
    o.Lock()
    defer o.Unlock()

    _, ok := o.values[h]
    if !ok {
        panic("Not balanced object state!")
    }

    o.values[h] -= 1
    if o.values[h] == 0 {
        if callback, ok := o.callbacks[h]; ok {
            callback()
            delete(o.callbacks, h)
        }
        delete(o.values, h)
    }
}

func (o *Opened) IsOpen(h string) bool {
    o.Lock()
    defer o.Unlock()

    _, ok := o.values[h]
    if !ok {
        return false
    } else {
        return o.values[h] > 0
    }
}

func (o *Opened) OnClose(h string, callback OpenedCallback) {
    o.Lock()
    defer o.Unlock()

    o.callbacks[h] = callback
}

func (o *Opened) Cancel(h string) {
    o.Lock()
    defer o.Unlock()

    if _, ok := o.callbacks[h]; ok {
        delete(o.callbacks, h)
    }
}
