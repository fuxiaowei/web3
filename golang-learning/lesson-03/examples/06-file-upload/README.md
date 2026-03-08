# 06 - 文件上传与下载

## 说明

演示如何处理文件上传、多文件上传、文件下载和静态文件服务。

## 知识点

- 单文件上传（`c.FormFile()`）
- 多文件上传（`c.MultipartForm()`）
- 保存上传的文件（`c.SaveUploadedFile()`）
- 文件下载（`c.File()`）
- 静态文件服务（`r.Static()`、`r.StaticFS()`）
- 文件路径安全处理

## 运行方式

```bash
go run main.go
```

程序会自动创建 `uploads` 目录用于存储上传的文件。

## 测试

### 1. 单文件上传

```bash
# 创建测试文件
echo "Hello, World!" > test.txt

# 上传文件
curl -X POST http://localhost:8080/api/upload \
  -F "file=@test.txt"
```

### 2. 多文件上传

```bash
# 创建多个测试文件
echo "File 1" > file1.txt
echo "File 2" > file2.txt
echo "File 3" > file3.txt

# 上传多个文件
curl -X POST http://localhost:8080/api/upload-multiple \
  -F "files=@file1.txt" \
  -F "files=@file2.txt" \
  -F "files=@file3.txt"
```

### 3. 文件下载

```bash
# 下载上传的文件
curl -O http://localhost:8080/api/download/test.txt

# 或者在浏览器中访问
# http://localhost:8080/api/download/test.txt
```

### 4. 访问静态文件

```bash
# 通过 /files 路径访问上传的文件
curl http://localhost:8080/files/test.txt

# 或在浏览器中查看
# http://localhost:8080/files/
```

## 使用 Postman 测试

### 单文件上传
1. Method: POST
2. URL: `http://localhost:8080/api/upload`
3. Body: form-data
4. Key: `file` (类型选择 File)
5. Value: 选择要上传的文件

### 多文件上传
1. Method: POST
2. URL: `http://localhost:8080/api/upload-multiple`
3. Body: form-data
4. Key: `files` (类型选择 File)
5. Value: 选择多个文件（可添加多个 `files` 字段）

## HTML 表单示例

```html
<!-- 单文件上传 -->
<form action="http://localhost:8080/api/upload" method="post" enctype="multipart/form-data">
  <input type="file" name="file">
  <button type="submit">上传</button>
</form>

<!-- 多文件上传 -->
<form action="http://localhost:8080/api/upload-multiple" method="post" enctype="multipart/form-data">
  <input type="file" name="files" multiple>
  <button type="submit">上传</button>
</form>
```

## 注意事项

1. **文件大小限制**：Gin 默认最大请求体大小为 32MB，可通过 `MaxMultipartMemory` 调整
2. **文件类型验证**：生产环境应添加文件类型和大小验证
3. **文件名安全**：使用 `filepath.Base()` 防止目录遍历攻击
4. **存储路径**：生产环境建议使用对象存储服务（如 S3、OSS）
5. **权限控制**：添加认证中间件保护上传和下载接口

## 扩展功能

可以添加的功能：
- 文件类型验证（MIME type）
- 文件大小限制
- 图片压缩和缩略图生成
- 文件去重（通过 MD5/SHA256）
- 上传进度显示
- 断点续传

