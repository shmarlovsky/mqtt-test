package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	TEMP_TOPIC      = "/sensors/temp"
	HUMIDITY_TOPIC  = "/sensors/hum"
	QOS_0           = 0
	QOS_1           = 1
	SERVERADDRESS   = "tcp://localhost:1883"
	NOTIFY_INTERVAL = time.Second

	WRITETOLOG = true // If true then published messages will be written to the console
)

type Sensor struct {
	time string
	Name string
}

func NewSensor() *Sensor {
	name := fmt.Sprintf("Sensor%v", rand.Intn(100))
	return &Sensor{
		Name: name,
	}
}

func (s *Sensor) SetTime(t string) {
	s.time = t
}

func (s *Sensor) Time() string {
	return fmt.Sprintf("%v time: %v", s.Name, s.time)
}

func (s *Sensor) Temperature() string {
	t := rand.Intn(50)
	return fmt.Sprintf("%v temp: %v", s.Name, t)
}

func (s *Sensor) Humidity() string {
	t := rand.Intn(100)
	return fmt.Sprintf("%v humidity: %v", s.Name, t)
}

func clientOptions(clientID string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(SERVERADDRESS)
	opts.SetClientID(clientID)

	opts.SetOrderMatters(
		false,
	) // Allow out of order messages (use this option unless in order delivery is essential)
	opts.ConnectTimeout = time.Second // Minimal delays on connect
	opts.WriteTimeout = time.Second   // Minimal delays on writes
	opts.KeepAlive = 10               // Keepalive every 10 seconds so we quickly detect network outages
	opts.PingTimeout = time.Second    // local broker so response should be quick

	// Automate connection management (will keep trying to connect and will reconnect if network drops)
	opts.ConnectRetry = true
	opts.AutoReconnect = true

	// Log events
	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		fmt.Println("connection lost")
	}
	opts.OnConnect = func(mqtt.Client) {
		fmt.Println("connection established")
	}
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		fmt.Println("attempting to reconnect")
	}
	return opts
}

func main() {
	// Enable logging by uncommenting the below
	// mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	// mqtt.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
	// mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	sensor := NewSensor()

	opts := clientOptions(sensor.Name)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connection is up")

	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		var count uint64
		for {
			select {
			case <-time.After(NOTIFY_INTERVAL):
				count += 1
				msg := sensor.Temperature()

				if WRITETOLOG {
					fmt.Printf("sending message: %s\n", msg)
				}
				retained := false
				t := client.Publish(TEMP_TOPIC, QOS_0, retained, msg)
				// Handle the token in a go routine so this loop keeps sending messages regardless of delivery status
				go func() {
					_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
					if err := t.Error(); err != nil {
						fmt.Printf("ERROR PUBLISHING: %s\n", err)
					}
				}()
			case <-done:
				fmt.Println("publisher done")
				wg.Done()
				return
			}
		}
	}()

	// Wait for a signal before exiting
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	signal.Notify(stopSignal, syscall.SIGTERM)

	<-stopSignal
	fmt.Println("signal caught - exiting")

	close(done)
	wg.Wait()
	fmt.Println("shutdown complete")
}
