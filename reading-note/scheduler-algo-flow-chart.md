
graph TB
    st(开始) --> predicate(预选)
    predicate --> parallel_predicate(go协程开启, 分析所有node)
    predicate --> priority(统计选择出node的分数)
