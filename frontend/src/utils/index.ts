/**
 * 格式化时间
 */
export const formatTime = (time: string | Date) => {
  if (!time) return '';
  const date = new Date(time);
  return date.toLocaleString('zh-CN');
};

/**
 * 去除字符串首尾空格
 */
export const trim = (str: string) => {
  return str ? str.trim() : '';
};

/**
 * 检查是否为空值
 */
export const isEmpty = (value: any) => {
  return value === null || value === undefined || value === '';
};

/**
 * 获取文件扩展名
 */
export const getFileExtension = (filename: string) => {
  return filename.slice(((filename.lastIndexOf('.') - 1) >>> 0) + 2);
};

/**
 * 格式化文件大小
 */
export const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};
