package jamyxgo

import (
    "net"
    "strconv"
    "log"
    "strings"
    "fmt"
    "sync"
)

// Session holds the connection to the server
type Session struct {
    connection net.Conn
    sendingCommand *sync.Mutex
}

// Connect to the jamyxer server
func (session* Session) Connect(ip string, port int) {
    conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
    if err != nil {
        log.Fatal(err)
    }
    session.connection = conn
    session.sendingCommand = &sync.Mutex{}
}

// Send a command to the jamyxer server
func (session* Session) SendCommand(cmd string, a ...interface{}) (reply string) {
    cmd_s := fmt.Sprintf(cmd, a...)

    // log.Println("Sending command:", []byte(cmd_s))
    log.Println("Sending command:", cmd_s)

    // Protect agaisnt multiple commands being sent at the same
    //   time due to SendCommand being called from multiple threads.
    session.sendingCommand.Lock()
    defer session.sendingCommand.Unlock()

    _, err := session.connection.Write([]byte(cmd_s))
    if err != nil {
        log.Fatal(err)
    }

    reply_b := make([]byte, 1024)
    numbytes, err := session.connection.Read(reply_b)
    if err != nil {
        log.Fatal(err)
    }

    return string(reply_b[:numbytes])
}

// Returns an array of strings representing the names of the input channels
func (session *Session) GetInputs() []string {
    reply := strings.Trim(session.SendCommand("gi\n"), "\n")
    return strings.Split(reply, "\n")
}

// Set volume for specified input channel
func (session *Session) VolumeInputSet(input string, volume float64) {
    session.SendCommand("vis \"%s\" %f\n", input, volume)
}

// Set volume for specified output channel
func (session *Session) VolumeOutputSet(output string, volume float64) {
    session.SendCommand("vos \"%s\" %f\n", output, volume)
}

// Get volume for specified input channel
func (session *Session) VolumeInputGet(input string) float64 {
    vol_s := session.SendCommand("vig \"%s\"\n", input)
    vol, _ := strconv.ParseFloat(vol_s, 64)
    return vol
}

// Get volume for specified output channel
func (session *Session) VolumeOutputGet(output string) float64 {
    vol_s := session.SendCommand("vog \"%s\"\n", output)
    vol, _ := strconv.ParseFloat(vol_s, 64)
    return vol
}

// Listen for volume for specified input channel
// This is a blocking call waiting for a change in volume and returning it
func (session *Session) VolumeInputListen(input string) float64 {
    vol_s := session.SendCommand("listni \"%s\"\n", input)
    vol, _ := strconv.ParseFloat(vol_s, 64)
    return vol
}

// Listen for volume for specified output channel
// This is a blocking call waiting for a change in volume and returning it
func (session *Session) VolumeOutputListen(output string) float64 {
    vol_s := session.SendCommand("listno \"%s\"\n", output)
    vol, _ := strconv.ParseFloat(vol_s, 64)
    return vol
}
