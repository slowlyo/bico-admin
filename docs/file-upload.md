# 文件上传功能说明

## 功能概述

admin模块提供了通用的文件上传功能，支持多文件上传、目录分类存储、文件类型验证等特性。

## 接口信息

- **路由**: `POST /admin/upload`
- **认证**: 需要登录认证，无需特殊权限
- **内容类型**: `multipart/form-data`

## 请求参数

### 表单参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| files | file[] | 是 | 上传的文件，支持多文件 |
| dir | string | 否 | 上传目录，不指定则存储在根目录 |

### URL参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| dir | string | 否 | 上传目录，与表单参数二选一 |

## 文件限制

### 文件大小
- 单个文件最大: **可配置**（默认10MB）
- 配置项: `upload.max_file_size`

### 文件数量
- 单次上传最大文件数: **可配置**（默认10个）
- 配置项: `upload.max_files`

### 支持的文件类型
- **可配置**：通过 `upload.allowed_types` 配置项设置
- **默认支持**：

#### 图片类型
- `.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`, `.webp`, `.svg`

#### 文档类型
- `.pdf`, `.doc`, `.docx`, `.xls`, `.xlsx`, `.ppt`, `.pptx`, `.txt`, `.csv`

#### 压缩文件
- `.zip`, `.rar`, `.7z`, `.tar`, `.gz`

#### 音视频文件
- `.mp3`, `.mp4`, `.avi`, `.mov`, `.wmv`, `.flv`, `.wav`

## 存储规则

### 目录结构
```
data/uploads/
├── [dir]/                    # 指定目录（可选）
│   ├── 20240715_143022_a1b2c3d4.jpg
│   └── 20240715_143025_e5f6g7h8.pdf
└── 20240715_143030_i9j0k1l2.txt  # 根目录文件
```

### 文件命名规则
- 格式: `{时间戳}_{UUID前8位}.{原扩展名}`
- 示例: `20240715_143022_a1b2c3d4.jpg`
- 避免文件名冲突，确保唯一性

## 响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "files": [
      {
        "file_name": "example.jpg",
        "file_path": "/uploads/images/20240715_143022_a1b2c3d4.jpg",
        "file_url": "http://localhost:8899/uploads/images/20240715_143022_a1b2c3d4.jpg",
        "file_size": 1024000,
        "file_type": "图片"
      }
    ],
    "total": 1
  }
}
```

### 部分成功响应
```json
{
  "code": 200,
  "message": "成功上传 1 个文件，1 个文件上传失败: large.mp4: 文件大小超过限制（最大10MB）",
  "data": {
    "files": [
      {
        "file_name": "small.jpg",
        "file_path": "/uploads/20240715_143022_a1b2c3d4.jpg",
        "file_size": 512000,
        "file_type": "图片"
      }
    ],
    "total": 1
  }
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "请选择要上传的文件"
}
```

## 文件访问

上传成功后，文件可通过以下方式访问：

### 推荐方式：使用 file_url
```javascript
// 直接使用响应中的完整URL（推荐）
const fileUrl = response.data.files[0].file_url;

// 在img标签中使用
<img src={fileUrl} alt="上传的图片" />
```

### 兼容方式：使用 file_path
```javascript
// 获取上传响应中的file_path
const filePath = response.data.files[0].file_path;

// 构建完整访问URL
const fileUrl = `${window.location.origin}${filePath}`;

// 在img标签中使用
<img src={fileUrl} alt="上传的图片" />
```

### 直接访问
- URL格式: `http://localhost:8899{file_path}`
- 示例: `http://localhost:8899/uploads/images/20240715_143022_a1b2c3d4.jpg`

## 使用示例

### JavaScript/Fetch
```javascript
const formData = new FormData();
formData.append('files', file1);
formData.append('files', file2);
formData.append('dir', 'images');

fetch('/admin/upload', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer your-token-here'
  },
  body: formData
})
.then(response => response.json())
.then(data => {
  console.log('上传成功:', data);
});
```

### cURL
```bash
curl -X POST \
  -H "Authorization: Bearer your-token-here" \
  -F "files=@/path/to/file1.jpg" \
  -F "files=@/path/to/file2.pdf" \
  -F "dir=documents" \
  http://localhost:8899/admin/upload
```

## 安全特性

1. **路径遍历防护**: 自动清理目录参数，防止 `../` 等路径遍历攻击
2. **文件类型验证**: 基于文件扩展名进行白名单验证
3. **文件大小限制**: 防止大文件上传导致服务器资源耗尽
4. **认证要求**: 需要有效的登录token才能上传文件
5. **唯一文件名**: 自动生成唯一文件名，避免文件覆盖

## 配置说明

文件上传功能支持通过配置文件进行自定义设置，配置项位于 `config/app.yml` 的 `upload` 部分：

```yaml
upload:
  max_file_size: "10MB"        # 单个文件最大大小: 支持 KB, MB, GB 单位
  max_files: 10                # 单次上传最大文件数量
  allowed_types:               # 允许的文件类型（扩展名）
    - ".jpg"
    - ".jpeg"
    - ".png"
    # ... 更多类型
  upload_dir: "data/uploads"   # 上传文件存储目录
  base_url: ""                 # 文件访问的基础URL，留空则自动使用请求的域名
```

### 配置项说明

#### base_url 配置
- **留空（默认）**: 自动使用请求的域名构建文件URL
  - 开发环境: `http://localhost:8899/uploads/xxx.jpg`
  - 生产环境: `https://api.example.com/uploads/xxx.jpg`
- **自定义域名**: 指定文件访问的基础URL
  - 示例: `base_url: "https://cdn.example.com"`
  - 结果: `https://cdn.example.com/uploads/xxx.jpg`
- **适用场景**:
  - 使用CDN加速文件访问
  - 前后端分离部署，需要指定后端域名
  - 文件服务器与API服务器分离

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `max_file_size` | string | "10MB" | 单个文件最大大小，支持KB/MB/GB单位 |
| `max_files` | int | 10 | 单次上传最大文件数量，0表示不限制 |
| `allowed_types` | []string | 见配置文件 | 允许的文件扩展名列表 |
| `upload_dir` | string | "data/uploads" | 文件存储根目录 |

## 注意事项

1. 上传的文件存储在配置的 `upload_dir` 目录（默认 `data/uploads/`）
2. 确保服务器有足够的磁盘空间存储上传文件
3. 生产环境建议配置反向代理（如Nginx）来处理静态文件服务
4. 定期清理不需要的上传文件以节省存储空间
5. 修改配置后需要重启服务才能生效
6. 文件类型验证基于文件扩展名，建议结合其他安全措施
