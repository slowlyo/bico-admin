import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
  formApi: {
    getValues: vi.fn(),
    resetForm: vi.fn(),
    setValues: vi.fn(),
    validate: vi.fn(),
  },
  message: { success: vi.fn() },
  modal: {} as any,
  modalApi: {
    close: vi.fn(),
    getData: vi.fn(),
    lock: vi.fn(),
    setState: vi.fn(),
    unlock: vi.fn(),
  },
}));

vi.mock('antdv-next', () => ({ message: mocks.message }));
vi.mock('@vben/common-ui', () => ({
  useVbenModal: (options: any) => {
    mocks.modal = options;
    return [{}, mocks.modalApi];
  },
}));
vi.mock('#/adapter/form', () => ({
  useVbenForm: () => [{}, mocks.formApi],
}));

const { useCrudForm } = await import('./use-crud-form');

interface Row {
  id: number;
  name?: string;
}

function setup() {
  const options = {
    create: vi.fn().mockResolvedValue(undefined),
    entityName: '用户',
    onSuccess: vi.fn(),
    schema: [],
    toValues: vi.fn((row?: Row) => ({ name: row?.name ?? '' })),
    update: vi.fn().mockResolvedValue(undefined),
  };
  useCrudForm<Row>(options);
  return options;
}

describe('useCrudForm 提交', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mocks.formApi.validate.mockResolvedValue({ valid: true });
    mocks.formApi.getValues.mockResolvedValue({ name: 'a' });
  });

  it('校验失败时不提交也不锁定', async () => {
    const options = setup();
    mocks.formApi.validate.mockResolvedValue({ valid: false });
    await mocks.modal.onConfirm();
    expect(options.create).not.toHaveBeenCalled();
    expect(mocks.modalApi.lock).not.toHaveBeenCalled();
  });

  it('没有 id 走创建，成功后提示、通知并关闭', async () => {
    const options = setup();
    mocks.modalApi.getData.mockReturnValue({});
    await mocks.modal.onConfirm();
    expect(options.create).toHaveBeenCalledWith({ name: 'a' });
    expect(options.update).not.toHaveBeenCalled();
    expect(mocks.message.success).toHaveBeenCalledWith('创建成功');
    expect(options.onSuccess).toHaveBeenCalled();
    expect(mocks.modalApi.close).toHaveBeenCalled();
    expect(mocks.modalApi.unlock).toHaveBeenCalled();
  });

  it('有 id 走更新并带上行数据', async () => {
    const options = setup();
    const row = { id: 9, name: 'old' };
    mocks.modalApi.getData.mockReturnValue(row);
    await mocks.modal.onConfirm();
    expect(options.update).toHaveBeenCalledWith(row, { name: 'a' });
    expect(mocks.message.success).toHaveBeenCalledWith('更新成功');
  });

  it('返回 false 表示前置校验失败，弹窗保持打开且仍然解锁', async () => {
    const options = setup();
    mocks.modalApi.getData.mockReturnValue({});
    options.create.mockResolvedValue(false);
    await mocks.modal.onConfirm();
    expect(mocks.message.success).not.toHaveBeenCalled();
    expect(options.onSuccess).not.toHaveBeenCalled();
    expect(mocks.modalApi.close).not.toHaveBeenCalled();
    expect(mocks.modalApi.unlock).toHaveBeenCalled();
  });

  it('接口抛错时不关闭弹窗，也不会遗留锁定状态', async () => {
    const options = setup();
    mocks.modalApi.getData.mockReturnValue({});
    options.create.mockRejectedValue(new Error('boom'));
    await expect(mocks.modal.onConfirm()).rejects.toThrow('boom');
    expect(mocks.modalApi.close).not.toHaveBeenCalled();
    expect(mocks.modalApi.unlock).toHaveBeenCalled();
  });
});

describe('useCrudForm 打开回填', () => {
  beforeEach(() => vi.clearAllMocks());

  it('关闭时不做任何事', async () => {
    setup();
    await mocks.modal.onOpenChange(false);
    expect(mocks.formApi.resetForm).not.toHaveBeenCalled();
  });

  it('新建：标题为新建，toValues 收到 undefined', async () => {
    const options = setup();
    mocks.modalApi.getData.mockReturnValue({});
    await mocks.modal.onOpenChange(true);
    expect(mocks.modalApi.setState).toHaveBeenCalledWith({ title: '新建用户' });
    expect(options.toValues).toHaveBeenCalledWith(undefined);
    expect(mocks.formApi.resetForm).toHaveBeenCalled();
    expect(mocks.formApi.setValues).toHaveBeenCalledWith({ name: '' });
  });

  it('编辑：标题为编辑，toValues 收到行数据', async () => {
    const options = setup();
    const row = { id: 3, name: 'x' };
    mocks.modalApi.getData.mockReturnValue(row);
    await mocks.modal.onOpenChange(true);
    expect(mocks.modalApi.setState).toHaveBeenCalledWith({ title: '编辑用户' });
    expect(options.toValues).toHaveBeenCalledWith(row);
    expect(mocks.formApi.setValues).toHaveBeenCalledWith({ name: 'x' });
  });
});
