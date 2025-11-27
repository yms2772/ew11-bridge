package ew11

import (
	"github.com/eclipse/paho.mqtt.golang"
)

type Device interface {
	On() error
	Off() error
	SetLevel(int) error
	SetStatus([]byte) error
	IsDevice([]byte) bool
	Init() error
}

type hasDeviceBase interface {
	Base() *DeviceBase
}

type DeviceBase struct {
	c *Communicator
}

func (d *DeviceBase) IsDebug() bool {
	return d.c.debug
}

func (d *DeviceBase) PublishToCustomTopic(topic string, payload any, retained bool) error {
	return d.c.client.Publish(topic, 1, retained, payload).Error()
}

func (d *DeviceBase) PublishToSendTopic(bridge Bridge, payload any, retained bool) error {
	return d.PublishToCustomTopic(d.c.topicMap[bridge].send, payload, retained)
}

func (d *DeviceBase) SubscribeFromCustomTopic(topic string, callback mqtt.MessageHandler) error {
	return d.c.client.Subscribe(topic, 1, callback).Error()
}

func (d *DeviceBase) setCommunicator(c *Communicator) {
	d.c = c
}
