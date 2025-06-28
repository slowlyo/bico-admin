import type { RefineThemedLayoutV2HeaderProps } from "@refinedev/antd";
import { useGetIdentity, useLogout, useTranslate } from "@refinedev/core";
import {
  Avatar,
  Layout as AntdLayout,
  Space,
  Switch,
  theme,
  Typography,
  Popover,
  Button,
} from "antd";
import { LogoutOutlined, UserOutlined, SettingOutlined } from "@ant-design/icons";
import React, { useContext } from "react";
import { useNavigate } from "react-router";
import { ColorModeContext } from "../../contexts/color-mode";
import "./header.css";

const { Text } = Typography;
const { useToken } = theme;

type IUser = {
  id: number;
  name: string;
  nickname?: string;
  email?: string;
  avatar: string;
};

export const Header: React.FC<RefineThemedLayoutV2HeaderProps> = ({
  sticky = true,
}) => {
  const { token } = useToken();
  const { data: user } = useGetIdentity<IUser>();
  const { mode, setMode } = useContext(ColorModeContext);
  const { mutate: logout } = useLogout();
  const translate = useTranslate();
  const navigate = useNavigate();

  // 创建用户弹出菜单内容
  const userMenuContent = (
    <div>
      <Button
        type="text"
        size="middle"
        icon={<SettingOutlined />}
        onClick={() => navigate("/profile")}
        style={{
          width: "100%",
          textAlign: "left",
          justifyContent: "flex-start",
          marginBottom: "4px",
        }}
      >
        个人资料
      </Button>
      <Button
        type="text"
        size="middle"
        icon={<LogoutOutlined />}
        onClick={() => logout()}
        style={{
          width: "100%",
          textAlign: "left",
          justifyContent: "flex-start",
        }}
      >
        {translate("buttons.logout", "退出登录")}
      </Button>
    </div>
  );

  const headerStyles: React.CSSProperties = {
    backgroundColor: token.colorBgElevated,
    display: "flex",
    justifyContent: "flex-end",
    alignItems: "center",
    padding: "0px 24px",
    height: "64px",
  };

  if (sticky) {
    headerStyles.position = "sticky";
    headerStyles.top = 0;
    headerStyles.zIndex = 1;
  }

  return (
    <AntdLayout.Header style={headerStyles}>
      <Space>
        <Switch
          checkedChildren="🌛"
          unCheckedChildren="🔆"
          onChange={() => setMode(mode === "light" ? "dark" : "light")}
          defaultChecked={mode === "dark"}
        />
        {user && (
          <Popover
            content={userMenuContent}
            placement="bottomRight"
            trigger="click"
            arrow={false}
          >
            <div
              style={{
                marginLeft: "8px",
                cursor: "pointer",
                padding: "4px 8px",
                borderRadius: "4px",
                transition: "background-color 0.2s",
                height: "40px",
                display: "flex",
                alignItems: "center",
                gap: "8px",
              }}
              className="user-dropdown-trigger"
            >
              <Avatar
                size={32}
                src={user.avatar}
                alt={user.nickname || user.name}
                icon={!user.avatar && <UserOutlined />}
              />
              {user.name && <Text strong>{user.nickname || user.name}</Text>}
            </div>
          </Popover>
        )}
      </Space>
    </AntdLayout.Header>
  );
};
