package imgcat

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

//Print print image
func Print(read io.ReadCloser) {
	cmd := exec.Command("imgcat")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, read)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)

}
