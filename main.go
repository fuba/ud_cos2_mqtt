package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/tarm/serial"
)

type SensorData struct {
	Time        int64   `json:"time"`
	CO2PPM      int     `json:"co2ppm"`
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

func (s SensorData) toString() string {
	return fmt.Sprintf("Time: %d, CO2PPM: %d, Humidity: %.2f, Temperature: %.2f", s.Time, s.CO2PPM, s.Humidity, s.Temperature)
}

var currentData SensorData
var running = true

func main() {
	mqttServer := flag.String("h", "localhost", "MQTT server hostname")
	mqttPort := flag.Int("p", 1883, "MQTT server port")
	mqttUsername := flag.String("u", "", "MQTT username")
	mqttPassword := flag.String("P", "", "MQTT password")
	flag.Parse()

	clientID := uuid.New().String()

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", *mqttServer, *mqttPort)).
		SetClientID(clientID).
		SetUsername(*mqttUsername).
		SetPassword(*mqttPassword)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error connecting to MQTT server:", token.Error())
		return
	}
	defer client.Disconnect(250)

	go startSerial("/dev/ttyACM0", client)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	go func() {
		for {
			sig := <-sigChan
			switch sig {
			case syscall.SIGHUP:
				fmt.Println(sig.String() + " signal received, stopping...")
				running = false
				return
			}
		}
	}()

	for running {
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Program stopped")
}

func startSerial(dev string, client mqtt.Client) {
	cfg := &serial.Config{Name: dev, Baud: 115200, ReadTimeout: time.Second * 6}
	s, err := serial.OpenPort(cfg)
	if err != nil {
		fmt.Println("Error opening serial port:", err)
		return
	}
	defer s.Close()

	_, err = s.Write([]byte("STA\r\n"))
	if err != nil {
		fmt.Println("Error writing to serial port:", err)
		return
	}
	buf := make([]byte, 128)
	_, err = s.Read(buf)
	if err != nil {
		fmt.Println("Error reading from serial port:", err)
		return
	}

	regex := regexp.MustCompile(`CO2=(?P<co2>\d+),HUM=(?P<hum>\d+\.\d+),TMP=(?P<tmp>-?\d+\.\d+)`)
	for running {
		n, err := s.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port:", err)
			}
			break
		}
		str := string(buf[:n])
		matches := regex.FindStringSubmatch(str)
		if matches != nil {
			co2, _ := strconv.Atoi(matches[1])
			hum, _ := strconv.ParseFloat(matches[2], 64)
			tmp, _ := strconv.ParseFloat(matches[3], 64)
			currentData = SensorData{
				Time:        time.Now().Unix(),
				CO2PPM:      co2,
				Humidity:    hum,
				Temperature: tmp,
			}

			fmt.Println(currentData.toString())

			jsonData, err := json.Marshal(currentData)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				continue
			}

			token := client.Publish("homeassistant/ud_cos2", 0, false, jsonData)
			token.Wait()
			if token.Error() != nil {
				fmt.Println("Error publishing to MQTT:", token.Error())
			}
		}
	}
	_, err = s.Write([]byte("STP\r\n"))
	if err != nil {
		fmt.Println("Error writing to serial port:", err)
	}
}
