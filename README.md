# mysql
基于xorm封装的MySQL库


# 功能
 - 简化MySQL实例初始化流程，基于配置自动对MySQL进行初始化
 - 支持连接池、多实例等场景
 - 使用该库，建议再看看ant-coder项目，可以基于数据表的字段和索引生成响应的model代码


# 基本使用
 - toml 配置文件
    ```
    [mysql.default]
        user = "root"
        pawd = "123456"
        host = "127.0.0.1"
        port = "3306"
        name = "business"
    [mysql.stats]
        user = "root"
        pawd = "123456"
        host = "127.0.0.1"
        port = "3306"
        name = "business"
    ```

 - 使用方法
	```golang
    // 初始化config包，参考config模块
    code...

    // 验证mysql实例的配置正确性与连通性。非必须
    if err = mysql.Valid(); err != nil {
        fmt.Printf("mysql error: %s\n", err)
        os.Exit(-1)
    }

    // 如下方式可以直接使用MySQL实例
    mysql.DefaultClient().SQL("SELECT * FROM t").Find(&rows)
    mysql.Client("default").SQL("SELECT * FROM t").Find(&rows)
    ```
