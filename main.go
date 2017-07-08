package jamyxgo

import (
    "net"
    "strconv"
    "log"
    "strings"
)


// Session holds the connection to the server
type Session struct {
    connection net.Conn
}

// Connect to the jamyxer server
func (session* Session) Connect(ip string, port int) {
    conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
    if err != nil {
        log.Fatal(err)
    }
    session.connection = conn
}

// Send a command to the jamyxer server
func (session* Session) SendCommand(cmd string) (reply string) {
    _, err := session.connection.Write([]byte(cmd))
    if err != nil {
        log.Fatal(err)
    }

    reply_b := make([]byte, 1024)
    _, err = session.connection.Read(reply_b)
    if err != nil {
        log.Fatal(err)
    }

    return string(reply_b)
}

// Returns an array of strings representing the names of the input channels
func (session *Session) GetInputs() []string {
    reply := strings.Trim(session.SendCommand("gi\n"), "\n")
    return strings.Split(reply, "\n")
}
