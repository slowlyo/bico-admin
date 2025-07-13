import { ModalForm, ProFormText, ProFormTextArea, ProFormSelect } from '@ant-design/pro-components';
import { Form } from 'antd';
import React from 'react';
import { RoleCreateRequest } from '@/services/role';

export interface CreateFormProps {
  modalVisible: boolean;
  onCancel: () => void;
  onSubmit: (values: RoleCreateRequest) => Promise<void>;
}

const CreateForm: React.FC<CreateFormProps> = (props) => {
  const { modalVisible, onCancel, onSubmit } = props;
  const [form] = Form.useForm();

  return (
    <ModalForm
      title="新建角色"
      width="400px"
      form={form}
      open={modalVisible}
      onOpenChange={(visible) => {
        if (!visible) {
          onCancel();
        }
      }}
      onFinish={async (value) => {
        await onSubmit(value);
        form.resetFields();
      }}
      modalProps={{
        destroyOnClose: true,
      }}
    >
      <ProFormText
        name="code"
        label="角色代码"
        placeholder="请输入角色代码"
        rules={[
          {
            required: true,
            message: '角色代码为必填项',
          },
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/,
            message: '角色代码只能包含字母、数字和下划线，且以字母开头',
          },
          {
            min: 2,
            max: 50,
            message: '角色代码长度为2-50个字符',
          },
        ]}
      />
      
      <ProFormText
        name="name"
        label="角色名称"
        placeholder="请输入角色名称"
        rules={[
          {
            required: true,
            message: '角色名称为必填项',
          },
          {
            max: 100,
            message: '角色名称长度不能超过100个字符',
          },
        ]}
      />
      
      <ProFormTextArea
        name="description"
        label="描述"
        placeholder="请输入角色描述"
        fieldProps={{
          rows: 4,
        }}
        rules={[
          {
            max: 500,
            message: '描述长度不能超过500个字符',
          },
        ]}
      />
      
      <ProFormSelect
        name="status"
        label="状态"
        placeholder="请选择状态"
        initialValue={1}
        options={[
          {
            label: '启用',
            value: 1,
          },
          {
            label: '禁用',
            value: 0,
          },
        ]}
        rules={[
          {
            required: true,
            message: '状态为必选项',
          },
        ]}
      />
    </ModalForm>
  );
};

export default CreateForm;
