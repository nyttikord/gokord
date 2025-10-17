package event

import (
	"context"
	"errors"

	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/logger"
)

var (
	ErrInvalidHandler = errors.New("invalid handler type")
)

// Handler is an interface for Discord events.
type Handler interface {
	// Type returns the type of event this handler belongs to.
	Type() string

	// Handle is called whenever an event of Type() happens.
	// It is the receivers responsibility to type assert that the interface
	// is the expected struct.
	Handle(context.Context, bot.Session, any)
}

// InterfaceProvider is an interface for providing empty interfaces for
// Discord events.
type InterfaceProvider interface {
	// Type is the type of event this handler belongs to.
	Type() string

	// New returns a new instance of the struct this event handler handles.
	// This is called once per event.
	// The struct is provided to all handlers of the same Type.
	New() any
}

// interfaceEventType is the event handler type for any events.
const interfaceEventType = "__INTERFACE__"

// interfaceHandler is an event handler for any events.
type interfaceHandler func(context.Context, bot.Session, any)

// Type returns the event type for any events.
func (eh interfaceHandler) Type() string {
	return interfaceEventType
}

// Handle is the handler for an any event.
func (eh interfaceHandler) Handle(ctx context.Context, s bot.Session, i any) {
	eh(ctx, s, i)
}

var registeredInterfaceProviders = map[string]InterfaceProvider{}

// registerInterfaceProvider registers a provider so that Gokord can access it's New() method.
func registerInterfaceProvider(eh InterfaceProvider) {
	if _, ok := registeredInterfaceProviders[eh.Type()]; ok {
		return
		// XXX:
		// if we should error here, we need to do something with it.
		// fmt.Errorf("event %s already registered", eh.Type())
	}
	registeredInterfaceProviders[eh.Type()] = eh
}

// GetInterfaceProvider returns the InterfaceProvider and true for the given type.
// It returns nil and false otherwise.
func GetInterfaceProvider(typ string) (InterfaceProvider, bool) {
	in, ok := registeredInterfaceProviders[typ]
	return in, ok
}

// eventHandlerInstance is a wrapper around an event handler, as functions cannot be compared directly.
type eventHandlerInstance struct {
	eventHandler Handler
}

// addEventHandler adds an event handler that will be fired anytime the Discord WSAPI matching Handler.Type fires.
func (e *Manager) addEventHandler(eventHandler Handler) func() {
	e.Lock()
	defer e.Unlock()

	if e.handlers == nil {
		e.handlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	e.handlers[eventHandler.Type()] = append(e.handlers[eventHandler.Type()], ehi)

	return func() {
		e.removeEventHandlerInstance(eventHandler.Type(), ehi)
	}
}

// addEventHandler adds an event handler that will be fired the next time
// the Discord WSAPI matching Handler.Type fires.
func (e *Manager) addEventHandlerOnce(eventHandler Handler) func() {
	e.Lock()
	defer e.Unlock()

	if e.onceHandlers == nil {
		e.onceHandlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	e.onceHandlers[eventHandler.Type()] = append(e.onceHandlers[eventHandler.Type()], ehi)

	return func() {
		e.removeEventHandlerInstance(eventHandler.Type(), ehi)
	}
}

// AddHandler allows you to add an event handler that will be fired anytime the Discord WSAPI event that matches the
// function fires.
// The first parameter is a Session, and the second parameter is a pointer to a struct corresponding to the event for
// which you want to listen.
//
// eg:
//
//	Session.AddHandler(func(s event.Session, m *discordgo.MessageCreate) {
//	})
//
// or:
//
//	Session.AddHandler(func(s event.Session, m *discordgo.PresenceUpdate) {
//	})
//
// List of events can be found at this page, with corresponding names in the library for each event:
// https://discord.com/developers/docs/topics/gateway#event-names
// There are also synthetic events fired by the library internally which are available for handling, like Connect,
// Disconnect, and RateLimit.
// events.go contains all the Discord WSAPI and synthetic events that can be handled.
//
// The return value of this method is a function, that when called will remove the event handler.
func (e *Manager) AddHandler(handler any) func() {
	eh := handlerForInterface(handler)

	if eh == nil {
		e.Logger().ErrorContext(
			logger.NewContext(context.Background(), 1),
			"handler will never be called",
			"error", ErrInvalidHandler,
		)
		return func() {}
	}

	return e.addEventHandler(eh)
}

// AddHandlerOnce allows you to add an event handler that will be fired the next time the Discord WSAPI event that
// matches the function fires.
//
// See AddHandler for more details.
func (e *Manager) AddHandlerOnce(handler any) func() {
	eh := handlerForInterface(handler)

	if eh == nil {
		e.Logger().ErrorContext(
			logger.NewContext(context.Background(), 1),
			"handler will never be called",
			"error", ErrInvalidHandler,
		)
		return func() {}
	}

	return e.addEventHandlerOnce(eh)
}

// removeEventHandlerInstance removes an event handler instance.
func (e *Manager) removeEventHandlerInstance(t string, ehi *eventHandlerInstance) {
	e.Lock()
	defer e.Unlock()

	handlers := e.handlers[t]
	for i := range handlers {
		if handlers[i] == ehi {
			e.handlers[t] = append(handlers[:i], handlers[i+1:]...)
		}
	}

	onceHandlers := e.onceHandlers[t]
	for i := range onceHandlers {
		if onceHandlers[i] == ehi {
			e.onceHandlers[t] = append(onceHandlers[:i], onceHandlers[i+1:]...)
		}
	}
}

// Handles calling permanent and once handlers for an event type.
func (e *Manager) handle(ctx context.Context, s bot.Session, t string, i any) {
	for _, eh := range e.handlers[t] {
		if e.SyncEvents {
			eh.eventHandler.Handle(ctx, s, i)
		} else {
			go eh.eventHandler.Handle(ctx, s, i)
		}
	}

	if len(e.onceHandlers[t]) > 0 {
		for _, eh := range e.onceHandlers[t] {
			if e.SyncEvents {
				eh.eventHandler.Handle(ctx, s, i)
			} else {
				go eh.eventHandler.Handle(ctx, s, i)
			}
		}
		e.onceHandlers[t] = nil
	}
}

// EmitEvent calls internal methods, fires handlers and fires the "any" event.
func (e *Manager) EmitEvent(ctx context.Context, s bot.Session, t string, i any) {
	e.RLock()
	defer e.RUnlock()

	// All events are dispatched internally first.
	e.onInterface(ctx, i)

	// Then they are dispatched to anyone handling any events.
	e.handle(ctx, s, interfaceEventType, i)

	// Finally they are dispatched to any typed handlers.
	e.handle(ctx, s, t, i)
}
