import { ModalForm } from '@ant-design/pro-components';
import type { ModalFormProps } from '@ant-design/pro-components';
import { message } from 'antd';
import React, { useEffect, useState } from 'react';

export interface CrudModalProps<T = any, CreateParams = any, UpdateParams = any>
  extends Omit<ModalFormProps, 'onFinish' | 'title'> {
  /** 弹窗标题前缀，如 "用户"，会自动生成 "新建用户" / "编辑用户" */
  title: string;
  /** 当前编辑的记录，有值时为编辑模式，无值时为创建模式 */
  record?: T;
  /** 创建 API */
  onCreate?: (values: CreateParams) => Promise<API.Response<any>>;
  /** 更新 API */
  onUpdate?: (id: number, values: UpdateParams) => Promise<API.Response<any>>;
  /** 成功回调 */
  onSuccess?: () => void;
  /** 将记录转换为表单初始值 */
  recordToValues?: (record: T) => any;
  /** 表单内容 */
  children: React.ReactNode;
}

function CrudModal<T extends { id: number } = any>({
  title,
  record,
  onCreate,
  onUpdate,
  onSuccess,
  recordToValues,
  children,
  open,
  onOpenChange,
  ...rest
}: CrudModalProps<T>) {
  const isEdit = !!record;
  const [initialValues, setInitialValues] = useState<any>({});

  useEffect(() => {
    if (open && record) {
      const values = recordToValues ? recordToValues(record) : record;
      setInitialValues(values);
    } else if (!open) {
      setInitialValues({});
    }
  }, [open, record, recordToValues]);

  const handleFinish = async (values: any) => {
    try {
      let res: API.Response<any>;
      if (isEdit && onUpdate) {
        res = await onUpdate(record.id, values);
      } else if (onCreate) {
        res = await onCreate(values);
      } else {
        return false;
      }

      if (res.code === 0) {
        message.success(isEdit ? '更新成功' : '创建成功');
        onSuccess?.();
        return true;
      }
      message.error(res.msg || (isEdit ? '更新失败' : '创建失败'));
      return false;
    } catch (error: any) {
      message.error(error.message || error.data?.msg || '操作失败');
      return false;
    }
  };

  return (
    <ModalForm
      title={isEdit ? `编辑${title}` : `新建${title}`}
      width={600}
      open={open}
      onOpenChange={onOpenChange}
      initialValues={initialValues}
      onFinish={handleFinish}
      modalProps={{
        destroyOnHidden: true,
      }}
      {...rest}
    >
      {children}
    </ModalForm>
  );
}

export default CrudModal;
