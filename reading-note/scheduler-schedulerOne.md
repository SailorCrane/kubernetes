sequenceDiagram

scheduler ->> factory :sched.config.nextPod()
    Note right of factory: pkg/scheduler/scheduler pkg/scheduler/factory/factory
    Note right of factory : config从ConfigFactory中来
    Note right of factory : podInformer和podQueue在ConfigFactory中
    Note right of factory : nextPod从nextPod中读取一个unsigned pod
factory ->> scheduler :return a unassigned pod
scheduler ->> scheduler: shced.schedule(pod)

scheduler ->> generic_scheduler : sched.config.Algorithm.Schedule(pod, sched.config.NodeLister)
    Note right of generic_scheduler : pkg/scheduler/core/generic_scheduler
    Note right of generic_scheduler : 为pod选择一个合适的node

    generic_scheduler ->> generic_scheduler : g.findNodesThatFit(pod, nodes)
        Note right of generic_scheduler : 先过滤出符合条件的node(cpu, 内容等符合条件)
        Note right of generic_scheduler : 即"predicate"步骤
    generic_scheduler ->> generic_scheduler : g.priorityMetaProducer(pod, g.cachedNodeInfoMap)
    generic_scheduler ->> generic_scheduler : PrioritizeNodes(pod, cachedNodeInfoMap, metaPrioritiesInterface)
        Note right of generic_scheduler : "priority"步骤
    generic_scheduler ->> generic_scheduler : g.selectHost(priorityList)
        Note right of generic_scheduler : 选择最合适的node主机, 并返回

generic_scheduler ->> scheduler : return host
generic_scheduler ->> metrics : metrics.SchedulingAlgorithmLatency.Observe(metrics.SinceInMicroseconds(start))
        Note right of metrics : 猜测是用来做profile/log/predicate等工作的(用于后期优化算法等等)

generic_scheduler ->> generic_scheduler : sched.assumeVolumes(assumedPod, scheduleResult.SuggestedHost)
        Note right of generic_scheduler : 猜测是用来为pod绑定volume
generic_scheduler ->> generic_scheduler : sched.assume(assumedPod, scheduleResult.SuggestedHost)
        Note right of generic_scheduler : 设置pod的NodeName为SuggestedHost

        generic_scheduler ->> pkg/scheduler/internal/cache/cache : sched.config.SchedulerCache.AssumePod(assumed)
            Note right of generic_scheduler : 在cache中记录 assumedPod, 表示正在bind中

generic_scheduler ->> generic_scheduler : sched.bind(assumedPod, v1.Binding(Name: host))
generic_scheduler ->> factory : sched.config.GetBinder(assumed).Bind(b)
    factory ->> generic_scheduler : getBinderFunc().Bind()
generic_scheduler ->> generic_scheduler : getBinderFunc().Bind()
        Note right of generic_scheduler : 这里把pod绑定到node
        Note right of generic_scheduler : 绑定失败的话, 会清除cache中的assumed
        Note right of generic_scheduler : 绑定失败的话, 也会在 recordSchedulingFailure()发送失败事件
        Note right of generic_scheduler : 绑定成功, 会通过sched.config.Record发送成功消息

        Note right of generic_scheduler : sched.config.Recorder.Eventf() // 绑定成功的消息

generic_scheduler ->> histogram : metrics.E2eSchedulingLatency.Observe(metrics.SinceInMicroseconds(start))
            Note right of histogram : 根据统计信息, 画直方图. 用来做profile
    

Note right of nothing1 : 只是为了mermaid显示完整所加, 不是源代码内容
Note right of nothing2 : 只是为了mermaid显示完整所加, 不是源代码内容
Note right of nothing3 : 只是为了mermaid显示完整所加, 不是源代码内generic_scheduler容
