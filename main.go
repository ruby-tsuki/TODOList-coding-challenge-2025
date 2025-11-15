package main

import (
	"J/DAO"
	"J/model"
	"J/service"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	// MySQL配置 - 请根据您的实际情况修改
	dataSourceName := "root:kongming123@tcp(localhost:3306)/todolist?parseTime=true"

	// 初始化DAO
	taskDAO, err := DAO.NewMySQLTaskDAO(dataSourceName)
	if err != nil {
		fmt.Printf("❌ 初始化MySQL失败: %v\n", err)
		return
	}

	// 初始化服务
	todoService := service.NewTodoService(taskDAO)
	defer todoService.Close()

	go DDLCheck(todoService)

	//整活
	fmt.Printf(`⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣦⡀⠀⠀⣠⣿⣿⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣠⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣤⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀⠀
⠀⠀⣰⣿⡿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠁⠀⠀⠀⠀
⠀⠀⣿⣿⣄⠈⢿⣿⣿⣿⣿⡿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠁⠀⠀⠀⠀⠀
⠀⠀⠈⠻⢿⣿⣿⠏⠉⠉⠉⠀⠀⠀⠈⠙⠻⠛⠃⠈⠛⠛⠉⠁⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
=======================================================
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⣿⢠⣾⠋⠀⠀⣿⡇⠀⣿⠀⠀⠀⣿⠛⣿⡆⠀⢀⣾⠟⠛⣷⡄⠀⠀⠀
⠀⠀⠀⣿⠻⣧⡀⠀⠀⣿⡇⢀⣿⠀⠀⠀⣿⠻⣯⡀⠀⠸⣿⡀⢀⣿⠇⠀⠀⠀
⠀⠀⠀⠛⠀⠙⠓⠀⠀⠈⠛⠛⠋⠀⠀⠀⠛⠀⠙⠓⠀⠀⠙⠛⠛⠋⠀⠀⠀⠀
=======================================================
`)
	fmt.Println("🎯 TodoList 应用")
	fmt.Println("命令说明: add, undo, urgent, done, update, delete, finish, deleteAll, clear, exit")

	runCLI(todoService)
}

func runCLI(service *service.TodoService) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n请输入命令: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "add":
			handleAddCommand(service, parts)
		case "undo":
			handleUndoCommand(service, parts)
		case "urgent":
			handleUrgentCommand(service, parts)
		case "done":
			handleDoneCommand(service)
		case "update":
			handleUpdateCommand(service, parts)
		case "delete":
			handleDeleteCommand(service, parts)
		case "finish":
			handleFinishCommand(service, parts)
		case "deleteAll":
			handleDeleteAllCommand(service)
		case "clear":
			handleClearCommand()
		case "exit":
			fmt.Println("再见!")
			return
		case "help":
			displayHelp()
		default:
			fmt.Println("未知命令，输入 'help' 查看可用命令")
		}
	}
}

func DDLCheck(service *service.TodoService) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			checkOneHourUrgentTasks(service)
		}
	}
}

func handleAddCommand(service *service.TodoService, parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: add <任务内容> <过期时间（可选）>")
		return
	}
	title := ""
	ddl := 24 * 60 //默认一天过期
	if len(parts) > 3 {
		title = strings.Join(parts[1:len(parts)-1], " ") //AI，还得是AI考虑的细
		ddlStr := parts[len(parts)-1]
		ddlInt, err := strconv.Atoi(ddlStr)
		if err != nil {
			fmt.Println("分钟数必须是整")
			return
		}
		if ddlInt < 0 {
			fmt.Println("分钟数不能为负数")
			return
		}
		ddl = ddlInt
	} else {
		title = parts[1]
	}

	task, err := service.AddTask(title, ddl)
	if err != nil {
		fmt.Printf("添加任务失败: %v\n", err)
	} else {
		fmt.Printf("✅ 任务添加成功! ID: %d\n", task.ID)
	}
}

// A!
func handleUndoCommand(service *service.TodoService, parts []string) {
	limit := 10
	if len(parts) > 1 {
		if l, err := strconv.Atoi(parts[1]); err == nil && l > 0 {
			limit = l
		}
	}

	tasks, err := service.GetRecentUndoTasks(limit)
	if err != nil {
		fmt.Printf("获取未完成任务失败: %v\n", err)
		return
	}
	fmt.Printf("\n🔄 最近 %d 个未完成任务:\n", limit)
	displayTasks("", tasks)
}

