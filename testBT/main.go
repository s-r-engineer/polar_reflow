package testBT

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	influxclient "polar_reflow/database/influxClient"
	"polar_reflow/models"
	"polar_reflow/tools"
	"time"
	"tinygo.org/x/bluetooth"
)

var err error
var file *os.File

func RunBT() {
	uid := uuid.New()
	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		log.Fatalf("Failed to enable Bluetooth adapter: %v", err)
	}
	deviceMAC := "A0:9E:1A:BF:02:06"
	bb := bluetooth.Address{}
	bb.Set(deviceMAC)
	var device bluetooth.Device
	for {
		device, err = adapter.Connect(bb, bluetooth.ConnectionParams{})
		if err == nil {
			break
		} else {
			log.Print("Failed to connect to device: %v", err)
		}
		time.Sleep(time.Second)
	}

	defer device.Disconnect()
	services, err := device.DiscoverServices(nil)
	if err != nil {
		log.Fatalf("Failed to discover services: %v", err)
	}
	for _, service := range services {
		log.Println("service")
		if service.String() == "0000180d-0000-1000-8000-00805f9b34fb" {
			log.Println("right service")
			characteristics, _ := service.DiscoverCharacteristics(nil)
			for _, char := range characteristics {
				log.Println("right char")
				time.Sleep(time.Second * 5)
				log.Println("right char 2")
				err = char.EnableNotifications(func(f []byte) {
					t := time.Now()
					result := binary.BigEndian.Uint16(f)
					write(models.BTHR{Session: uid, Value: result, TimePoint: t})
				})
				tools.Dumper(err)
			}
		}
	}
	time.Sleep(time.Hour * 2)
}

func init() {
	file, err = os.OpenFile("./data", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
}

func write(bthr models.BTHR) {
	tools.Dumper(bthr)
	_, err = file.Write([]byte(fmt.Sprintf("%d %s %d\n", bthr.TimePoint.UnixNano(), bthr.Session.String(), bthr.Value)))
	influxclient.WriteBTHR(bthr)
	if err != nil {
		tools.Dumper(err)
	}
}
