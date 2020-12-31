package impls

import (
	"alphawolf.com/alphaservice/models"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type TerminalCommand struct {
	UUID    string
	Command string
	Data    map[string]interface{}
}

func NewSshPty(uuid string, conn *websocket.Conn, userGroupIds []int, user *models.User, target *models.Target, targetOption *models.TargetOption) *SshPty {
	sshPty := &SshPty{
		uuid:         uuid,
		remoteAddr:   conn.RemoteAddr().String(),
		conn:         conn,
		user:         user,
		target:       target,
		targetOption: targetOption,
		userGroupIds: userGroupIds,
	}
	return sshPty
}

type SshPty struct {
	uuid               string
	remoteAddr         string
	userGroupIds       []int
	user               *models.User
	target             *models.Target
	targetOption       *models.TargetOption
	conn               *websocket.Conn
	sshWriteChan       chan []byte
	websocketWriteChan chan []byte
	sshClient          *ssh.Client
	stop               int32
	lastCmd            string
	lastCmdPos         int
	cmdPrompt          string
	whenPromptReplace  int32
	lastCmdMutex       sync.Mutex
}

func (this *SshPty) Stop() {
	atomic.StoreInt32(&this.stop, 1)
	this.conn.Close()
	time.Sleep(1 * time.Second)
	close(this.websocketWriteChan)
	close(this.sshWriteChan)
}

func (this *SshPty) lastCmdAppend(data []byte) {
	this.lastCmdMutex.Lock()
	defer this.lastCmdMutex.Unlock()
	cLen := len(this.lastCmd)
	if this.lastCmdPos > cLen {
		logs.Error("error last cmd pos")
		this.lastCmdPos = len(this.lastCmd)
	}
	if this.lastCmdPos == cLen {
		this.lastCmd += string(data)
		this.lastCmdPos = len(this.lastCmd)
	} else {
		this.lastCmd = this.lastCmd[0:this.lastCmdPos] + string(data) + this.lastCmd[this.lastCmdPos:]
		this.lastCmdPos += len(string(data))
	}

}

func (this *SshPty) lastCmdSet(data []byte) {
	this.lastCmdMutex.Lock()
	defer this.lastCmdMutex.Unlock()
	if data == nil {
		this.lastCmd = ""
		this.lastCmdPos = 0
	} else {
		this.lastCmd = string(data)
		this.lastCmdPos = len(this.lastCmd)
	}

}

func (this *SshPty) lastCmdGet() string {
	this.lastCmdMutex.Lock()
	defer this.lastCmdMutex.Unlock()
	return this.lastCmd
}

func (this *SshPty) lastCmdBackspace() {
	this.lastCmdMutex.Lock()
	defer this.lastCmdMutex.Unlock()
	cLen := len(this.lastCmd)
	if this.lastCmdPos > cLen {
		logs.Error("error last cmd pos")
		this.lastCmdPos = len(this.lastCmd)
	}
	if this.lastCmdPos > 0 {
		if this.lastCmdPos == cLen {
			this.lastCmd = this.lastCmd[:cLen-1]

		} else {
			this.lastCmd = this.lastCmd[:this.lastCmdPos-1] + this.lastCmd[this.lastCmdPos:]
		}
		this.lastCmdPos -= 1
	}
}

func (this *SshPty) lastCmdPosAdd() {
	this.lastCmdMutex.Lock()
	defer this.lastCmdMutex.Unlock()
	if this.lastCmdPos+1 >= len(this.lastCmd) {
		this.lastCmdPos = len(this.lastCmd)
	} else {
		this.lastCmdPos++
	}
}

func (this *SshPty) lastCmdPosSub() {
	this.lastCmdMutex.Lock()
	defer this.lastCmdMutex.Unlock()
	if this.lastCmdPos-1 < 0 {
		this.lastCmdPos = 0
	} else {
		this.lastCmdPos--
	}
}

func (this *SshPty) Run() error {
	this.sshWriteChan = make(chan []byte, 10)
	this.websocketWriteChan = make(chan []byte, 10)

	//defer close(this.websocketWriteChan)
	//defer close(this.sshWriteChan)

	addr := this.target.Address
	port, err := strconv.Atoi(this.targetOption.Port)
	if err != nil {
		return errors.New("invalid target port:" + this.targetOption.Port)
	}

	if port < 1 || port > 65535 {
		port = 22
	}

	addr += ":" + strconv.Itoa(port)

	pwd, err := PwdDecrypt(this.targetOption.Password)
	if err != nil {
		return errors.New("invalid password")
	}

	sshConn, channel, reqs, err := DialSSH(addr, this.targetOption.Username, 5, ssh.Password(pwd), ssh.KeyboardInteractive(
		func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
			if len(questions) == 0 {
				return []string{}, nil
			}
			return []string{pwd}, nil
		}))

	if err != nil {
		return err
	}

	this.sshClient = ssh.NewClient(sshConn, channel, reqs)
	defer this.sshClient.Close()

	// Set up new Session between server and host terminal via ssh
	session, err := this.sshClient.NewSession()
	if err != nil {
		return errors.New("create ssh session failed, " + err.Error())
	}
	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 30, modes); err != nil {
		return errors.New("request for pseudo terminal failed, " + err.Error())
	}

	//set io.Reader and io.Writer from terminal session
	sshReader, err := session.StdoutPipe()
	if err != nil {
		return errors.New("set ssh session stdout pipe failed, " + err.Error())
	}

	sshWriter, err := session.StdinPipe()
	if err != nil {
		return errors.New("set ssh session stdin pipe failed, " + err.Error())
	}

	defer sshWriter.Close()

	//websocket writer
	go func() {
		//defer logs.Trace("websocket writer exit ...")
		for {
			stop := atomic.LoadInt32(&this.stop)
			if stop == 1 {
				return
			}

			select {
			case payload, ok := <-this.websocketWriteChan:
				if !ok {
					return
				}
				err := this.conn.WriteMessage(websocket.BinaryMessage, payload)
				if err != nil {
					logs.Error(err)
					return
				}
			}
		}
	}()

	//ssh writer
	go func() {
		//defer logs.Trace("ssh writer exit ...")
		for {
			stop := atomic.LoadInt32(&this.stop)
			if stop == 1 {
				return
			}
			select {
			case payload, ok := <-this.sshWriteChan:
				if !ok {
					return
				}
				_, err := sshWriter.Write(payload)
				if err != nil {
					logs.Error(err)
					this.websocketWriteChan <- []byte(err.Error())
					return
				}
			}
		}
	}()

	//ssh reader
	go func() {
		//defer logs.Trace("ssh reader exit ...")
		for {
			stop := atomic.LoadInt32(&this.stop)
			if stop == 1 {
				return
			}

			buf := make([]byte, 4096)
			n, err := sshReader.Read(buf)
			if err != nil {
				logs.Error(err)
				return
			}

			OperateAuditSrv.Add(this.uuid, 3, this.user, this.target, this.targetOption, this.remoteAddr, string(buf[:n]))

			replaceLastCmd := atomic.LoadInt32(&this.whenPromptReplace)
			if replaceLastCmd == 1 {
				this.lastCmdSet(buf[:n])
				atomic.StoreInt32(&this.whenPromptReplace, 0)
			}

			//logs.Error("ssh cmd output:",string(buf[:n]))

			output, err := OperateAclSrv.CmdOutputCheck(string(buf[:n]), this.targetOption.Id, this.target.Id, this.target.GroupId, this.user.Id, this.userGroupIds)
			if err != nil {
				this.sshWriteChan <- []byte{0x0d}
				this.websocketWriteChan <- []byte("\r\ncmd output check failed, " + err.Error())
				this.sshWriteChan <- []byte{0x0d}
				OperateAuditSrv.Add(this.uuid, 2, this.user, this.target, this.targetOption, this.remoteAddr, "cmd output check failed, "+err.Error())
			} else {
				this.websocketWriteChan <- []byte(output)
			}
		}
	}()

	//websocket reader
	go func() {
		//defer logs.Trace("websocket reader exit ...")
		for {
			//logs.Trace("websocket read ...")
			stop := atomic.LoadInt32(&this.stop)
			if stop == 1 {
				return
			}

			// set up io.Reader of websocket
			_, reader, err := this.conn.NextReader()
			if err != nil {
				logs.Error("websocket read error, ", err.Error())
				return
			}
			// read first byte to determine whether to pass data or resize terminal
			dataTypeBuf := make([]byte, 1)
			_, err = reader.Read(dataTypeBuf)
			if err != nil {
				logs.Error(err)
				return
			}

			switch dataTypeBuf[0] {
			// when pass data
			case 0:
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)
				if err != nil {
					logs.Error(err)
					return
				}
				//logs.Trace("--cmd:",string(buf[:n]),"\n",hex.Dump(buf[:n]))

				OperateAuditSrv.Add(this.uuid, 2, this.user, this.target, this.targetOption, this.remoteAddr, string(buf[:n]))

				switch string(buf[:n]) {
				case string(0x0d):
					cmd := this.lastCmdGet()
					//logs.Trace("++cmd :",cmd)
					this.lastCmdSet(nil)
					cmdLen := len(cmd)
					cmd = strings.TrimSpace(cmd)

					ok, err := OperateAclSrv.CmdInputCheck(cmd, this.targetOption.Id, this.target.Id, this.target.GroupId, this.user.Id, this.userGroupIds)
					if err != nil {
						tmp := make([]byte, cmdLen)
						space := ""
						for i := 0; i < cmdLen; i++ {
							tmp[i] = 0x7f
							space += " "
						}
						this.sshWriteChan <- tmp

						this.sshWriteChan <- []byte{0x0d}
						this.websocketWriteChan <- []byte("\r\ncommand check failed, " + err.Error() + space)
						this.sshWriteChan <- []byte{0x0d}
						OperateAuditSrv.Add(this.uuid, 2, this.user, this.target, this.targetOption, this.remoteAddr, "command check failed, "+err.Error())
						continue
					} else if !ok {
						tmp := make([]byte, cmdLen)
						space := ""
						for i := 0; i < cmdLen; i++ {
							tmp[i] = 0x7f
							space += " "
						}
						this.sshWriteChan <- tmp
						//this.websocketWriteChan <- []byte{0x7f}
						this.sshWriteChan <- []byte{0x0d}
						this.websocketWriteChan <- []byte("\r\nnot permit command" + space)
						this.sshWriteChan <- []byte{0x0d}
						OperateAuditSrv.Add(this.uuid, 2, this.user, this.target, this.targetOption, this.remoteAddr, "not permit command")
						continue
					}
				case string(0x7f):
					this.lastCmdBackspace()
				case string([]byte{0x1b, 0x5b, 0x41}), string([]byte{0x1b, 0x5b, 0x42}):
					atomic.StoreInt32(&this.whenPromptReplace, 1)
				case string([]byte{0x1b, 0x5b, 0x44}):
					this.lastCmdPosSub()
				case string([]byte{0x1b, 0x5b, 0x43}):
					this.lastCmdPosAdd()
				default:
					this.lastCmdAppend(buf[:n])
				}

				this.sshWriteChan <- buf[:n]

			// when resize terminal
			case 1:
				decoder := json.NewDecoder(reader)
				terminalCmd := TerminalCommand{}
				err := decoder.Decode(&terminalCmd)
				if err != nil {
					errMsg := "invalid function command, " + err.Error()
					this.websocketWriteChan <- []byte(errMsg)
					continue
				}

				//cmdString,_ := util.JsonIter.Marshal(terminalCmd)
				//logs.Trace("cmdString:",string(cmdString))

				switch terminalCmd.Command {
				case "resize":
					rowsObj, ok := terminalCmd.Data["rows"]
					if !ok {
						errMsg := "invalid rows parameter"
						logs.Error(errMsg)
						this.websocketWriteChan <- []byte(errMsg)
						continue
					}

					colsObj, ok := terminalCmd.Data["cols"]
					if !ok {
						errMsg := "invalid cols parameter"
						logs.Error(errMsg)
						this.websocketWriteChan <- []byte(errMsg)
						continue
					}

					//logs.Trace(reflect.TypeOf(colsObj).String())

					rows, ok := rowsObj.(float64)
					if !ok {
						errMsg := "invalid rows int parameter"
						logs.Error(errMsg)
						this.websocketWriteChan <- []byte(errMsg)
						continue
					}

					cols, ok := colsObj.(float64)
					if !ok {
						errMsg := "invalid cols int parameter"
						logs.Error(errMsg)
						this.websocketWriteChan <- []byte(errMsg)
						continue
					}

					err = session.WindowChange(int(rows), int(cols))
					if err != nil {
						errMsg := "resize window change failed, " + err.Error()
						logs.Error(errMsg)
						this.websocketWriteChan <- []byte(errMsg)
						continue
					}

				default:
					errMsg := "unknown function command: " + terminalCmd.Command
					logs.Error(errMsg)
					this.websocketWriteChan <- []byte(errMsg)
					continue
				}

			// unexpected data
			default:
				errMsg := "unexpected data type"
				logs.Error(errMsg)
				this.websocketWriteChan <- []byte(errMsg)
			}
		}
	}()

	// Start remote shell
	if err := session.Shell(); err != nil {
		return errors.New("start shell failed, " + err.Error())
	}

	OperateAuditSrv.Add(this.uuid, 1, this.user, this.target, this.targetOption, this.remoteAddr, "ssh connected")

	return session.Wait()
}

var SshTerminalSrv SshTerminal

type SshTerminal struct {
	TerminalMap sync.Map
}