func handleUrgentCommand(service *service.TodoService, parts []string) {
	limit := 5
	if len(parts) > 1 {
		if l, err := strconv.Atoi(parts[1]); err == nil && l > 0 {
			limit = l
		}
	}

	tasks, err := service.GetUrgentTasks(limit)
	if err != nil {
		fmt.Printf("获取紧迫任务失败: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("🎯 没有紧迫的DDL任务")
	} else {
		fmt.Printf("\n🚨 最紧迫的 %d 个DDL任务:\n", limit)
		displayTasksWithDeadline(tasks)
	}
}

func checkOneHourUrgentTasks(service *service.TodoService) {
	tasks, err := service.GetUrgentTasks(5)
	if err != nil {
		fmt.Printf("检查DDL时获取任务列表失败:%v\n", err)
		return
	}
	now := time.Now()
	warningTime := time.Hour

	expiringTasks := make([]*model.Task, 0)
	for _, task := range tasks {
		timeUntilDeadline := task.DeadLine.Sub(now)
		if timeUntilDeadline > 0 && timeUntilDeadline <= warningTime {
			expiringTasks = append(expiringTasks, task)
		}
	}
	if len(expiringTasks) > 0 {
		fmt.Println("\n🚨🚨🚨 DDL 警报！以下任务将在1小时内到期：")
		fmt.Println("=========================================")
		for i, task := range expiringTasks {
			minutesLeft := int(task.DeadLine.Sub(now).Minutes())
			fmt.Printf("%d. [ID:%d] %s\n", i+1, task.ID, task.Title)
			fmt.Printf("   剩余时间: %d分钟 | 到期时间: %s\n", minutesLeft, task.DeadLine.Format("15:04:05"))
			fmt.Println()
		}
		fmt.Println("=========================================")
	}
}

func handleDoneCommand(service *service.TodoService) {
	tasks, err := service.ShowDoneTasks()
	if err != nil {
		fmt.Printf("获取已完成任务失败: %v\n", err)
		return
	}
	displayTasks("✅ 已完成任务", tasks)
}

func handleUpdateCommand(service *service.TodoService, parts []string) {
	// 命令参数长度校验：允许 4 个（不更新DDL）或 5 个（更新DDL）参数
	if len(parts) < 4 || len(parts) > 5 {
		fmt.Println("用法: update <任务ID> <新标题> <完成状态(true/false)> [相对当前的分钟数(可选，用于更新DDL)]")
		return
	}

	taskID, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("任务ID必须是数字")
		return
	}

	newTitle := parts[2]

	done, err := strconv.ParseBool(parts[3])
	if err != nil {
		fmt.Println("完成状态必须是 true 或 false")
		return
	}

	var ddl time.Time // 默认为零值（表示不更新DDL）
	if len(parts) == 5 {
		// 存在分钟数参数，解析为整数
		minutes, err := strconv.Atoi(parts[4])
		if err != nil {
			fmt.Println("分钟数必须是整数（例如：30 表示30分钟后）")
			return
		}
		// 校验分钟数非负（避免设置过去的时间，根据业务需求可调整）
		if minutes < 0 {
			fmt.Println("分钟数不能为负数（请输入相对于当前时间的未来分钟数）")
			return
		}
		// 计算DDL：当前时间 + 分钟数
		ddl = time.Now().Add(time.Duration(minutes) * time.Minute)
	}

	err = service.UpdateTask(taskID, newTitle, done, ddl)
	if err != nil {
		fmt.Printf("更新任务失败: %v\n", err)
	} else {
		if len(parts) == 5 {
			fmt.Printf("📝 任务 %d 更新成功（新DDL：%s）\n", taskID, ddl.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("📝 任务 %d 更新成功（未修改DDL）\n", taskID)
		}
	}
}
func handleDeleteCommand(service *service.TodoService, parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: delete <任务ID>")
		return
	}
	taskID, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("任务ID必须是数字")
		return
	}
	err = service.DeleteTask(taskID)
	if err != nil {
		fmt.Printf("删除任务失败: %v\n", err)
	} else {
		fmt.Printf("🗑️  任务 %d 已删除\n", taskID)
	}
}

func handleFinishCommand(service *service.TodoService, parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: finish <任务ID>")
		return
	}
	taskID, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("任务ID必须是数字")
		return
	}
	err = service.FinishedTask(taskID)
	if err != nil {
		fmt.Printf("标记任务完成失败: %v\n", err)
	} else {
		fmt.Printf("✅ 任务 %d 已完成\n", taskID)
	}
}

