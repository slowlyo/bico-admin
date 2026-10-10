import type { VbenFormSchema } from '#/adapter/form';

import { useVbenModal } from '@vben/common-ui';

import { message } from 'antdv-next';

import { useVbenForm } from '#/adapter/form';

interface CrudFormOptions<T extends { id: number }> {
  // 新建。返回 false 表示前置校验未通过，弹窗保持打开且不提示成功。
  create: (values: Record<string, any>) => Promise<unknown>;
  // 用于弹窗标题，如“用户”。
  entityName: string;
  // 提交成功后通知列表刷新，通常是 emit('success')。
  onSuccess: () => void;
  schema: VbenFormSchema[];
  // 把行数据转成表单值；新建时 row 为空，在这里给出默认值。
  toValues: (row?: T) => Record<string, any>;
  // 编辑。返回 false 的含义同 create。
  update: (row: T, values: Record<string, any>) => Promise<unknown>;
}

// 标准新建/编辑弹窗：打开时回填、确认时校验并提交、提交期间锁定防重复。
// 列表页通过 setData(row) 传入行数据表示编辑，setData({}) 表示新建。
export function useCrudForm<T extends { id: number }>(
  options: CrudFormOptions<T>,
) {
  const [Form, formApi] = useVbenForm({
    schema: options.schema,
    showDefaultActions: false,
  });

  const [Modal, modalApi] = useVbenModal({
    async onConfirm() {
      const { valid } = await formApi.validate();
      if (!valid) {
        return;
      }
      const values = await formApi.getValues();
      const row = modalApi.getData<T>();
      const editing = Boolean(row?.id);
      modalApi.lock();
      try {
        const result = editing
          ? await options.update(row, values)
          : await options.create(values);
        if (result === false) {
          // 业务前置校验失败：由调用方已提示原因，保持弹窗打开让用户修改。
          return;
        }
        message.success(editing ? '更新成功' : '创建成功');
        options.onSuccess();
        modalApi.close();
      } finally {
        modalApi.unlock();
      }
    },
    async onOpenChange(isOpen) {
      if (!isOpen) {
        return;
      }
      const row = modalApi.getData<T>();
      const editing = Boolean(row?.id);
      modalApi.setState({
        title: `${editing ? '编辑' : '新建'}${options.entityName}`,
      });
      await formApi.resetForm();
      await formApi.setValues(options.toValues(editing ? row : undefined));
    },
  });

  return { Form, formApi, Modal, modalApi };
}
