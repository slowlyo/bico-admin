// 全局共享数据
import { DEFAULT_NAME } from '@/constants';
import { useState } from 'react';

const useUser = () => {
  const [name, setName] = useState<string>(DEFAULT_NAME);
  const [currentUser, setCurrentUser] = useState<API.CurrentUser | undefined>();

  return {
    name,
    setName,
    currentUser,
    setCurrentUser,
  };
};

export default useUser;
