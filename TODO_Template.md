# TODO List 项目说明文档（示例模板）

## 1. 技术选型
- **方案初衷**：TODOlist是基于终端、单主机，单用户。因此不需要网络以及高并发上的需求。
- **编程语言**：GO。熟悉并且初步设计中需要使用到goroutine 
- **框架/库**：使用DAO设计模式，数据库操作和业务分离。并且使用interface便于拓展其它数据库。  
- **数据库/存储**：mysql。字段设计很适合用于todolist，并且可以本地持久化，而且还可以实现排序等操作。 
- 替代方案对比：reids，不适合多字段，除非使用hashmap。快是快，而且像是todolist这种小规模数据也不用考虑内存不够。但是其本地化保存不比mysql方便，而且mysql还提供诸多方法（如排序、模糊查询等）

## 2. 项目结构设计
- 整体架构说明:  
- 目录结构：  
  ```
  src/
    DAO/
      /task_dao_interface.go
      /task_dao_mysql.go
    model/
      /task.go
    service/
      /service.go
    main,go
    go.mod
    TODO_Template.md

  ```  
- 模块职责说明:
- model层-数据模型定义：定义数据模块
- DAO层-数据访问层：task_dao_interface.go定义接口，方法标签Create、GetList、Update、Delete、Count、Close。task_dao_mysql.go是基于mysql数据库的接口实现。
- service层-业务逻辑层：任务管理：AddTask、ShowUndoTasks、ShowDoneTasks。任务更新和删除：UpdateTask、DeleteTask、FinishedTask。批量操作：ClearAllTask、GetRecentUndoTasks。连接释放：Close
- main-主程序：命令行界面(CLI)的实现、用户交互和输入处理
## 3. 需求细节与决策
- 任务的添加、删除、标记完成、查询未完成/已完成任务。
- 实现全部删除、行数显示限制
- 使用mysql，可实现持久化

## 4. AI 使用说明
- 使用到chatGPT、Deepseek 
- 使用 AI 的环节：  
  - main函数中终端显示以及输入处理
  - 文档初稿编写
  - 代码检查及优化建议
- AI 输出如何修改：AI设计的GetList是整体显示，按照大多数人使用习惯我将其查分成查询undo和done。

## 5. 运行与测试方式
- 本地终端运行
  - 执行go run .\main.go
  - 输入help查询使用说明以及示例
- 已测试过的环境（windows）。  
- 已知问题与不足
  - getlist时id倒叙
  - 没有设计ddl，应该添加ddl字段，配合watchdog实现自动提醒
  - 应该设计一个终端显示清除功能

## 6. 总结与反思
- 如果有更多时间，你会如何改进？ 
  - 根据上面的问题和不足进行版本优化
- 你觉得这个实现的最大亮点是什么？
  - 基于DAO模式设计，项目结构清晰，后期维护以及拓展方便

