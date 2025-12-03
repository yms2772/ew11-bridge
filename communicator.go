package ew11

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type topic struct {
	recv string
	send string
}

type Communicator struct {
	client   mqtt.Client
	debug    bool
	devices  []Device
	topicMap map[Bridge]topic
	home     struct {
		building int
		unit     int
	}
}

func NewCommunicator() (*Communicator, error) {
	var c Communicator
	c.debug, _ = strconv.ParseBool(os.Getenv("EW11_DEBUG"))
	c.topicMap = map[Bridge]topic{
		Bridge1: {
			recv: os.Getenv("EW11_RECEIVE_TOPIC"),
			send: os.Getenv("EW11_SEND_TOPIC"),
		},
		Bridge2: {
			recv: os.Getenv("EW11_2_RECEIVE_TOPIC"),
			send: os.Getenv("EW11_2_SEND_TOPIC"),
		},
	}

	building, err := strconv.Atoi(os.Getenv("EW11_HOME_BUILDING_NUMBER"))
	if err != nil {
		return nil, fmt.Errorf("동 번호가 잘못되었습니다: %v", err)
	}

	unit, err := strconv.Atoi(os.Getenv("EW11_HOME_UNIT_NUMBER"))
	if err != nil {
		return nil, fmt.Errorf("호 번호가 잘못되었습니다: %v", err)
	}

	c.home.building = building
	c.home.unit = unit

	c.client = mqtt.NewClient(mqtt.NewClientOptions().
		AddBroker(os.Getenv("EW11_MQTT_BROKER_URL")).
		SetClientID(os.Getenv("EW11_MQTT_BROKER_CLIENT_ID")).
		SetUsername(os.Getenv("EW11_MQTT_BROKER_USERNAME")).
		SetPassword(os.Getenv("EW11_MQTT_BROKER_PASSWORD")).
		SetCleanSession(true).
		SetOrderMatters(false).
		SetKeepAlive(30 * time.Second))
	tok := c.client.Connect()
	if ok := tok.WaitTimeout(5 * time.Second); !ok {
		return nil, fmt.Errorf("타임아웃: MQTT 브로커에 연결할 수 없습니다")
	}
	if err := tok.Error(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Communicator) Disconnect() {
	c.client.Disconnect(250)
}

func (c *Communicator) AddDevice(device Device) error {
	base, ok := device.(hasDeviceBase)
	if !ok {
		return fmt.Errorf("디바이스 추가 실패: 'Base() *DeviceBase'을 포함해야 합니다")
	}

	if b := base.Base(); b != nil {
		b.c = c
	} else {
		return fmt.Errorf("디바이스 추가 실패: 'Base() *DeviceBase'이 nil을 반환했습니다")
	}

	c.devices = append(c.devices, device)
	return device.Init()
}

func (c *Communicator) StartAndWait() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	callback := func(client mqtt.Client, message mqtt.Message) {
		payload := message.Payload()
		for _, device := range c.devices {
			if device.IsDevice(payload) {
				if err := device.SetStatus(payload); err != nil && c.debug {
					log.Printf("디바이스 상태 설정 실패: %v (%s)", err, PrettyHex(payload))
				}
			}
		}
	}

	c.client.Subscribe(c.topicMap[Bridge1].recv, 1, callback)
	c.client.Subscribe(c.topicMap[Bridge2].recv, 1, callback)

	<-ctx.Done()
}
