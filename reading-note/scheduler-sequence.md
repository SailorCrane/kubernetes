sequenceDiagram

cmd/scheduler->> app/server: NewSchedulerCommand()
    Note right of app/server : 设置命令行选项到option中, 并创建一个pflag cmd出来
app/server->> cmd/scheduler: return a "cmd" object

cmd/scheduler ->> app/server: cmd.Execute()
app/server ->> app/server: RunCommand(cmd, options)
    app/server ->> app/options: opts.Config()
        Note right of app/options : 将options转换为scheduler's config
        Note right of app/options : config内含kube client(用来连接kube-apiserver)
    app/options ->> app/server: return config
    app/server ->> pkg/scheduler/algorithm: algorithmprovider.ApplyFeatureGates()
        Note right of pkg/scheduler/algorithm: 根据feature gate选项, 配置scheduler开启的选项

app/server ->> app/server: Run(cc/completedConfig, stopCh)
        app/server ->> pkg/scheduler/scheduler: scheduler.New(cc.PodInformer, cc.InformerFactory)
            Note right of pkg/scheduler/scheduler: scheduler后期会执行scheduleOne()
            Note right of pkg/scheduler/scheduler: 创建一个Scheduler对象pkg/scheduler/scheduler出来
            Note right of pkg/scheduler/scheduler: Scheduler对象其实就是包含一个config(config由FactoryConfig创建)
            Note right of pkg/scheduler/scheduler: FactoryConfig包含podInformer和相应pod添加时的回调函数
            Note right of pkg/scheduler/scheduler: FactoryConfig's podInformer回调函数addPodToSchedulingQueue()
            Note right of pkg/scheduler/scheduler: addPodToSchedulingQueue() 添加 unassigned到podQueue中(informer添加回调时有filter过滤器),
            Note right of pkg/scheduler/scheduler: 这里过滤没有名字的pod(没有到达final state的pod)
        pkg/scheduler/scheduler ->> app/server: return sched

        app/server ->> pkg/scheduler/scheduler: sched.Run()
            Note right of app/server : 在协程中运行 sched.Run()
            Note right of app/server : sched.Run() 循环调用 scheduleOne()
            Note right of app/server : shceduleOne就是真正把pod调度到node的代码所在
            pkg/scheduler/scheduler->> pkg/scheduler/scheduler: sched.scheduleOne()



Note right of nothing1 : 只是为了mermaid显示完整所加, 不是源代码内容
Note right of nothing2 : 只是为了mermaid显示完整所加, 不是源代码内容
Note right of nothing3 : 只是为了mermaid显示完整所加, 不是源代码内容
