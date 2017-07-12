package jamyxgo

import (
    "net"
    "strconv"
    "log"
    "strings"
    "fmt"
    "sync"
)

// Session holds the connection to the server.
type Session struct {
    connection net.Conn
    sendingCommand *sync.Mutex
}

// Connect to the jamyxer server.
func (session* Session) Connect(ip string, port int) {
    conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
    if err != nil {
        log.Fatal(err)
    }
    session.connection = conn
    session.sendingCommand = &sync.Mutex{}
}

// Send a command to the jamyxer server.
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

func getTargetType(isinput bool) string {
    if isinput { return "i" }
    return "o"
}

// ==== Get Channels ====

// Returns an array of strings representing the names of the input/output channels.
func (session *Session) GetChannels(isinput bool) []string {
    reply := strings.Trim(session.SendCommand("g%s\n", getTargetType(isinput)), "\n")
    return strings.Split(reply, "\n")
}
// Returns an array of strings representing the names of the input channels.
func (session *Session) GetInputs() []string { return session.GetChannels(true) }
// Returns an array of strings representing the names of the output channels.
func (session *Session) GetOutputs() []string { return session.GetChannels(false) }

// ==== Set Volume ====

// Set volume for specified input/output channel.
func (session *Session) VolumeSet(isinput bool, channel string, volume float64) {
    session.SendCommand("v%ss \"%s\" %f\n", getTargetType(isinput), channel, volume)
}
// Set volume for specified input channel.
func (session *Session) VolumeInputSet(input string, volume float64) { session.VolumeSet(true, input, volume) }
// Set volume for specified output channel.
func (session *Session) VolumeOutputSet(output string, volume float64) { session.VolumeSet(false, output, volume) }

// ==== Get Volume ====

// Get volume for specified input/output channel.
func (session *Session) VolumeGet(isinput bool, channel string) float64 {
    vol_s := session.SendCommand("v%sg \"%s\"\n", getTargetType(isinput), channel)
    vol, _ := strconv.ParseFloat(vol_s, 64)
    return vol
}
// Get volume for specified input channel.
func (session *Session) VolumeInputGet(input string) float64 { return session.VolumeGet(true, input) }
// Get volume for specified output channel.
func (session *Session) VolumeOutputGet(output string) float64 { return session.VolumeGet(false, output) }

// ==== Listeners ====
// Listen for volume change for specified channel.
// This is a blocking call waiting for a change in volume and returning it.
func (session *Session) VolumeListen(isinput bool, channel string) float64 {
    vol_s := session.SendCommand("v%sln \"%s\"\n", getTargetType(isinput), channel)
    vol, _ := strconv.ParseFloat(vol_s, 64)
    return vol
}
// Listen for volume for specified input channel.
// This is a blocking call waiting for a change in volume and returning it.
func (session *Session) VolumeInputListen(input string) float64 { return session.VolumeListen(true, input) }
// Listen for volume for specified output channel.
// This is a blocking call waiting for a change in volume and returning it.
func (session *Session) VolumeOutputListen(output string) float64 { return session.VolumeListen(false, output) }
