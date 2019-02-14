package controllers

import (
	"net/http"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"io"
	"golang.org/x/crypto/ssh"
	"time"
	"os"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/ziutek/telnet"
	"fmt"
)

type ShellWsController struct {
	beego.Controller
}

type wsWrapper struct {
	*websocket.Conn
}

func (wsw *wsWrapper) Write(p []byte) (n int, err error) {
	writer, err := wsw.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	return writer.Write(p)
}

func (wsw *wsWrapper) Read(p []byte) (n int, err error) {
	for {
		msgType, reader, err := wsw.Conn.NextReader()
		if err != nil {
			return 0, err
		}
		if msgType != websocket.TextMessage {
			continue
		}
		return reader.Read(p)
	}
}

func sshHandle(rw io.ReadWriter,ip,port,user,passwd string,errhandle func(string)){
	sshConfig :=  &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(passwd)},
		Timeout:6*time.Second,
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	if port == "" {
		ip = ip + ":22"
	} else {
		ip = ip + ":" + port
	}
	client, err := ssh.Dial("tcp",ip, sshConfig)
	if err != nil {
		beego.Debug(err.Error())
		errhandle(err.Error())
		return
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		beego.Debug(err.Error())
		errhandle(err.Error())
		return
	}
	defer session.Close()
	fd := int(os.Stdin.Fd())
	session.Stdout = rw
	session.Stderr = rw
	session.Stdin = rw
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	termWidth, termHeight, err := terminal.GetSize(fd)
	err = session.RequestPty("xterm",termHeight,termWidth, modes)
	if err != nil {
		errhandle(err.Error())
	}
	err = session.Shell()
	if err != nil {
		errhandle(err.Error())
	}
	err = session.Wait()
	if err != nil {
		errhandle(err.Error())
	}
	return
}

func telnetHandle (rw io.ReadWriter,ip,port string,errhandle func(string)){
	if port == "" {
		ip = ip + ":23"
	} else {
		ip = ip + ":"+port
	}
	con,err := telnet.Dial("tcp",ip)
	if err != nil {
		errhandle(err.Error())
		return
	}
	defer con.Close()
	buf := make([]byte,16*1024)
	for ;; {
		n,err := con.Read(buf)
		if err != nil {
			errhandle(err.Error())
			break
		}
		_,err = rw.Write(buf[:n])
		if err != nil {
			errhandle(err.Error())
			break
		}
		n,err = rw.Read(buf)
		if err != nil {
			errhandle(err.Error())
			break
		}
		if buf[0] == 13 {
			data := []byte{telnet.CR,telnet.LF}
			_,err = con.Write(data)
		} else {
			_,err = con.Write(buf[:n])
		}
		if err != nil {
			errhandle(err.Error())
			break
		}
	}
	return
}

func websocketHandle(con *websocket.Conn,shellinfo shellInfoStruct){
	rw := io.ReadWriter(&wsWrapper{con})
	webprintln := func(data string){
		rw.Write([]byte(data+"\r\n"))
	}
	con.SetCloseHandler(func(code int, text string) error{
		con.Close()
		return nil
	})
	switch (shellinfo.Proto) {
	case "ssh":
		sshHandle(rw,shellinfo.IpAddr,shellinfo.Port,shellinfo.User,shellinfo.Passwd,webprintln)
	case "telnet":
		telnetHandle(rw,shellinfo.IpAddr,shellinfo.Port,webprintln)
	default:
		webprintln("Not Support Protocol '"+shellinfo.Proto +"'")
	}
	return
}

func (c *ShellWsController)Get(){
	c.EnableRender = false
	rethandle := func(err error,status int){
		reply := "ok"
		if err != nil {
			reply = err.Error()
		}
		c.Data["json"] = reply
		c.Ctx.Output.Status = status
		c.ServeJSON()
		return
	}
	v := c.GetSession("keyinfo")
	shellinfo,ok := v.(shellInfoStruct)
	if ok == false {
		fmt.Println("Can't Found session info")
		rethandle(errors.New("Can't Found session info"),http.StatusOK)
		return
	}
	info := shellinfo
	c.DelSession("keyinfo")
	con,err := upgrader.Upgrade(c.Ctx.ResponseWriter,c.Ctx.Request,nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		beego.Error("Not a websocket connection")
		rethandle(errors.New("Not a websocket handshake"),http.StatusOK)
		return
	} else if err != nil {
		rethandle(err,http.StatusOK)
		return
	}
	go websocketHandle(con,info)
	return
}

func (c *ShellWsController) Delete() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}

func (c *ShellWsController) Put() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}

func (c *ShellWsController) Post() {
	http.Error(c.Ctx.ResponseWriter, "Method Denied", 405)
}
