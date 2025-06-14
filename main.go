package main

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"trpc.group/trpc-go/trpc-go"
	// proto package 的路径请读者自行调整
	"hello/proto/simplest"

	_ "github.com/go-sql-driver/mysql" // 匿名导入，注册 MySQL 驱动
)

func main() {
	// 打开数据库连接池
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db?charset=utf8mb4")
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 测试连接
	if err = db.Ping(); err != nil {
		log.Fatal("数据库ping失败:", err)
	}
	s := trpc.NewServer()
	simplest.RegisterHelloWorldService(s, helloWorldImpl{})
	_ = s.Serve()
}

type helloWorldImpl struct{}

func (helloWorldImpl) Hello(ctx context.Context, req *simplest.HelloRequest) (*simplest.HelloResponse, error) {
	// 1. 从请求中获取ID
	id := req.GetId()

	// 2. 查询数据库
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db?charset=utf8mb4")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "数据库连接失败: %v", err)
	}
	defer db.Close()

	// 3. 使用上下文进行查询（支持超时取消）
	var name string
	err = db.QueryRowContext(ctx, "SELECT name FROM name WHERE id = ?", id).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "未找到ID为%d的记录", id)
		}
		return nil, status.Errorf(codes.Internal, "数据库查询失败: %v", err)
	}

	// 4. 返回响应
	return &simplest.HelloResponse{Name: name}, nil
}
