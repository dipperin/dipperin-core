// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.


package util

import (
	"flag"
	"fmt"
	"github.com/json-iterator/go"
	"math/rand"
	"net"
	"os"
	"os/user"
	"reflect"
	"strings"
	"time"
)

// GetUniqueId get the unique id associated with the timestamp
func GetUniqueId() string {
	curNano := time.Now().UnixNano()
	r := rand.New(rand.NewSource(curNano))
	return fmt.Sprintf("%d%06v", curNano, r.Int31n(1000000))
}

// ParseJson parsing json strings
func ParseJson(data string, result interface{}) error {
	return ParseJsonFromBytes([]byte(data), result)
}

// StringifyJson json to string
func StringifyJson(data interface{}) string {
	return string(StringifyJsonToBytes(data))
}

// ParseJsonFromBytes parsing json bytes
func ParseJsonFromBytes(data []byte, result interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, result)
}

// json bytes to string
func StringifyJsonToBytes(data interface{}) []byte {
	b, _ := StringifyJsonToBytesWithErr(data)
	return b
}

func StringifyJsonToBytesWithErr(data interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(&data)
	return b, err
}

// LoadJSON reads the given file and unmarshals its content.
//func LoadJSON(file string, val interface{}) error {
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//
//	content, err := ioutil.ReadFile(file)
//	if err != nil {
//		return err
//	}
//	if err := json.Unmarshal(content, val); err != nil {
//		//if syntaxerr, ok := err.(*json.SyntaxError); ok {
//		//	line := findLine(content, syntaxerr.Offset)
//		//	return fmt.Errorf("JSON syntax error at %v:%v: %v", file, line, err)
//		//}
//		return fmt.Errorf("JSON unmarshal error in %v: %v", file, err)
//	}
//	return nil
//}


// findLine returns the line number for the given offset into data.
//func findLine(data []byte, offset int64) (line int) {
//	line = 1
//	for i, r := range string(data) {
//		if int64(i) >= offset {
//			return
//		}
//		if r == '\n' {
//			line++
//		}
//	}
//	return
//}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

//func AbsolutePath(Datadir string, filename string) string {
//	if filepath.IsAbs(filename) {
//		return filename
//	}
//	return filepath.Join(Datadir, filename)
//}

// get home dir
func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

// kill some app
//func KillProcess(appName string) {
//	exec.Command("/bin/bash", "-c", fmt.Sprintf("ps aux|grep %v|awk '{print $2}'|xargs kill ", appName)).Output()
//	//cslog.Info().Err(err).Str("result", string(rb)).Str("app Name", appName).Msg("execute kill program command ")
//}
//
//// kill a remote application
//func KillRemoteProcess(sshConnStr, appName string) {
//	ExecCmdBySSH(sshConnStr, fmt.Sprintf("ps aux|grep %v|awk '{print $2}'|xargs kill ", appName))
//}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// GetCurPcIp get the local ip according to the matching conditions
func GetCurPcIp(matcher string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		if strings.Contains(a.String(), matcher) {
			return strings.Split(a.String(), "/")[0]
		}
	}
	return ""
}

// determine if it is a remote host
//func IsRemoteHost(host string, ipMatcher string) bool {
//	curIp := GetCurPcIp(ipMatcher)
//	// not any local ip configuration then it is remote
//	if host != "" && host != "localhost" && host != "127.0.0.1" && host != curIp {
//		return true
//	}
//	return false
//}
//
//// be sure not to enter a username
//func CpFileToTargetHost(filePath string, loginUser string, targetHost string, targetPath string) {
//	scpCmdStr := fmt.Sprintf("scp %v %v@%v:%v", filePath, loginUser, targetHost, targetPath)
//	log.Debug("execute copy file to remote", "cmd", scpCmdStr)
//	rb, err := exec.Command("/bin/bash", "-c", scpCmdStr).Output()
//	if err != nil {
//		log.Warn("execute copy file to remote error ", "err", err, "result", string(rb))
//	}
//}
//
//// execute commands with ssh
//func ExecCmdBySSH(sshConnStr string, cmdStr string) {
//	//rb, err := exec.Command("ssh", sshConnStr, `"` + cmdStr + `"`).Output()
//	tmpShFile := filepath.Join("/tmp", GetUniqueId() + "_tmp_cs.sh")
//	ioutil.WriteFile(tmpShFile, []byte(cmdStr), 0755)
//
//	cmdExecStr := fmt.Sprintf(`ssh %v bash -s < %v`, sshConnStr, tmpShFile)
//	log.Debug("run the remote ssh command", "cmd exec str", cmdExecStr, "cmdStr", cmdStr)
//	rb, err := exec.Command("/bin/bash", "-c", cmdExecStr).Output()
//	if err != nil {
//		log.Warn("run the remote ssh command error", "err", err, "result", string(rb))
//	}
//	os.RemoveAll(tmpShFile)
//}

// determine if it is a test environment
func IsTestEnv() bool {
	return flag.Lookup("test.v") != nil
}

// determine if you are not registering prometheus
//func IgnorePrometheus() bool {
//	// m_port
//	return flag.Lookup("m_port") == nil || IsTestEnv()
//}
//
//// check if the slice contains str
//func StrContainsInSlice(ss []string, str string) bool {
//	for _, s := range ss {
//		if s == str {
//			return true
//		}
//	}
//	return false
//}

func StopChanClosed(stop chan struct{}) bool {
	if stop == nil {
		return true
	}
	select {
	case _, ok := <- stop:
		return !ok
	default:
		return false
	}
}

func ExecuteFuncWithTimeout(f func(), t time.Duration) {
	finish := make(chan struct{})
	go func() {
		f()
		close(finish)
	}()

	timer := time.NewTimer(t)
	// The caller will block here until f is executed or the timer triggers timeout.
	select {
	case <- finish:
		timer.Stop()
		return
	case <- timer.C:
		panic("exec func timeout")
	}
}

// If the finishFunc is not called outside, the timeout will be triggered.
func SetTimeout(timeoutFunc func(), dur time.Duration) (finishFunc func()) {
	finishChan := make(chan struct{})

	timer := time.NewTimer(dur)
	// This goroutine is a timer. If you don't do anything outside, it will also trigger at a fixed time.
	go func() {
		select {
		case <- finishChan:
			timer.Stop()
			return
		case <- timer.C:
			timeoutFunc()
		}
	}()

	return func() {
		close(finishChan)
	}
}

// like: []interface{} -> []*Block, or []*Block -> []interface{}
func InterfaceSliceCopy(to, from interface{}) {
	toV := reflect.ValueOf(to)
	fromV := reflect.ValueOf(from)
	fLen := fromV.Len()
	for i := 0; i < fLen; i++ {
		// support bothway copy
		toV.Index(i).Set(reflect.ValueOf(fromV.Index(i).Interface()))
	}
}

// determine if the interface is nil
func InterfaceIsNil(i interface{}) bool {
	v := reflect.ValueOf(i)
	if v.IsValid() && !v.IsNil() {
		return false
	}
	return true
}