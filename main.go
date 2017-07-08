package jamyxgo

import (
    "net"
    "strconv"
    "log"
    "strings"
)


type Session struct {
    connection net.Conn
}

func (session* Session) Connect(ip string, port int) {
    conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
    if err != nil {
        log.Fatal(err)
    }
    session.connection = conn
}

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

func (session *Session) GetInputs() []string {
    reply := strings.Trim(session.SendCommand("gi\n"), "\n")
    return strings.Split(reply, "\n")
}
