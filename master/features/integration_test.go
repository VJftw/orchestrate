package integration

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/features/user"
	"github.com/vjftw/orchestrate/master/features/utils"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	Convey("Building test binary", t, func() {
		cmd := exec.Command("go", "build", "-o", "integrationTest", "../")
		stderr, err := cmd.StderrPipe()
		if err != nil {
			t.Error(err)
		}
		So(err, ShouldBeNil)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Error(err)
		}
		So(err, ShouldBeNil)

		err = cmd.Start()
		if err != nil {
			t.Error(err)
		}
		So(err, ShouldBeNil)

		io.Copy(os.Stdout, stdout)
		errBuf, _ := ioutil.ReadAll(stderr)

		err = cmd.Wait()
		if err != nil {
			t.Error(string(errBuf))
		}
		So(err, ShouldBeNil)

		if _, err = os.Stat("integrationTest"); os.IsNotExist(err) {
			t.Error(err)
		}
		So(err, ShouldBeNil)

		Convey("Starting test binary", func() {
			cmd := exec.Command("./integrationTest")
			err := cmd.Start()
			if err != nil {
				t.Error(err)
			}
			So(err, ShouldBeNil)

			apiClient := utils.APIClient{}
			apiClient.BaseURI = "http://localhost:8734"

			userManagement := user.UserManagementTests{}

			userManagement.RunTests(t, apiClient)

			Convey("Stopping test binary", func() {
				cmd.Process.Kill()
			})
		})
	})

}
