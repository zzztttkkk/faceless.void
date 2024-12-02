package evts

import "reflect"

type EvtOnUserCreated struct {
	Uid string
}

var (
	typeOfOnUserCreated = reflect.TypeOf(EvtOnUserCreated{})
)

func EmitOnUserCreated(evt EvtOnUserCreated) {
	bus.Emit(typeOfOnUserCreated, evt)
}

func OnUserCreated(fnc func(evt EvtOnUserCreated)) {
	bus.AddListener(typeOfOnUserCreated, func(at int64, evtany any) {
		fnc(evtany.(EvtOnUserCreated))
	})
}
