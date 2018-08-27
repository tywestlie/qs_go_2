package main

import (
  
)

func main() {
    a := App{}
    a.Initialize("dbname=qs_go_dev sslmode=disable")

    a.Run(":3000")
}
