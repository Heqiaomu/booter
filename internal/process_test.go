package internal

import (
	"fmt"
	"github.com/agiledragon/gomonkey"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/google/gops/goprocess"
	"github.com/olekukonko/tablewriter"
)

func init(){
	resp, err := http.Get("http://nexus.hyperchain.cn/repository/blocface/blocface/booter/blocface-booter-test")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = ioutil.WriteFile("test/blocface-test",b,777)
	if err != nil{
		panic(err)
	}
}

func TestGops(t *testing.T) {
	as := goprocess.FindAll()
	fmt.Println("====", as)
	fmt.Println("start : ", StartAt("87746"))
}

func tp() {
	data := [][]string{
		{"blocface-core", "running", "2233", "$10.98"},
		{"blocface-gateway", "down", "-", "$54.95"},
		{"blocface-pay", "down", "-", "$51.00"},
		{"blocface-monitor", "running", "2233", "$30.00"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"App", "State", "Pid", "StartAt"})
	// table.SetFooter([]string{"", "", "Total", "$146.93"}) // Add Footer
	table.SetBorder(false) // Set Border to false
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}


type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) String() string {
	return string(*s)
}


func Test_getProcess(t *testing.T) {
	type args struct {
		c *cobra.Command
	}
	tests := []struct {
		name  string
		args  args
		wantP *process
	}{
		// TODO: Add test cases.
		{
			name: "test-ok",
		},
	}

	f := &pflag.Flag{}
	var dir string
	var work string
	gomonkey.ApplyMethod(reflect.TypeOf(&cobra.Command{}),"Flag", func(command *cobra. Command, cmd string) (flag *pflag.Flag){
		if cmd == "dir" {
			f.Value = newStringValue("test",&dir)
		}else if cmd == "work" {
			f.Value = newStringValue("test",&work)
		}
		return f
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getProcess(tt.args.c)
		})
	}
}

func Test_process_filter(t *testing.T) {

	f := &pflag.Flag{}
	var dir string
	var work string
	gomonkey.ApplyMethod(reflect.TypeOf(&cobra.Command{}),"Flag", func(command *cobra. Command, cmd string) (flag *pflag.Flag){
		if cmd == "dir" {
			f.Value = newStringValue("test",&dir)
		}else if cmd == "work" {
			f.Value = newStringValue("test",&work)
		}
		return f
	})

	type args struct {
		names []string
	}
	tests := []struct {
		name    string
		args    args
		wantApsNum int
	}{
		// TODO: Add test cases.
		{
			name: "test-app ok",
			args: args{names: []string{"test"}},
			wantApsNum: 1,
		},
		{
			name: "test-app all",
			args: args{names: []string{"all"}},
			wantApsNum: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := getProcess(&cobra.Command{})
			if gotAps := p.filter(tt.args.names); len(gotAps) != tt.wantApsNum {
				t.Errorf("filter() = %v, want num:%v", gotAps, tt.wantApsNum)
			}
		})
	}
}

func Test_process_restart(t *testing.T) {
	f := &pflag.Flag{}
	var dir string
	var work string
	gomonkey.ApplyMethod(reflect.TypeOf(&cobra.Command{}),"Flag", func(command *cobra. Command, cmd string) (flag *pflag.Flag){
		if cmd == "dir" {
			f.Value = newStringValue("test",&dir)
		}else if cmd == "work" {
			f.Value = newStringValue("test",&work)
		}
		return f
	})

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := getProcess(&cobra.Command{})
			if err := p.restart(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("restart() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(1111)
		})
	}
}

func Test_process_start(t *testing.T) {

	f := &pflag.Flag{}
	var dir string
	var work string
	gomonkey.ApplyMethod(reflect.TypeOf(&cobra.Command{}),"Flag", func(command *cobra. Command, cmd string) (flag *pflag.Flag){
		if cmd == "dir" {
			f.Value = newStringValue("test",&dir)
		}else if cmd == "work" {
			f.Value = newStringValue("test",&work)
		}
		return f
	})

	type args struct {
		console bool
		args    []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := getProcess(&cobra.Command{})
			if err := p.start(tt.args.console, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_process_status(t *testing.T) {

	f := &pflag.Flag{}
	var dir string
	var work string
	gomonkey.ApplyMethod(reflect.TypeOf(&cobra.Command{}),"Flag", func(command *cobra. Command, cmd string) (flag *pflag.Flag){
		if cmd == "dir" {
			f.Value = newStringValue("test",&dir)
		}else if cmd == "work" {
			f.Value = newStringValue("test",&work)
		}
		return f
	})

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := getProcess(&cobra.Command{})
			if err := p.status(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("status() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_process_stop(t *testing.T) {

	f := &pflag.Flag{}
	var dir string
	var work string
	gomonkey.ApplyMethod(reflect.TypeOf(&cobra.Command{}),"Flag", func(command *cobra. Command, cmd string) (flag *pflag.Flag){
		if cmd == "dir" {
			f.Value = newStringValue("test",&dir)
		}else if cmd == "work" {
			f.Value = newStringValue("test",&work)
		}
		return f
	})

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-ok",
		},
		{
			name: "test-stop-ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			p := getProcess(&cobra.Command{})
			if tt.name == "test-stop-ok" {
				p.start(false,nil)
			}
			if err := p.stop(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_process_tail(t *testing.T) {
	f := &pflag.Flag{}
	var dir string
	var work string
	gomonkey.ApplyMethod(reflect.TypeOf(&cobra.Command{}),"Flag", func(command *cobra. Command, cmd string) (flag *pflag.Flag){
		if cmd == "dir" {
			f.Value = newStringValue("test",&dir)
		}else if cmd == "work" {
			f.Value = newStringValue("test",&work)
		}
		return f
	})
	gomonkey.ApplyMethod(reflect.TypeOf(&exec.Cmd{}),"Run", func(cmd *exec.Cmd) error {return nil})
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-need app name ",
			wantErr: true,
		},
		{
			name: "test-ok",
			args: args{
				args: []string{"test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := getProcess(&cobra.Command{})
			if err := p.tail(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("tail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

