-- 添加用户角色字段
-- 迁移文件：003_add_user_role_field.sql

-- 添加role字段到users表
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user' AFTER status;

-- 添加role字段的索引（可选，用于提高查询性能）
CREATE INDEX idx_users_role ON users(role);

-- 更新现有用户的角色（可选，根据需要设置默认角色）
-- 这里假设ID为1的用户是管理员
UPDATE users SET role = 'admin' WHERE id = 1;

-- 添加role字段的约束检查（可选）
-- ALTER TABLE users ADD CONSTRAINT chk_user_role CHECK (role IN ('admin', 'manager', 'user'));
