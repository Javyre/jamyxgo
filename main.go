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
    if numbytes == 0 {
        log.Fatal("Numbytes is 0!")
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

// ==== Set Connected ====

// Connect input with output
func (session *Session) ConnectIO(input, output string) {
    session.SendCommand(`c "%s" "%s"`, input, output)
}
// Toggle connection between input and output
func (session *Session) ToggleConnectionIO(input, output string) {
    session.SendCommand(`tc "%s" "%s"`, input, output)
}
// Disconnect input and output
func (session *Session) DisconnectIO(input, output string) {
    session.SendCommand(`dc "%s" "%s"`, input, output)
}

// ==== Get Connected ====

// Return true if output & input are connected
func (session *Session) GetConnectedIO(input, output string) bool {
    ret := session.SendCommand(`gc "%s" "%s"`, output, input)
    return ret == "1"
}

// ==== Set Monitor ====
func (session *Session) SetMonitor(isinput bool, channel string) {
    t := 'o'; if isinput { t = 'i' }
    session.SendCommand(`mn%c "%s"`, t, channel)
}

// ==== Get Monitor ====

// Get name of the channel being monitored
func (session *Session) GetMonitorChannel() string {
    return session.SendCommand("gmn");
}
// Get type of the channel being monitored
func (session *Session) MonitorIsInput() bool {
    return session.SendCommand("gmni") == "1"
}
// Get the channel being monitored name and type
func (session *Session) GetMonitor() (isinput bool, channel string) {
    isinput = session.MonitorIsInput()
    channel = session.GetMonitorChannel()
    return isinput, channel
}

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
