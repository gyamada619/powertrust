package command

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"github.com/kyokomi/emoji"
)

type ServiceCommand struct {
	Meta
}

func (c *ServiceCommand) Run(args []string) int {

	http.HandleFunc("/", downloadHandler)
	http.HandleFunc("/sign", signHandler)
	http.HandleFunc("/delete", deleteHandler)

	started := emoji.Sprint("Now serving your scrupulous script-signing API :white_check_mark:")
	fmt.Println(started)

	stop := emoji.Sprint("To kill this process, press Crtl + C :skull:")
	fmt.Println(stop)

	http.ListenAndServe(":7974", nil)
	return 0
}

func deleteHandler(writer http.ResponseWriter, request *http.Request) {

	servePath := os.Getenv("POWERTRUST_SERVE")
	Filename := request.URL.Query().Get("file")
	err := os.Remove(servePath + Filename)
	if err != nil {
		//File not found, send 404
		http.Error(writer, "File not found.", 404)
		return
	}

}

// HandleClient handles request for file from client.
func downloadHandler(writer http.ResponseWriter, request *http.Request) {

	servePath := os.Getenv("POWERTRUST_SERVE")
	//First of check if Get is set in the URL
	Filename := request.URL.Query().Get("file")
	if Filename == "" {
		//Get not set, send a 400 bad request
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}
	fmt.Println("A client requested a download of: " + Filename)

	//Check if file exists and open
	Openfile, err := os.Open(servePath + Filename)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(writer, "File not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	writer.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	writer.Header().Set("Content-Type", FileContentType)
	writer.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(writer, Openfile) //'Copy' the file to the client
	return
}

func signHandler(w http.ResponseWriter, r *http.Request) {

	servePath := os.Getenv("POWERTRUST_SERVE")
	// The FormFile function takes in the POST with the file.
	file, header, err := r.FormFile("fileUploadName")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	// Create file on disk.
	out, err := os.Create(servePath + header.Filename)
	if err != nil {
		fmt.Fprintf(w, "Unable to write file to disk. Ensure you have write access.")
		return
	}

	// Write the content from the POST request to the new file.
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	out.Close()

	// Local PowerShell "backend" for execution.
	back := &backend.Local{}

	// Start a local PowerShell process.
	shell, err := ps.New(back)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	defer shell.Exit()

	// Sign the PowerShell script and panic if there is an error.
	stdout, stderr, err := shell.Execute("Set-AuthenticodeSignature -FilePath " + servePath + header.Filename + " -Certificate (Get-ChildItem -Path 'Cert:\\LocalMachine\\My' -CodeSigningCert)")
	if err != nil {
		fmt.Fprintln(w, err)
		fmt.Print(stderr)
	}
	_ = stdout

	// Return an HTTP response, assuming the above code did not fail.
	Message := emoji.Sprint("\nScript signed successfully :white_check_mark:")
	fmt.Fprint(w, Message)

}

func (c *ServiceCommand) Synopsis() string {
	return "Start a server on port 7974 (PWSH) to listen for PowerShell script uploads to sign."
}

func (c *ServiceCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
