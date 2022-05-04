# image-pipeline-server

image pipeline支持各种方便快捷的图片处理任务的pipeline处理，


图片压缩的处理依赖于[tiny](https://github.com/vicanso/tiny)，有需要单独启用此服务。 

```bash
docker run -d --restart=always \
  -p 6002:7002 \
  --name=tiny \
  vicanso/tiny
```

## ENV配置

- `IP_REDIS`: redis数据库连接串，如`IP_REDIS=redis://[username:password@]host1[:port1][,...hostN[:portN]]/?poolSize=10`，如果有配置redis连接，则会使用redis作为二级缓存，缓存处理后的图片数据
- `IP_CACHE_SIZE`: 内存缓存大小(MB)，如`IP_CACHE_SIZE=100`表示设置内存缓存为100MB，如果不指定则是20MB
- 以`IP_FINDER_`开头的环境变量，表示设置的image pipeline的finder
- 其它以`IP_`开头的环境变量，表示设置image pipeline任务的别名

## 启用服务

```bash
docker run -d --restart=always \
  -p 8001:7001 \
  -e IP_tiny=optimize/172.16.23.175:6002 \
  -e IP_FINDER_oss="aliyun://oss-cn-beijing.aliyuncs.com?accessKey=key&secretKey=secret" \
  -e IP_FINDER_minio="minio://172.16.214.137:9000/?accessKey=test&secretKey=testabcd" \
  -e IP_FINDER_gridfs="mongodb://test:testabcd@172.16.214.137:27017/admin" \
  --name image-pipeline \
  vicanso/image-pipeline-server
```

- `IP_tiny=optimize/172.16.23.175:6002` 指定一个名为`tiny`的图片优化任务别名，它使用`172.16.23.175:6002`的tiny服务处理图片优化
- `IP_FINDER_oss="aliyun://oss-cn-beijing.aliyuncs.com?accessKey=key&secretKey=secret"` 指定一个名为`oss`的阿里云oss存储服务，其中`key`与`secret`需要调整oss对应的配置
- `IP_FINDER_minio="minio://172.16.214.137:9000/?accessKey=test&secretKey=testabcd"` 指定一个名为minio的minio oss存储服务，其中`key`与`secret`需要调整oss对应的配置
- `IP_FINDER_gridfs="mongodb://test:testabcd@172.16.214.137:27017/admin"` 指定一个名为gridfs的mongodb gridfs存储服务，其中用户名、密码与db需要调整为对应的配置

获取图片并转换为webp，质量为90，对应的请求如下：

```bash
http://127.0.0.1:8001/?proxy/https%3A%2F%2Fwww.baidu.com%2Fimg%2FPCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png|tiny/90/webp
```

从oss服务中获取图片(bucket:tinysite)，转换为webp，质量为90，对应的请求如下：

```bash
http://127.0.0.1:8001/?oss/tinysite/go-echarts.jpg|tiny/90/webp
```

从minio服务中获取图片(bucket:bigtree)，转换为webp，质量为90，对应的请求如下：

```bash
http://127.0.0.1:8001/?minio/bigtree/go-charts.png|tiny/90/webp
```

从mongodb gridfs服务中获取图片(bucket:bigtree)，转换为webp，质量为90，对应的请求如下：

```bash
http://127.0.0.1:8001/?gridfs/6242c7f1e08b32ac7b550673|tiny/80/webp
```
