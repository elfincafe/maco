package mmaco

import (
	"fmt"
	"math"
	"strings"
)

type (
	help struct {
		Desc string
	}
)

func (cmd *help) Init() {
	cmd.Desc = "this help."
}

func (cmd *help) Validate() error {
	return nil
}

func (cmd *help) Run(ctx *Context) error {
	// Either Root Command or Sub Command
	name := ""
	if ctx.NumArg() > 0 {
		arg := ctx.Arg(0)
		if _, ok := ctx.subCmds[arg]; ok {
			name = ctx.subCmds[arg].Name
		}
	}

	line := ""
	sb := new(strings.Builder)
	sb.WriteString("Usage:\n")
	if name != "" {
		line = fmt.Sprintf("    %s [options] %s [sub command options] [arg] ...\n", ctx.cmd.Name, name)
	} else {
		line = fmt.Sprintf("    %s [options] <sub command> [sub command options] [arg] ...\n", ctx.cmd.Name)
	}
	sb.WriteString(line)
	sb.WriteString("\nOptions:\n")

	// Root Options
	max := 0
	for _, o := range ctx.cmd.opts {
		max = int(math.Max(float64(max), float64(len(o.Long))))
	}
	format := fmt.Sprintf("    -%%-s, --%%-%ds   %%s\n", max)
	for _, o := range ctx.cmd.opts {
		sb.WriteString(fmt.Sprintf(format, o.Short, o.Long, o.Desc))
	}

	// Sub Commands
	if name == "" {
		max += 4
		sb.WriteString("\nSub Commands:\n")
		for _, sc := range ctx.scOrder {
			max = int(math.Max(float64(max), float64(len(sc))))
		}
		format = fmt.Sprintf("    %%-%ds   %%s\n", max)
		for _, sc := range ctx.scOrder {
			sb.WriteString(fmt.Sprintf(format, sc, ctx.subCmds[sc].Desc))
		}
	}

	// Sub Command Options
	if name != "" {
		sb.WriteString("\nSub Command Options:\n")
		max = 0
		for _, o := range ctx.subCmds[name].opts {
			max = int(math.Max(float64(max), float64(len(o.Long))))
		}
		format := fmt.Sprintf("    -%%-s, --%%-%ds   %%s\n", max)
		for _, o := range ctx.subCmds[name].opts {
			sb.WriteString(fmt.Sprintf(format, o.Short, o.Long, o.Desc))
		}
	}

	println(sb.String())
	return nil
}
