import './style.css';

import {
  AlignCenterOutlined,
  AlignLeftOutlined,
  AlignRightOutlined,
  BoldOutlined,
  ClearOutlined,
  CodeOutlined,
  ItalicOutlined,
  LinkOutlined,
  OrderedListOutlined,
  PictureOutlined,
  RedoOutlined,
  StrikethroughOutlined,
  UnderlineOutlined,
  UndoOutlined,
  UnorderedListOutlined,
  VideoCameraOutlined,
} from '@ant-design/icons';
import { Node, mergeAttributes } from '@tiptap/core';
import type { Editor } from '@tiptap/react';
import { EditorContent, useEditor } from '@tiptap/react';
import Image from '@tiptap/extension-image';
import Link from '@tiptap/extension-link';
import Placeholder from '@tiptap/extension-placeholder';
import TextAlign from '@tiptap/extension-text-align';
import Underline from '@tiptap/extension-underline';
import StarterKit from '@tiptap/starter-kit';
import { Button, Divider, Space, Tooltip, Typography, Upload, message } from 'antd';
import type { UploadProps } from 'antd';
import React, { useEffect, useMemo, useState } from 'react';

import { uploadForEditor } from '@/services/common/upload';

export interface RichEditorProps {
  /** HTML 内容 */
  value?: string;
  /** 内容变更回调 */
  onChange?: (html: string) => void;
  /** 是否禁用 */
  disabled?: boolean;
  /** 编辑器高度 */
  height?: number;
}

interface ToolbarButtonProps {
  title: string;
  icon?: React.ReactNode;
  text?: string;
  active?: boolean;
  disabled?: boolean;
  onClick: () => void;
}

type UploadActionType = 'image' | 'video';

/**
 * 为 TipTap 增加视频节点，保持富文本 HTML 里可直接保存和回显 video 标签。
 */
const Video = Node.create({
  name: 'video',
  group: 'block',
  atom: true,
  draggable: true,

  /**
   * 声明视频节点需要持久化的属性，避免回显已有 HTML 时丢失播放配置。
   */
  addAttributes() {
    return {
      src: {
        default: null,
      },
      controls: {
        default: true,
        parseHTML: (element) => element.hasAttribute('controls'),
        renderHTML: (attributes) => {
          // controls 为 false 时不输出属性，和原生 video 标签行为保持一致。
          if (!attributes.controls) {
            return {};
          }
          return { controls: '' };
        },
      },
    };
  },

  /**
   * 只解析带 src 的 video，避免空标签进入文档后无法播放。
   */
  parseHTML() {
    return [{ tag: 'video[src]' }];
  },

  /**
   * 输出标准 video 标签，后端保存 HTML 后无需额外转换。
   */
  renderHTML({ HTMLAttributes }) {
    return ['video', mergeAttributes({ class: 'rich-editor-video' }, HTMLAttributes)];
  },

  /**
   * 暴露插入视频命令，上传完成后由工具栏直接写入当前光标位置。
   */
  addCommands() {
    return {
      setVideo:
        (attrs: { src: string }) =>
        ({ commands }: { commands: Editor['commands'] }) => {
          return commands.insertContent({
            type: this.name,
            attrs: {
              src: attrs.src,
              controls: true,
            },
          });
        },
    };
  },
});

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    video: {
      setVideo: (attrs: { src: string }) => ReturnType;
    };
  }
}

/**
 * 工具栏按钮统一处理激活态和禁用态，保证按钮交互风格一致。
 */
const ToolbarButton: React.FC<ToolbarButtonProps> = ({
  title,
  icon,
  text,
  active,
  disabled,
  onClick,
}) => {
  return (
    <Tooltip title={title}>
      <Button
        size="small"
        type={active ? 'primary' : 'text'}
        icon={icon}
        disabled={disabled}
        onMouseDown={(event) => {
          // 阻止按钮抢走编辑器焦点，否则部分链式命令会丢失当前选区。
          event.preventDefault();
          onClick();
        }}
      >
        {text}
      </Button>
    </Tooltip>
  );
};

