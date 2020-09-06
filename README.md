# Golang 公用类库

## 日志类库

### 说明

日志类库是对 Uber 的 [zap](https://github.com/uber-go/zap) 日志库进行了一层封装。

所有组件统一使用该类库来记录日志。 

### 使用方式

1. 首选使用方式

    ```
    # 初始化一个 logger，底层实际上是 zap 的 SugaredLogger
    logger := log.NewLogger()
    defer logger.FlushLogger()
    
    logger.Infow("test", "foo", "bar")
    ```

2. 特殊情况使用（主要为了兼容之前已经写好的代码）

    ```
    # 初始化一个 logger，底层实际上是 zap 的 Logger
    logger := log.NewOriginLogger()
    defer logger.FlushLogger()
    
    logger.Debug("test", zap.String("foo", "bar"))
    ```

### 配置

#### 1.设置日志级别

在配置文件或环境变量中给以下任一 key 赋值即可：

- LOG_LEVEL
- logLevel

可选的值为：（按从低到高级别排序）

`debug/info/warning/error/dpanic/panic/fatal`

如果不赋值则默认为 `debug` 级别。

#### 2.设置日志文件路径

在配置文件或环境变量中给以下任一 key 赋值即可：

- LOG_FILE
- logFile

如果不赋值则表示不写入到文件。

## MongoDB

### 1. 配置

使用之前需要配置以下变量的值，配置方式不限，能通过 [viper](https://github.com/spf13/viper) 类库获取到即可。推荐配置在系统环境变量中。

- MONGODB_USER：用户名
- MONGODB_PASSWORD：密码
- MONGODB_HOST：主机域名或 IP
- MONGODB_PORT：MongoDB 端口
- MONGODB_SSL：是否使用 https 连接。可选值为 true 或 false，字符串。
- MONGODB_DBNAME：默认连接的数据库名称。也可以通过 `SetDatabase` 方法设置。
- MONGODB_CONN_TIMEOUT：连接 MongoDB 的超时时间
- MONGODB_OP_TIMEOUT：各个方法的操作超时时间

### 2. 初始化一个连接

```
mongo, err := mongodb.NewMongoClient(logger)
if err != nil {
    panic(err)
}
mongo.SetDatabase("amf")
```

### 3. 各个方法的具体使用方式请见方法的注释

## DB (GORM)

### 1. 配置

使用之前需要配置以下变量的值，配置方式不限，能通过 [viper](https://github.com/spf13/viper) 类库获取到即可。推荐配置在系统环境变量中。

- DB_ENGINE：数据库类型。可选值为为 `mysql/postgres/sqlite3/mssql` 等。
- DB_USER：用户名
- DB_PASSWORD：密码
- DB_HOST：主机域名或 IP
- DB_PORT：MongoDB 端口
- DB_CHARSET：数据库字符集
- DB_NAME：默认连接的数据库名称

### 2. 初始化一个连接

```
dbClient, err := db.NewDB(*logger)
if err != nil {
    panic(err)
}
```

### 3. 操作数据库的具体方式请见官方文档

- https://gorm.io/
- https://github.com/jinzhu/gorm
