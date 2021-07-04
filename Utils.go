package main

import "fmt"
import "strconv"
import "time"
import "math/rand"
import (
    "crypto/sha512"
    "encoding/base64"
)




func check_error(err error) {
    if err != nil {
      fmt.Println("Error: " , err)
    }
}


func random(maxNumber int) int {
    s1 := rand.NewSource(time.Now().UnixNano());
    r1 := rand.New(s1);
    return r1.Intn(maxNumber);
}


func now() int {
    return int(time.Now().Unix());
}


func generate_UUID() string {
    return hash(strconv.Itoa(random(4294967296)));
}


func hash(input string) string {
    hasher := sha512.New();

    hasher.Write([]byte(input));
    return base64.URLEncoding.EncodeToString(hasher.Sum(nil));
}