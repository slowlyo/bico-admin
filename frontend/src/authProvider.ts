import type { AuthProvider } from "@refinedev/core";
import { authAPI, setToken, clearAuth } from "./utils/request";

export const TOKEN_KEY = "refine-auth";

export const authProvider: AuthProvider = {
  login: async ({ username, email, password }) => {
    try {
      // 调用后端登录API
      const response = await authAPI.login({
        username: username || email,
        password,
      });

      // 后端返回格式: { code: 200, message: "Success", data: { user, token } }
      if (response.code === 200 && response.data) {
        // 存储token和用户信息
        setToken(response.data.token);
        localStorage.setItem(TOKEN_KEY, response.data.token);
        localStorage.setItem("user", JSON.stringify(response.data.user));

        return {
          success: true,
          redirectTo: "/",
        };
      }

      return {
        success: false,
        error: {
          name: "LoginError",
          message: response.message || "登录失败",
        },
      };
    } catch (error: any) {
      console.error("Login error:", error);
      return {
        success: false,
        error: {
          name: "LoginError",
          message: error.message || "网络错误，请稍后重试",
        },
      };
    }
  },
  logout: async () => {
    try {
      // 调用后端登出API
      await authAPI.logout();
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      // 无论API调用是否成功，都清除本地存储
      clearAuth();
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem("user");
    }

    return {
      success: true,
      redirectTo: "/login",
    };
  },
  check: async () => {
    const token = localStorage.getItem(TOKEN_KEY);
    if (token) {
      return {
        authenticated: true,
      };
    }

    return {
      authenticated: false,
      redirectTo: "/login",
    };
  },
  getPermissions: async () => null,
  getIdentity: async () => {
    const token = localStorage.getItem(TOKEN_KEY);
    if (token) {
      try {
        // 尝试从localStorage获取用户信息
        const userStr = localStorage.getItem("user");
        if (userStr) {
          const user = JSON.parse(userStr);
          return {
            id: user.id,
            name: user.nickname || user.username,
            email: user.email,
            avatar: user.avatar || "https://i.pravatar.cc/300",
          };
        }

        // 如果本地没有用户信息，尝试从API获取
        const response = await authAPI.getProfile();
        if (response.code === 200 && response.data) {
          localStorage.setItem("user", JSON.stringify(response.data));
          return {
            id: response.data.id,
            name: response.data.nickname || response.data.username,
            email: response.data.email,
            avatar: response.data.avatar || "https://i.pravatar.cc/300",
          };
        }
      } catch (error) {
        console.error("Get identity error:", error);
        // 如果获取用户信息失败，清除认证状态
        clearAuth();
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem("user");
      }
    }
    return null;
  },
  onError: async (error) => {
    console.error(error);
    return { error };
  },
};
