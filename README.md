# TodoList 应用说明
一个基于 Go 语言的命令行任务管理工具
## 功能特性
- 添加、更新、删除任务
- DDL自动提醒（自动检查1小时内到期任务）
- 多样化查询方式（基于添加顺序、基于DDL、查看未完成、查看已完成）
- 任务状态跟踪（未完成/已完成）
- 批量任务清理
- 数据持久化
## 技术栈
- 后端: Go
- 数据库: MySQL
- 架构: DAO + Service 分层架构
## 使用方法
|命令 | 说明 | 示例 |
|--|--|--|
| add <内容> [ddl] | 添加任务（ddl单位为分钟，不设置默认24h） |add “秋招” 60|
| undo [数量] | 显示未完成任务 |undo|
| urgent [数量] | 显示紧迫的DDL任务	 |urgent 3|
| done | 显示已完成任务	 |done|
| update <ID> <标题> <状态> [ddl]	 | 更新任务	 |update 1 "新标题" true 30|
| delete <ID> | 删除指定任务	 |delete 1|
| finish <ID>	 | 标记任务完成	 |finish 1|
| deleteAll | 删除所有任务	 |deleteAll|
| clear | 清空终端屏幕	 |clear|
| exit | 退出程序	 |exit|
| help | 显示帮助信息		 |help|
## 其它
启动界面有猫猫画面