/**
 * 通用富文本编辑器组件，内部使用 TipTap，外部仍以 HTML 字符串对接表单和后端。
 */
const RichEditor: React.FC<RichEditorProps> = ({
  value = '',
  onChange,
  disabled,
  height = 420,
}) => {
  const [messageApi, contextHolder] = message.useMessage();
  const [uploadingType, setUploadingType] = useState<UploadActionType | null>(null);

  const extensions = useMemo(() => {
    return [
      StarterKit,
      Underline,
      Link.configure({
        openOnClick: false,
        HTMLAttributes: {
          rel: 'noopener noreferrer',
          target: '_blank',
        },
      }),
      Image.configure({
        HTMLAttributes: {
          class: 'rich-editor-image',
        },
      }),
      TextAlign.configure({
        types: ['heading', 'paragraph'],
      }),
      Placeholder.configure({
        placeholder: '请输入内容...',
      }),
      Video,
    ];
  }, []);

  const editor = useEditor({
    extensions,
    content: value,
    editable: !disabled,
    immediatelyRender: false,
    onUpdate: ({ editor: currentEditor }) => {
      // TipTap 内部维护文档结构，对外只同步 HTML，方便复用现有接口。
      onChange?.(currentEditor.getHTML());
    },
  });

  useEffect(() => {
    if (!editor) {
      return;
    }

    const nextValue = value || '';
    const currentValue = editor.getHTML();
    // 外部 value 变化时才覆盖编辑器内容，避免输入过程中重置光标。
    if (nextValue !== currentValue) {
      editor.commands.setContent(nextValue, { emitUpdate: false });
    }
  }, [editor, value]);

  useEffect(() => {
    if (!editor) {
      return;
    }

    editor.setEditable(!disabled);
  }, [disabled, editor]);

  /**
   * 处理链接设置；已有链接会被当前输入替换，空输入代表取消链接。
   */
  const handleSetLink = () => {
    if (!editor) {
      return;
    }

    const previousUrl = editor.getAttributes('link').href as string | undefined;
    const url = window.prompt('请输入链接地址', previousUrl || '');
    // 用户取消输入时不改动当前选区，避免误删已有链接。
    if (url === null) {
      return;
    }

    // 空字符串表示主动清除链接，与取消弹窗区分开。
    if (url.trim() === '') {
      editor.chain().focus().extendMarkRange('link').unsetLink().run();
      return;
    }

    editor.chain().focus().extendMarkRange('link').setLink({ href: url.trim() }).run();
  };

  /**
   * 通过现有上传接口上传媒体文件，并插入对应 TipTap 节点。
   */
  const uploadMedia = async (file: File, type: UploadActionType) => {
    if (!editor) {
      return;
    }

    setUploadingType(type);
    try {
      const resp = await uploadForEditor(file, type);
      // 上传接口沿用编辑器响应结构，缺少 url 时不能写入富文本内容。
      if (!resp || resp.errno !== 0 || !resp.data?.url) {
        throw new Error(resp?.message || '上传失败');
      }

      // 图片和视频对应不同节点，分别调用扩展命令以保留语义。
      if (type === 'image') {
        editor.chain().focus().setImage({ src: resp.data.url }).run();
      } else {
        editor.chain().focus().setVideo({ src: resp.data.url }).run();
      }
    } catch (error: any) {
      messageApi.error(error?.message || '上传失败');
    } finally {
      setUploadingType(null);
    }
  };

  /**
   * 生成 Ant Design Upload 配置，禁止自动列表展示，由编辑器内容承载结果。
   */
  const buildUploadProps = (type: UploadActionType): UploadProps => {
    return {
      accept: type === 'image' ? 'image/*' : 'video/*',
      showUploadList: false,
      disabled: !editor || disabled || uploadingType !== null,
      beforeUpload: (file) => {
        void uploadMedia(file, type);
        return Upload.LIST_IGNORE;
      },
    };
  };

  const toolbarDisabled = !editor || disabled;

  return (
    <div className="rich-editor">
      {contextHolder}
      <div className="rich-editor-toolbar">
        <Space size={4} wrap>
          <ToolbarButton
            title="撤销"
            icon={<UndoOutlined />}
            disabled={toolbarDisabled || !editor?.can().undo()}
            onClick={() => editor?.chain().focus().undo().run()}
          />
          <ToolbarButton
            title="重做"
            icon={<RedoOutlined />}
            disabled={toolbarDisabled || !editor?.can().redo()}
            onClick={() => editor?.chain().focus().redo().run()}
          />
          <Divider type="vertical" />
          <ToolbarButton
            title="加粗"
            icon={<BoldOutlined />}
            active={editor?.isActive('bold')}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleBold().run()}
          />
          <ToolbarButton
            title="斜体"
            icon={<ItalicOutlined />}
            active={editor?.isActive('italic')}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleItalic().run()}
          />
          <ToolbarButton
            title="下划线"
            icon={<UnderlineOutlined />}
            active={editor?.isActive('underline')}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleUnderline().run()}
          />
          <ToolbarButton
            title="删除线"
            icon={<StrikethroughOutlined />}
            active={editor?.isActive('strike')}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleStrike().run()}
          />
          <ToolbarButton
            title="行内代码"
            icon={<CodeOutlined />}
            active={editor?.isActive('code')}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleCode().run()}
          />
          <Divider type="vertical" />
          <ToolbarButton
            title="标题 1"
            text="H1"
            active={editor?.isActive('heading', { level: 1 })}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleHeading({ level: 1 }).run()}
          />
          <ToolbarButton
            title="标题 2"
            text="H2"
            active={editor?.isActive('heading', { level: 2 })}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleHeading({ level: 2 }).run()}
          />
          <ToolbarButton
            title="引用"
            text="“”"
            active={editor?.isActive('blockquote')}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleBlockquote().run()}
          />
          <ToolbarButton
            title="无序列表"
            icon={<UnorderedListOutlined />}
            active={editor?.isActive('bulletList')}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleBulletList().run()}
          />
          <ToolbarButton
            title="有序列表"
            icon={<OrderedListOutlined />}
            active={editor?.isActive('orderedList')}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().toggleOrderedList().run()}
          />
          <Divider type="vertical" />
          <ToolbarButton
            title="左对齐"
            icon={<AlignLeftOutlined />}
            active={editor?.isActive({ textAlign: 'left' })}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().setTextAlign('left').run()}
          />
          <ToolbarButton
            title="居中"
            icon={<AlignCenterOutlined />}
            active={editor?.isActive({ textAlign: 'center' })}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().setTextAlign('center').run()}
          />
          <ToolbarButton
            title="右对齐"
            icon={<AlignRightOutlined />}
            active={editor?.isActive({ textAlign: 'right' })}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().setTextAlign('right').run()}
          />
          <Divider type="vertical" />
          <ToolbarButton
            title="链接"
            icon={<LinkOutlined />}
            active={editor?.isActive('link')}
            disabled={toolbarDisabled}
            onClick={handleSetLink}
          />
          <Upload {...buildUploadProps('image')}>
            <Tooltip title="上传图片">
              <Button
                size="small"
                type="text"
                icon={<PictureOutlined />}
                disabled={!editor || disabled || uploadingType !== null}
                loading={uploadingType === 'image'}
              />
            </Tooltip>
          </Upload>
          <Upload {...buildUploadProps('video')}>
            <Tooltip title="上传视频">
              <Button
                size="small"
                type="text"
                icon={<VideoCameraOutlined />}
                disabled={!editor || disabled || uploadingType !== null}
                loading={uploadingType === 'video'}
              />
            </Tooltip>
          </Upload>
          <Divider type="vertical" />
          <ToolbarButton
            title="清除格式"
            icon={<ClearOutlined />}
            disabled={toolbarDisabled}
            onClick={() => editor?.chain().focus().unsetAllMarks().clearNodes().run()}
          />
        </Space>
      </div>

      <div className="rich-editor-content" style={{ minHeight: height }}>
        {editor ? (
          <EditorContent editor={editor} />
        ) : (
          <Typography.Text type="secondary">编辑器加载中...</Typography.Text>
        )}
      </div>
    </div>
  );
};

export default RichEditor;
