package tui

import (
	"fmt"
	"strings"
	"sort"
	"encoding/json"
	"github.com/google/shlex"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	value "mc/value"
)

type MxTui struct {
	app *tview.Application
	inputs []string
	input_desc map[string]string
	outputs []string
	output_desc map[string]string
	ioview *tview.TextView
	prompt *tview.InputField
	node string
	cmd_channel chan string
	done_channel chan bool
}

func MakeTui(iface string, cmd_ch chan string, done_ch chan bool) *MxTui {
	lines := strings.Split(iface, "\n")
	outdesc := make(map[string]string)
	indesc := make(map[string]string)
	inputs := make([]string, 0)
	outputs := make([]string, 0)
	node := strings.TrimSpace(lines[0])
	
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		elts := strings.SplitN(line," ",3)
		if elts[0] == "input" {
			indesc[elts[1]] = elts[2]
			inputs = append(inputs, elts[1])
		} else if elts[0] == "output" {
			outdesc[elts[1]] = elts[2]
			outputs = append(outputs, elts[1])
		}
	}
	sort.Strings(inputs)
	sort.Strings(outputs)
	t := &MxTui {
		app: tview.NewApplication(),
		inputs: inputs,
		outputs: outputs,
		input_desc: indesc,
		output_desc: outdesc,
		node: node,
		cmd_channel: cmd_ch,
		done_channel: done_ch}

	data_input := tview.NewInputField()
	data_input.SetLabel("> ")
	data_input.SetFieldWidth(0)
	//data_input.SetAutocompleteFunc(t.Complete)
	data_input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			data := data_input.GetText()
			if data == "quit" || data == "exit" {
				t.done_channel <- true
				t.app.Stop()
			} else if data == "help" || data == "?" {
				t.GotInput("help", "")
				data_input.SetText("")
				t.GotInternal(`mx shell:
help|? -- print this message
quit|exit -- exit\n
<cmd> <arg1> <val1> <arg2> <val2> ... -- send command <cmd> with arguments arg1=val1, arg2=val2, ...`)
			} else {
				t.Enter()
			}
		}
	})
	
	ifspec_lines := make([]string, 0)
	ifspec_lines = append(ifspec_lines, "Inputs: ")
	for _, i := range t.inputs {
		ifspec_lines = append(ifspec_lines, fmt.Sprintf("- %s %s",i,t.input_desc[i]))
	}
	ifspec_lines = append(ifspec_lines, "Outputs: ")
	for _, o := range t.outputs {
		ifspec_lines = append(ifspec_lines, fmt.Sprintf("- %s %s",o,t.output_desc[o]))
	}
	ifspec := strings.Join(ifspec_lines, "\n")
	
	header_view := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("%s: connected", t.node))
	io_view := tview.NewTextView().SetTextAlign(tview.AlignLeft).SetText("").SetDynamicColors(true)
	help_view := tview.NewTextView().SetTextAlign(tview.AlignLeft).SetText(ifspec)
	
	grid := tview.NewGrid()
	grid.SetRows(1, 0, 1) // rows: header, middle, input 
	grid.SetColumns(0, 0) // cols: data, help
	grid.SetBorders(true)

	grid.AddItem(header_view, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(io_view, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(help_view, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(data_input, 2, 0, 1, 2, 0, 0, false)
	t.app.SetRoot(grid, true).SetFocus(data_input)
	t.prompt = data_input
	t.ioview = io_view
	return t
}

func (t *MxTui) Enter() {
	inp := t.prompt.GetText()
	toks, err := shlex.Split(inp)
	if err != nil {
		t.GotError(fmt.Sprintf("Failed to split input: %s",inp))
		return
	}
	if len(toks) == 0 {
		return
	}
	cmd := toks[0]
	args := toks[1:]
	if len(args) % 2 != 0 {
		t.GotError("ERROR: Odd number of args. Expected <args> of the form `arg val arg val ...`")
		return
	}

	arg_val, _ := value.FromObject(make(map[string]interface{}))
	arg_name := ""
	for i, arg := range args {
		if i%2 == 0 {
			arg_name = arg
		} else {
			v, err := value.Parse(arg)
			if err != nil {
				v = value.MakeStringValue(arg)
			}
			arg_val.MapVal[arg_name] = v
		}
	}

	data := make(map[string]interface{})
	data["operation"] = "send"
	data["command"] = cmd
	data["args"] = arg_val.ToObject()
	bs, err := json.Marshal(data)
	if err != nil {
		t.GotError(fmt.Sprintf("ERROR: Failed to serialize data: %v\n",data))
		return
	}
	//fmt.Println("Sending",string(bs))
	t.cmd_channel <- string(bs)
	t.GotInput(cmd, arg_val.ToPrettyString())
	t.prompt.SetText("")
}

func (t *MxTui) GotOutput(command, args string) {
	arg := args
	d := map[string]interface{}{}
	json.Unmarshal([]byte(args), &d)
	arg_val, err := value.FromObject(d)
	if err == nil {
		arg = arg_val.ToPrettyString()
	} else {
		arg = fmt.Sprintf("%s (%v)", arg, err)
	}
	t.ioview.Write([]byte(fmt.Sprintf("< [yellow]%s %s[white]\n",command, arg)))
	t.app.Draw()
}
	
func (t *MxTui) GotInternal(data string) {
	t.ioview.Write([]byte(fmt.Sprintf("< [gray]%s[white]\n",data)))
}
	
func (t *MxTui) GotError(data string) {
	t.ioview.Write([]byte(fmt.Sprintf("! [red]%s[white]\n",data)))
}

func (t *MxTui) GotInput(command, args string) {
	t.ioview.Write([]byte(fmt.Sprintf("> [green]%s %s[white]\n",command, args)))
}
	
func (t *MxTui) Complete(prefix string) []string {
	ans := make([]string, 0)
	for _, inp := range t.inputs {
		if strings.HasPrefix(inp, prefix) {
			ans = append(ans, inp)
		}
	}
	return ans
}

func (t *MxTui) Run() {
	if err := t.app.Run(); err != nil {
		panic(err)
	}
	fmt.Println("done")
}
