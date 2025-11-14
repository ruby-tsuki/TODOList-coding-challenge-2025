package main

import (
	"J/DAO"
	"J/model"
	"J/service"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	fmt.Println("🎯 TodoList 应用")
	fmt.Println("命令说明: add, undo, done, update, delete, finish, clear, recent, exit")

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
			handleUndoCommand(service)
		case "done":
			handleDoneCommand(service)
		case "update":
			handleUpdateCommand(service, parts)
		case "delete":
			handleDeleteCommand(service, parts)
		case "finish":
			handleFinishCommand(service, parts)
		case "clear":
			handleClearCommand(service)
		case "recent":
			handleRecentCommand(service, parts)
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

func handleAddCommand(service *service.TodoService, parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: add <任务内容>")
		return
	}
	title := strings.Join(parts[1:], " ")
	task, err := service.AddTask(title)
	if err != nil {
		fmt.Printf("添加任务失败: %v\n", err)
	} else {
		fmt.Printf("✅ 任务添加成功! ID: %d\n", task.ID)
	}
}

func handleUndoCommand(service *service.TodoService) {
	tasks, err := service.ShowUndoTasks()
	if err != nil {
		fmt.Printf("获取未完成任务失败: %v\n", err)
		return
	}
	displayTasks("🔄 未完成任务", tasks)
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
	if len(parts) < 4 {
		fmt.Println("用法: update <任务ID> <新标题> <完成状态(true/false)>")
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

	err = service.UpdateTask(taskID, newTitle, done)
	if err != nil {
		fmt.Printf("更新任务失败: %v\n", err)
	} else {
		fmt.Printf("📝 任务 %d 更新成功\n", taskID)
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

func handleClearCommand(service *service.TodoService) {
	fmt.Print("确定要清空所有任务吗？(y/N): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		confirm := strings.TrimSpace(scanner.Text())
		if strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes" {
			err := service.ClearAllTasks()
			if err != nil {
				fmt.Printf("清空任务失败: %v\n", err)
			} else {
				fmt.Println("🗑️  所有任务已清空")
			}
		} else {
			fmt.Println("取消清空操作")
		}
	}
}

func handleRecentCommand(service *service.TodoService, parts []string) {
	limit := 10
	if len(parts) > 1 {
		if l, err := strconv.Atoi(parts[1]); err == nil && l > 0 {
			limit = l
		}
	}

	tasks, err := service.GetRecentUndoTasks(limit)
	if err != nil {
		fmt.Printf("获取最近未完成任务失败: %v\n", err)
		return
	}

	fmt.Printf("\n🔄 最近 %d 个未完成任务:\n", limit)
	displayTasks("", tasks)
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
		fmt.Printf("%s [%d] %s (创建: %s, 更新: %s)\n",
			status, task.ID, task.Title,
			task.CreateAt.Format("2006-01-02 15:04"),
			task.UpdateAt.Format("2006-01-02 15:04"))
	}
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("总计: %d 个任务\n", len(tasks))
}

func displayHelp() {
	fmt.Println(`
可用命令:
  add <任务内容>        - 添加新任务
  undo                 - 显示未完成任务
  done                 - 显示已完成任务
  update <ID> <标题> <状态> - 更新任务(标题和状态)
  delete <ID>          - 删除任务
  finish <ID>          - 标记任务为已完成
  clear                - 清空所有任务
  recent [数量]        - 显示最近未完成任务
  exit                 - 退出程序
  help                 - 显示此帮助信息

示例:
  add 学习Go语言        # 添加任务
  undo                 # 查看未完成任务
  finish 1             # 将ID为1的任务标记为完成
  update 1 "学习Golang" true  # 更新ID为1的任务
  delete 1             # 删除ID为1的任务
  recent 5             # 显示最近5个未完成任务
	`)
}
