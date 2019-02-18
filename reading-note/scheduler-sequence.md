sequenceDiagram

cmd/scheduler->> app/server: NewSchedulerCommand()
Note right of app/server : 设置命令行选项到option中, 并创建一个pflag cmd出来
app/server->> cmd/scheduler: return cmd

cmd/scheduler ->> app/server: cmd.Execute()
app/server ->> app/server: RunCommand(cmd, options)
    app/server -> app/options: opts.Config()
Note right of app/options : 将options转换为scheduler的config