func handleDeleteAllCommand(service *service.TodoService) {
	fmt.Print("确定要删除所有任务吗？(y/N): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		confirm := strings.TrimSpace(scanner.Text())
		if strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes" {
			err := service.ClearAllTasks()
			if err != nil {
				fmt.Printf("删除所有任务失败: %v\n", err)
			} else {
				fmt.Println("🗑️  所有任务已删除")
			}
		} else {
			fmt.Println("取消删除操作")
		}
	}
}

func handleClearCommand() {
	// 清空终端屏幕
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		// 如果不支持清屏，至少输出一些空行
		fmt.Print("\033[2J\033[H")
	}
}

func displayTasks(title string, tasks []*model.Task) {
	if len(tasks) == 0 {
		fmt.Println("📝 当前没有任务")
		return
	}

	if title != "" {
		fmt.Printf("\n%s:\n", title)
	}
	fmt.Println(strings.Repeat("-", 60))
	for _, task := range tasks {
		status := "❌"
		if task.Done {
			status = "✅"
		}
		fmt.Printf("%s [%d] %s (创建: %s)\n",
			status, task.ID, task.Title,
			task.CreateAt.Format("2006-01-02 15:04"))
	}
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("总计: %d 个任务\n", len(tasks))
}

func displayTasksWithDeadline(tasks []*model.Task) {
	if len(tasks) == 0 {
		fmt.Println("📝 当前没有任务")
		return
	}

	fmt.Println(strings.Repeat("-", 80))
	now := time.Now()

	for _, task := range tasks {
		status := "❌"
		if task.Done {
			status = "✅"
		}

		// 计算剩余时间
		var timeInfo string
		if !task.DeadLine.IsZero() {
			if task.DeadLine.Before(now) {
				// 已过期
				overdue := now.Sub(task.DeadLine)
				timeInfo = fmt.Sprintf("(已过期 %v)", formatDuration(overdue))
			} else {
				// 未过期
				remaining := task.DeadLine.Sub(now)
				timeInfo = fmt.Sprintf("(剩余 %v)", formatDuration(remaining))
			}
		}

		fmt.Printf("%s [%d] %s\n", status, task.ID, task.Title)
		if !task.DeadLine.IsZero() {
			fmt.Printf("   📅 DDL: %s %s\n",
				task.DeadLine.Format("2006-01-02 15:04"), timeInfo)
		}
		fmt.Printf("   🕒 创建: %s\n", task.CreateAt.Format("2006-01-02 15:04"))
		fmt.Println()
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("总计: %d 个任务\n", len(tasks))
}

// 辅助函数：格式化时间间隔
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "不到1分钟"
	} else if d < time.Hour {
		return fmt.Sprintf("%.0f分钟", d.Minutes())
	} else if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if minutes > 0 {
			return fmt.Sprintf("%d小时%d分钟", hours, minutes)
		}
		return fmt.Sprintf("%d小时", hours)
	} else {
		days := int(d.Hours() / 24)
		hours := int(d.Hours()) % 24
		if hours > 0 {
			return fmt.Sprintf("%d天%d小时", days, hours)
		}
		return fmt.Sprintf("%d天", days)
	}
}

func displayHelp() {
	fmt.Println(`
可用命令:
  add <任务内容> <ddl(可选)>       - 添加新任务
  undo [数量]          			 - 显示最近未完成任务（按创建时间）
  urgent [数量]        			 - 显示最紧迫的DDL任务
  done                 			 - 显示已完成任务
  update <ID> <标题> <状态> <ddl>  - 更新任务(标题和状态)
  delete <ID>         			 - 删除指定任务
  finish <ID>         			 - 标记任务为已完成
  deleteAll          			 - 删除所有任务
  clear               			 - 清空终端屏幕
  exit                			 - 退出程序
  help                			 - 显示此帮助信息

示例:
  add 学习Go语言        # 添加任务
  undo                 # 查看最近10个未完成任务
  undo 5               # 查看最近5个未完成任务
  urgent               # 查看最紧迫的5个DDL任务
  urgent 3             # 查看最紧迫的3个DDL任务
  finish 1             # 将ID为1的任务标记为完成
  update 1 "学习Golang" true  # 更新ID为1的任务
  delete 1             # 删除ID为1的任务
  deleteAll            # 删除所有任务
  clear                # 清空终端屏幕
	`)
}
