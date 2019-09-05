package main

import (
	"log"
    "fmt"
    "bufio"
    "net/http"
    "os"
    "errors"
    "sync"
    "strings"
    "math/big"
    "time"
    "encoding/json"
)

// Use mutex when read/write file
// to avoid race condition
type ReadWriteMutex struct {
	mux sync.Mutex
}

func main() {
    // Just use one CPU core here, (required from the requirements)
    runtime.GOMAXPROCS(1)
    fmt.Println("Running server at port 8080")
    mux := http.NewServeMux()
    rw1 := ReadWriteMutex{}
    rw2 := ReadWriteMutex{}
    mux.HandleFunc("/current/", Current)
    mux.HandleFunc("/next/", rw1.Next)
    mux.HandleFunc("/prev/", rw2.Prev)
    log.Fatal(http.ListenAndServe(":8080", RequestLogger(mux)))
}

// Return the current fibonacci num from data.log
func Current(w http.ResponseWriter, r *http.Request) {
    current, _, err := ReadPrevAndCurrent("data.log")
    // something wrong when reading file
    if err != nil {
        JsonResponse(w, "", 500, err.Error())
    } else {
        str_current := current.String()
        // /current/ Should always succeed
        JsonResponse(w, str_current, 200, "00000")
    }
}

// Return the next fibonacci num from data.log
func (rw1 *ReadWriteMutex) Next(w http.ResponseWriter, r *http.Request) {
    // Use mutex when read/write file
    rw1.mux.Lock()
    defer rw1.mux.Unlock()
    current, next, err := ReadPrevAndCurrent("data.log")
    // Read data.log error
    if err != nil {
        JsonResponse(w, "", 500, err.Error())
    } else {
        str_next := next.String()
        // fibonacci calculation
        current.Add(current, next)
        current, next = next, current
        err := UpdatePrevAndCurrent("data.log", current, next)
        // Write data.log error
        if err != nil {
            JsonResponse(w, "", 500, err.Error())
        } else {
            JsonResponse(w, str_next, 200, "00000")
        }
   }
}

// Return the prev fibonacci num from data.log
func (rw2 *ReadWriteMutex) Prev(w http.ResponseWriter, r *http.Request) {
    rw2.mux.Lock()
    defer rw2.mux.Unlock()
    current, next, err := ReadPrevAndCurrent("data.log")
    if err != nil {
        JsonResponse(w, "", 500, err.Error())
    } else {
        // Return 400 when we want to get the prev value of 0
        zero := big.NewInt(0)
        if (current.Cmp(zero) == 0) {
            JsonResponse(w, "Nil", 400, "10000")
        } else {
            next.Sub(next, current)
            str_prev := next.String()
            current, next = next, current
            err := UpdatePrevAndCurrent("data.log", current, next)
            if err != nil {
                JsonResponse(w, "", 500, err.Error())
            } else {
                JsonResponse(w, str_prev, 200, "00000")
            }
        }
    }
}

// Construct response from input data, the expected response format looks like
// {
//     "code":"00000", 
//     "data":"1",
//     "msg":"Succeed"
//  }
// We use "code" and "msg" to tell frontend what's going wrong in the server
// The "data" is the fibonacci val we need 
func JsonResponse(w http.ResponseWriter, data string, statusCode int, code string) {
    codeToMessage := map[string]string{
        "00000": "Succeed",
        "10000": "There is no previous value for 0",
        "10001": "data.log not found",
        "10002": "Can't read the first line from data.log",
        "10003": "Can't read the second line from data.log",
        "10004": "Can't update data.log",
    }
    if (statusCode == 200) {
        msg := codeToMessage[code]
        FormatJson(w, data, statusCode, code, msg)
    } else {
        msg := codeToMessage[code]
        FormatJson(w, data, statusCode, code, msg)
    }
}


// Construct json format
func FormatJson(w http.ResponseWriter, data string, statusCode int, code string, msg string) {
    jsonMap := make(map[string]string)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    jsonMap["code"] = code
    jsonMap["msg"] = msg
    jsonMap["data"] = data
    json.NewEncoder(w).Encode(jsonMap)
}

// Read the first and second line from "name"
// then convert them to big.int
func ReadPrevAndCurrent(name string) (*big.Int, *big.Int, error){
    f, err := os.Open(name) // os.OpenFile has more options if you need them
    if err != nil {
        err := errors.New("10001")
        zero := big.NewInt(0)
        return zero, zero, err
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // Get first line
    current, err := rd.ReadString('\n')
    current = strings.TrimSuffix(current, "\n")
    if current == "" || err != nil {
        err := errors.New("10002")
        zero := big.NewInt(0)
        return zero, zero, err
    }
    int_current := StringToInt(current)

    // Get second line
    next, err := rd.ReadString('\n')
    next = strings.TrimSuffix(next, "\n")
    if next == "" || err != nil {
        err := errors.New("10003")
        zero := big.NewInt(0)
        return zero, zero, err
    }
    int_next := StringToInt(next)
    return int_current, int_next, nil
}

// Rewrite the first and second line from "name"
func UpdatePrevAndCurrent(name string, current *big.Int, next *big.Int) (error) {
    f, err := os.OpenFile(name, os.O_WRONLY, 0644)
    if err != nil {
        err := errors.New("10004")
        return err
    }
    defer f.Close()
    str_current := current.String()
    str_next := next.String()
    _, err = f.WriteString(str_current + "\n")
    if err != nil {
        err := errors.New("10004")
        return err
    }
    _, err = f.WriteString(str_next + "\n")
    if err != nil {
        err := errors.New("10004")
        return err
    }
    return nil
}

// Logging for debug 
func RequestLogger(targetMux http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        targetMux.ServeHTTP(w, r)
        requesterIP := r.RemoteAddr
        log.Printf(
            "%s\t\t%s\t\t%s\t\t%v",
            r.Method,
            r.RequestURI,
            requesterIP,
            time.Since(start),
        )
    })
}

// Convert string to big.int
func StringToInt(str string) (*big.Int) {
    n := new(big.Int)
    n, err := n.SetString(str, 10)
    if !err {
        log.Fatal(err)
    }
    return n
}
