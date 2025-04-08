# dev-helper

## 项目日志记录

1、使用gin初始化了项目，新建的 main.go、go.mod, 下一步完善书籍展示功能。  --- 2025.03.26 23:0

2、增加 redis client, 定时任务。  --- 2025.04.06 22:19

### 三方包

1、redis操作使用 `github.com/redis/go-redis/v9` 包
2、sql操作使用 `gorm.io/gorm` 包


### go-gin使用手册
- https://gin-gonic.com/zh-cn/docs/


### brew redis 操作

- brew services start redis        #启动redis服务
- brew services stop redis         #关闭redis服务
- brew services restart redis      #重启redis服务
