/**
 * 简单的加密/解密工具
 * 用于记住密码功能
 */

const SECRET_KEY = 'bico-admin-secret-key';

/**
 * 加密字符串
 */
export function encrypt(text: string): string {
  try {
    // 使用 Base64 + 简单异或加密
    const encrypted = btoa(
      text
        .split('')
        .map((char, i) => 
          String.fromCharCode(char.charCodeAt(0) ^ SECRET_KEY.charCodeAt(i % SECRET_KEY.length))
        )
        .join('')
    );
    return encrypted;
  } catch (e) {
    console.error('加密失败:', e);
    return '';
  }
}

/**
 * 解密字符串
 */
export function decrypt(encrypted: string): string {
  try {
    if (!encrypted) return '';
    
    const decrypted = atob(encrypted)
      .split('')
      .map((char, i) => 
        String.fromCharCode(char.charCodeAt(0) ^ SECRET_KEY.charCodeAt(i % SECRET_KEY.length))
      )
      .join('');
    return decrypted;
  } catch (e) {
    console.error('解密失败:', e);
    return '';
  }
}

/**
 * 保存记住的账号密码
 */
export function saveCredentials(username: string, password: string) {
  try {
    localStorage.setItem('remembered_username', encrypt(username));
    localStorage.setItem('remembered_password', encrypt(password));
  } catch (e) {
    console.error('保存凭证失败:', e);
  }
}

/**
 * 获取记住的账号密码
 */
export function getCredentials(): { username: string; password: string } | null {
  try {
    const encryptedUsername = localStorage.getItem('remembered_username');
    const encryptedPassword = localStorage.getItem('remembered_password');
    
    if (!encryptedUsername || !encryptedPassword) {
      return null;
    }
    
    return {
      username: decrypt(encryptedUsername),
      password: decrypt(encryptedPassword),
    };
  } catch (e) {
    console.error('获取凭证失败:', e);
    return null;
  }
}

/**
 * 清除记住的账号密码
 */
export function clearCredentials() {
  try {
    localStorage.removeItem('remembered_username');
    localStorage.removeItem('remembered_password');
  } catch (e) {
    console.error('清除凭证失败:', e);
  }
}
