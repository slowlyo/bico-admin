import i18n from "i18next";
import { initReactI18next } from "react-i18next";

// 中文翻译资源
const resources = {
  zh: {
    common: {
      pages: {
        login: {
          title: "登录您的账户",
          signin: "登录",
          divider: "或",
          fields: {
            username: "用户名",
            password: "密码"
          },
          errors: {
            requiredUsername: "用户名是必填项",
            requiredPassword: "密码是必填项"
          },
          buttons: {
            submit: "登录",
            rememberMe: "记住我"
          }
        },
        error: {
          info: "您可能忘记将 {{action}} 组件添加到 {{resource}} 资源中。",
          404: "抱歉，您访问的页面不存在。",
          resource404: "您确定已经创建了 {{resource}} 资源吗？",
          backHome: "返回首页"
        }
      },
      actions: {
        list: "列表",
        create: "创建",
        edit: "编辑",
        show: "查看"
      },
      buttons: {
        create: "创建",
        save: "保存",
        logout: "退出登录",
        delete: "删除",
        edit: "编辑",
        cancel: "取消",
        confirm: "您确定吗？",
        filter: "筛选",
        clear: "清除",
        refresh: "刷新",
        show: "查看",
        undo: "撤销",
        import: "导入",
        clone: "克隆",
        notAccessTitle: "您没有访问权限"
      },
      warnWhenUnsavedChanges: "您确定要离开吗？您有未保存的更改。",
      notifications: {
        success: "成功",
        error: "错误 (状态码: {{statusCode}})",
        undoable: "您有 {{seconds}} 秒时间撤销",
        createSuccess: "成功创建 {{resource}}",
        createError: "创建 {{resource}} 时出错 (状态码: {{statusCode}})",
        deleteSuccess: "成功删除 {{resource}}",
        deleteError: "删除 {{resource}} 时出错 (状态码: {{statusCode}})",
        editSuccess: "成功编辑 {{resource}}",
        editError: "编辑 {{resource}} 时出错 (状态码: {{statusCode}})",
        importProgress: "导入中: {{processed}}/{{total}}"
      },
      loading: "加载中",
      tags: {
        clone: "克隆"
      },
      dashboard: {
        title: "控制台"
      },
      table: {
        actions: "操作"
      },
      documentTitle: {
        default: "Bico Admin",
        suffix: " | Bico Admin",
        dashboard: "控制台 | Bico Admin"
      },
      autoSave: {
        success: "已保存",
        error: "自动保存失败",
        loading: "保存中...",
        idle: "等待更改"
      }
    }
  }
};

i18n
  .use(initReactI18next)
  .init({
    resources,
    lng: "zh", // 默认语言设置为中文
    fallbackLng: "zh",
    ns: ["common"],
    defaultNS: "common",
    interpolation: {
      escapeValue: false
    }
  });

export default i18n;
