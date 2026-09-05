package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"bico-admin/internal/projectinit"

	"github.com/spf13/cobra"
)

var (
	initModule string
	initName   string
	initTitle  string
	initDryRun bool
	initYes    bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "将本仓库身份替换为新项目名称",
	Long: `拷贝 bico-admin 后，用新的 module / 二进制名 / 展示名替换仓库内残留身份。

示例:
  go run ./cmd init --module github.com/acme/foo --name foo --title "Foo Admin"
  go run ./cmd init --dry-run --module github.com/acme/foo --name foo --yes`,
	SilenceUsage: true,
	RunE:         runInit,
}

func init() {
	initCmd.Flags().StringVar(&initModule, "module", "", "Go module 路径，例如 github.com/acme/foo")
	initCmd.Flags().StringVar(&initName, "name", "", "二进制 / 镜像 / package.json name，例如 foo")
	initCmd.Flags().StringVar(&initTitle, "title", "", "展示标题，例如 Foo Admin")
	initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "只打印将改路径，不写盘")
	initCmd.Flags().BoolVar(&initYes, "yes", false, "非交互跳过确认")
	rootCmd.AddCommand(initCmd)
}

// runInit 收集身份、确认后执行清名。缺 flag 时在交互终端补齐。
func runInit(cmd *cobra.Command, args []string) error {
	in := cmd.InOrStdin()
	out := cmd.OutOrStdout()

	id, err := resolveIdentity(in, out)
	if err != nil {
		return err
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "module: %s\nname:   %s\ntitle:  %s\ndb:     %s\nenv:    %s_\n",
		id.Module, id.Name, id.Title, id.SnakeName(), id.EnvPrefix())

	if !initYes && !initDryRun {
		ok, err := confirm(in, out, "确认按以上身份改写当前仓库? [y/N] ")
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("已取消")
		}
	}

	result, err := projectinit.Run(root, id, projectinit.Secrets{}, initDryRun)
	if err != nil {
		return err
	}

	if len(result.Files) == 0 {
		fmt.Fprintln(out, "没有需要改写的文件")
		return nil
	}

	fmt.Fprintln(out, "将改写的文件:")
	for _, file := range result.Files {
		fmt.Fprintln(out, "  "+file)
	}
	fmt.Fprintf(out, "共 %d 个文件\n", len(result.Files))
	if result.DryRun {
		fmt.Fprintln(out, "dry-run: 未写盘")
		return nil
	}
	fmt.Fprintln(out, "已生成新的 jwt.secret 与前端 crypto 密钥")
	return nil
}

// resolveIdentity 用 flags 与交互输入拼出完整身份。
// --yes 时不允许再提问，缺必填项直接失败。
func resolveIdentity(in io.Reader, out io.Writer) (projectinit.Identity, error) {
	id := projectinit.Identity{
		Module: initModule,
		Name:   initName,
		Title:  initTitle,
	}

	needPrompt := id.Module == "" || id.Name == "" || id.Title == ""
	if needPrompt {
		if initYes {
			normalized, err := id.Normalize()
			if err != nil {
				return projectinit.Identity{}, fmt.Errorf("非交互模式需要完整参数: %w", err)
			}
			return normalized, nil
		}
		if !isInteractive(in) {
			normalized, err := id.Normalize()
			if err != nil {
				return projectinit.Identity{}, fmt.Errorf("非交互模式需要 --module，或补齐 --name/--title: %w", err)
			}
			return normalized, nil
		}
		prompted, err := promptIdentity(in, out, id)
		if err != nil {
			return projectinit.Identity{}, err
		}
		id = prompted
	}
	return id.Normalize()
}

// promptIdentity 在终端逐项询问缺失字段。
func promptIdentity(in io.Reader, out io.Writer, id projectinit.Identity) (projectinit.Identity, error) {
	scanner := bufio.NewScanner(in)
	module, err := promptLine(scanner, out, "Go module", id.Module)
	if err != nil {
		return projectinit.Identity{}, err
	}
	id.Module = module

	defaultName := id.Name
	if defaultName == "" {
		if normalized, err := id.Normalize(); err == nil {
			defaultName = normalized.Name
		}
	}
	name, err := promptLine(scanner, out, "项目名 (二进制/镜像)", defaultName)
	if err != nil {
		return projectinit.Identity{}, err
	}
	id.Name = name

	defaultTitle := id.Title
	if defaultTitle == "" {
		if normalized, err := id.Normalize(); err == nil {
			defaultTitle = normalized.Title
		}
	}
	title, err := promptLine(scanner, out, "展示标题", defaultTitle)
	if err != nil {
		return projectinit.Identity{}, err
	}
	id.Title = title
	return id, nil
}

// promptLine 读取一行输入，空回车使用默认值。
func promptLine(scanner *bufio.Scanner, out io.Writer, label, def string) (string, error) {
	if def != "" {
		fmt.Fprintf(out, "%s [%s]: ", label, def)
	} else {
		fmt.Fprintf(out, "%s: ", label)
	}
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		if def != "" {
			return def, nil
		}
		return "", fmt.Errorf("未输入 %s", label)
	}
	line := strings.TrimSpace(scanner.Text())
	if line == "" {
		return def, nil
	}
	return line, nil
}

// confirm 读取 y/yes 作为确认，其他输入视为取消。
func confirm(in io.Reader, out io.Writer, question string) (bool, error) {
	fmt.Fprint(out, question)
	scanner := bufio.NewScanner(in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}
		return false, nil
	}
	switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}

// isInteractive 判断 stdin 是否是终端，管道输入时不再追问。
func isInteractive(in io.Reader) bool {
	file, ok := in.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
