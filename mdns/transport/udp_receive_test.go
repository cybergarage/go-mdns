package transport

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

func TestUDPMessageRetainsReceivedBytes(t *testing.T) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	socket := NewUDPSocket(dns.TransportUDP)
	socket.SetConn(conn)
	defer socket.Close()
	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	sender, err := net.DialUDP("udp4", nil, conn.LocalAddr().(*net.UDPAddr))
	if err != nil {
		t.Fatal(err)
	}
	defer sender.Close()

	receive := func(id uint) (dns.Message, []byte) {
		t.Helper()
		wire := dns.NewRequestMessage(dns.WithMessageID(id)).Bytes()
		if _, err := sender.Write(wire); err != nil {
			t.Fatal(err)
		}
		msg, err := socket.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}
		return msg, wire
	}
	first, wire := receive(1)
	receive(2)
	if !bytes.Equal(first.Bytes(), wire) {
		t.Fatal("first message bytes changed after the second receive")
	}
}

func TestUDPReceiveDiscardsInvalidDatagram(t *testing.T) {
	multicast := NewMulticastServer()
	unicast := NewUnicastServer()
	for _, tc := range []struct {
		name      string
		socket    *UDPSocket
		processor func(dns.MessageProcessor)
		run       func(chan any)
	}{
		{"multicast", multicast.UDPSocket, multicast.SetMessageProcessor, func(done chan any) {
			handleMulticastConnection(multicast, done)
		}},
		{"unicast", unicast.UDPSocket.UDPSocket, unicast.SetMessageProcessor, func(done chan any) {
			handleUnicastUDPConnection(unicast, done)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			listener, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
			if err != nil {
				t.Fatal(err)
			}
			tc.socket.SetConn(listener)
			sender, err := net.DialUDP("udp4", nil, listener.LocalAddr().(*net.UDPAddr))
			if err != nil {
				listener.Close()
				t.Fatal(err)
			}
			defer sender.Close()

			received := make(chan uint, 2)
			tc.processor(func(msg dns.Message) (dns.Message, error) {
				received <- msg.ID()
				return nil, nil
			})
			stop := make(chan any)
			exited := make(chan struct{})
			go func() {
				defer close(exited)
				tc.run(stop)
			}()
			defer func() {
				close(stop)
				tc.socket.Close()
				<-exited
			}()

			checkMessage := func(id uint) {
				t.Helper()
				msg := dns.NewRequestMessage(dns.WithMessageID(id))
				if _, err := sender.Write(msg.Bytes()); err != nil {
					t.Fatal(err)
				}
				select {
				case got := <-received:
					if got != id {
						t.Fatalf("received message ID %d, want %d", got, id)
					}
				case <-exited:
					t.Fatal("receive loop stopped")
				case <-time.After(time.Second):
					t.Fatal("receiver did not process the valid message")
				}
			}
			checkMessage(1)
			if _, err := sender.Write(nil); err != nil {
				t.Fatal(err)
			}
			checkMessage(2)
		})
	}
}
