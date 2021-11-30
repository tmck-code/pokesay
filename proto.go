package main

import (
    "fmt"
    "log"
	"os"

    "google.golang.org/protobuf/proto"
)

func main() {
  elliot := &Person{
      Name: "Elliot",
      Age:  24,
  }

  data, err := proto.Marshal(elliot)
  if err != nil {
      log.Fatal("marshaling error: ", err)
  }

  // printing out our raw protobuf object
  fmt.Println(data)
  err = os.WriteFile("data.txt", data, 0644)

  newPerson := &Person{}
  data2, _ := os.ReadFile("data.txt")
  fmt.Println("read file", data2)
  err = proto.Unmarshal(data2, newPerson)
  if err != nil {
      log.Fatal("unmarshaling error: ", err)
  }

  fmt.Println(newPerson)
}

