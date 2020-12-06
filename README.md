#### memstore
通过sync.map简单实现一个内存数据库，主要实现的操作有
- Get(key) 根据key获取一个值
- Put(key, value, ttl) 设置一个值，并设置ttl
- Garbage 实现简单的过期数据回收的机制