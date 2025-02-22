package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/siuyin/dflt"
)

type Person struct {
	Name        string
	DateOfBirth time.Time
	MassKg      float32
}

var (
	nc *nats.Conn
	ns *server.Server
)

func init() {
	nc, ns = runServer()
}

func main() {
	kv := newKeyValueStore()
	defer shutdownKeyValueStore()

	putData(kv)

	person := get(kv, "KitSiew")
	fmt.Println(person.Name)
	select {}
}

func putData(kv jetstream.KeyValue) {
	p1 := Person{Name: "SiuYin"}
	p2 := Person{Name: "KitSiew"}
	put(kv, &p1)
	put(kv, &p2)
}

func runServer() (*nats.Conn, *server.Server) {
	opts := &server.Options{JetStream: true,
		DontListen: false,
		StoreDir:   "/tmp/mystore"}
	ns, err := server.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	if dflt.EnvString("LOGGING", "") != "" {
		ns.ConfigureLogger()
	}
	go ns.Start()
	if !ns.ReadyForConnections(5 * time.Second) {
		log.Fatal("could not bring up embedded NATS server")
	}

	nc, err := nats.Connect(ns.ClientURL(), nats.InProcessServer(ns))
	if err != nil {
		log.Fatal(err)
	}

	return nc, ns
}

func newKeyValueStore() jetstream.KeyValue {
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	cfg := jetstream.KeyValueConfig{Bucket: "mykv"}
	kv, err := js.CreateKeyValue(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	return kv
}

func shutdownKeyValueStore() {
	ns.WaitForShutdown()
}

func personBytes(p *Person) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(p)
	return buf.Bytes()
}

func bytesPerson(b []byte) Person {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	var p Person
	if err := dec.Decode(&p); err != nil {
		log.Fatal(err)
	}

	return p
}

func put(kv jetstream.KeyValue, p *Person) {
	if _, err := kv.Put(context.Background(),
		p.Name, personBytes(p)); err != nil {
		log.Fatal(err)
	}
}

func get(kv jetstream.KeyValue, key string) Person {
	obj, err := kv.Get(context.Background(), key)
	if err != nil {
		log.Fatal(err)
	}

	p := bytesPerson(obj.Value())
	return p
}
