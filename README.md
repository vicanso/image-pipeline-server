# image-pipeline-server

## ENV配置

- `IP_REDIS`: redis数据库连接串，如`IP_REDIS=redis://[username:password@]host1[:port1][,...hostN[:portN]]/?poolSize=10`
- `IP_CACHE_SIZE`: 内存缓存大小(MB)，如`IP_CACHE_SIZE=100`表示设置内存缓存为100MB，如果不指定则是20MB
