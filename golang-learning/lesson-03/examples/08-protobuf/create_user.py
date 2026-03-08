#!/usr/bin/env python3
"""
生成 Protobuf 格式的创建用户请求文件

使用方法:
    python create_user.py

依赖:
    pip install protobuf

生成文件:
    create_user.bin - Protobuf 二进制格式的请求文件
"""

import sys
import os

# 导入生成的 protobuf 模块
try:
    import user_pb2
except ImportError:
    print("错误: 找不到 user_pb2 模块")
    print("请先运行以下命令生成 Python protobuf 代码:")
    print("  protoc --python_out=. user.proto")
    sys.exit(1)


def create_user_request(username="testuser", email="test@example.com", age=25):
    """创建用户请求并序列化为二进制文件"""
    # 创建请求对象
    req = user_pb2.CreateUserRequest()
    req.username = username
    req.email = email
    req.age = age

    # 序列化为二进制文件
    output_file = "create_user.bin"
    with open(output_file, "wb") as f:
        f.write(req.SerializeToString())

    print(f"✓ 成功生成 Protobuf 请求文件: {output_file}")
    print(f"  用户名: {username}")
    print(f"  邮箱: {email}")
    print(f"  年龄: {age}")
    print(f"\n可以使用以下命令发送请求:")
    print(f"  curl -X POST http://localhost:8080/api/proto/user \\")
    print(f"    -H 'Content-Type: application/x-protobuf' \\")
    print(f"    -H 'Accept: application/x-protobuf' \\")
    print(f"    --data-binary @{output_file} \\")
    print(f"    --output response.bin")


if __name__ == "__main__":
    # 可以从命令行参数读取，也可以使用默认值
    if len(sys.argv) >= 4:
        username = sys.argv[1]
        email = sys.argv[2]
        age = int(sys.argv[3])
        create_user_request(username, email, age)
    else:
        # 使用默认值
        create_user_request()

